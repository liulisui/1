package teamControllers

import (
	"CMS/app/models"
	"CMS/app/service/codeservice"
	teamservice "CMS/app/service/teamservice"
	userService "CMS/app/service/userservice"
	"CMS/app/utils"
	"CMS/config/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateTeamData struct {
	ID         uint   `json:"team_id"`
	TeamName   string `json:"team_name" binding:"required"`
	TeamLeader string `json:"leader_name" binding:"required"`
	Number     int    `json:"number" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Total      int    `json:"total"`
	State      int    `json:"state"`
}

func CreateTeam(c *gin.Context) {
	var data CreateTeamData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}
	//默认状态为未报名，设为1
	data.State = 1

	//默认总数是1位成员
	data.Total = 1

	if data.Number < 4 || data.Number > 6 {
		utils.JsonErrorResponse(c, 5, "人数只能为4到6的整数！")
		return
	}
	if data.TeamName == "null" {
		utils.JsonErrorResponse(c, 5, "不能设置该团队名！")
		return
	}
	err = teamservice.CheckTeamExistByName(data.TeamName)
	if err == nil {
		utils.JsonErrorResponse(c, 4, "该团队已经存在！")
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		utils.JsonInternalServerErrorResponse(c)
		return
	}
	//队长是否存在
	err = userService.CheckUserExistByName(data.TeamLeader)
	if err != nil {
		utils.JsonErrorResponse(c, 3, "队长不存在！")
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		utils.JsonInternalServerErrorResponse(c)
		return
	}
	//判断密码长度大于8位
	flag := userService.CheckPasswordLength(data.Password)
	if flag {
		utils.JsonErrorResponse(c, 1, "密码长度最少为8位！")
		return
	}
	//判断密码格式是否正确
	flag = userService.CheckPasswordFormat(data.Password)
	if !flag {
		utils.JsonErrorResponse(c, 1, "密码需同时包含大写字母、小写字母和数字,并且不能包含空格！")
		return
	}
	err = teamservice.CreateTeam(models.Team{
		ID:         data.ID,
		TeamName:   data.TeamName,
		TeamLeader: data.TeamLeader,
		Number:     data.Number,
		Password:   codeservice.AesEncryptCBC([]byte(data.Password)),
		Total:      data.Total,
		State:      data.State,
	})

	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}

	var user *models.User
	database.DB.Table("users").Where("user_name =?", data.TeamLeader).First(&user)
	user.Teambelonging = data.TeamName
	database.DB.Save(&user)
	utils.JsonSuccessResponse(c, nil)
}

type UpdateTeamData struct {
	ID         uint   `json:"team_id" `
	TeamName   string `json:"team_name" binding:"required"`
	Number     int    `json:"number" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Repassword string `json:"Repassword" binding:"required"`
}

// 更新团队信息
func UpdateTeam_(c *gin.Context) {
	var data UpdateTeamData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}

	result := database.DB.Table("teams").Where("team_name =?", data.TeamName).First(&models.Team{})
	if result.Error == nil {
		utils.JsonErrorResponse(c, 4, "该队伍已存在！")
		return
	}

	if data.Number < 4 || data.Number > 6 {
		utils.JsonErrorResponse(c, 5, "人数只能为4到6的整数！")
		return
	}

	if !userService.CheckPasswordFormat(data.Password) {
		utils.JsonErrorResponse(c, 1, "密码需同时包含大写字母、小写字母和数字,并且不能包含空格！")
		return
	}

	if userService.CheckPasswordLength(data.Password) {
		utils.JsonErrorResponse(c, 1, "密码长度小于8位！")
		return
	}

	if data.Password != data.Repassword {
		utils.JsonErrorResponse(c, 5, "两次密码不相同！")
		return
	}
	err = teamservice.UpdateTeam(models.Team{
		ID:       data.ID,
		TeamName: data.TeamName,
		Password: codeservice.AesEncryptCBC([]byte(data.Password)),
		Number:   data.Number,
	})
	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type DeleteTeamData struct {
	TeamName string `json:"team_name" binding:"required"`
	UserName string `json:"user_name" binding:"required"`
}

func DeleteTeam(c *gin.Context) {
	var data DeleteTeamData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}
	var team *models.Team
	team, err = teamservice.CheckTeams(data.TeamName)
	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}
	if team.TeamLeader != data.UserName {
		utils.JsonErrorResponse(c, 5, "只有队长可以解散团队!")
	}

	err = teamservice.DeleteTeam(data.UserName, data.TeamName)
	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type GetTeamData struct {
	UserAccount string `json:"user_account" binding:"required"`
}
type TeamData struct {
	ID         uint   `json:"team_id"`
	TeamName   string `json:"team_name" binding:"required"`
	TeamLeader string `json:"leader_name" binding:"required"`
	Number     int    `json:"number" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Total      int    `json:"total"`
	State      int    `json:"state"`
	Member     []models.User
}

func GetTeam(c *gin.Context) {
	var data GetTeamData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}
	var User models.User
	result := database.DB.Where("account = ?", data.UserAccount).First(&User)
	if result.Error != nil {
		utils.JsonErrorResponse(c, 3, "该用户不存在！")
		return
	}

	TeamName := User.Teambelonging
	if TeamName == "null" {
		utils.JsonErrorResponse(c, 3, "你不在团队中！")
		return
	}

	var Teamorigin models.Team
	result = database.DB.Where("team_name =?", User.Teambelonging).First(&Teamorigin)
	if result.Error != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}

	var memberList []models.User
	memberList, err = teamservice.GetTeamMember(TeamName)
	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}
	TeamNew := TeamData{
		ID:         Teamorigin.ID,
		TeamName:   Teamorigin.TeamName,
		TeamLeader: Teamorigin.TeamLeader,
		Number:     Teamorigin.Number,
		Password:   string(codeservice.AesDecryptCBC(Teamorigin.Password)),
		Member:     memberList,
		State:      Teamorigin.State,
		Total:      Teamorigin.Total,
	}

	utils.JsonSuccessResponse(c, TeamNew)
}

type Joindata struct {
	TeamName    string `json:"team_name" binding:"required"`
	UserAccount string `json:"user_account" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

func Join(c *gin.Context) {
	var data Joindata
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}
	err = teamservice.CheckTeamExistByName(data.TeamName)
	if err != nil {
		utils.JsonErrorResponse(c, 3, "团队不存在！")
		return
	}
	var user models.User
	result := database.DB.Where("account = ?", data.UserAccount).First(&user)
	if result.Error != nil {
		utils.JsonErrorResponse(c, 3, "用户不存在！")
		return
	}
	var team *models.Team
	team, err = teamservice.CheckTeams(data.TeamName)
	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}
	flag := teamservice.CompareTeamPwd(string(codeservice.AesDecryptCBC(team.Password)), data.Password)
	if !flag {
		utils.JsonErrorResponse(c, 2, "密码错误!")
		return
	}

	if user.Teambelonging == data.TeamName {
		utils.JsonErrorResponse(c, 5, "你已报名该团队！")
		return
	}

	if team.Number == team.Total {
		utils.JsonErrorResponse(c, 5, "该团队人数已满！")
		return
	}
	err = teamservice.JoinTeam(data.TeamName, data.UserAccount)
	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}

	team.Total += 1
	utils.JsonSuccessResponse(c, nil)
}

type signdata struct {
	TeamName string `json:"team_name" binding:"required"`
}

func Sign(c *gin.Context) {
	var data signdata
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}

	var team models.Team
	result := database.DB.Table("teams").Where("team_name =?", data.TeamName).First(&team)
	if result.Error != nil {
		utils.JsonErrorResponse(c, 3, "队伍不存在！")
		return
	}

	if team.State == 2 {
		utils.JsonErrorResponse(c, 4, "已经报名了！")
		return
	}

	team.State = 2
	database.DB.Save(&team)

	utils.JsonSuccessResponse(c, nil)
}

type canceldata struct {
	TeamName string `json:"team_name" binding:"required"`
}

func Cancel(c *gin.Context) {
	var data canceldata
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}

	var team models.Team
	result := database.DB.Table("teams").Where("team_name =?", data.TeamName).First(&team)
	if result.Error != nil {
		utils.JsonErrorResponse(c, 3, "队伍不存在！")
		return
	}

	if team.State == 1 {
		utils.JsonErrorResponse(c, 4, "未处于报名状态!")
		return
	}

	team.State = 1
	database.DB.Save(&team)

	utils.JsonSuccessResponse(c, nil)
}

type DeleteTeammatedata struct {
	Account  string `json:"account" binding:"required"`
	Teamname string `json:"team_name" binding:"required"`
}

func DeleteTeammate(c *gin.Context) {
	var data DeleteTeamData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}

	err = userService.CheckUserExistByAccount(data.UserName)
	if err != nil {
		utils.JsonErrorResponse(c, 3, "用户不存在")
		return
	}

	var user models.User
	database.DB.Table("users").Where("user_name =?", data.UserName).First(&user)
	if user.Teambelonging != data.TeamName {
		utils.JsonErrorResponse(c, 3, "用户不在队伍中！")
		return
	}

	user.Teambelonging = "null"
	database.DB.Save(&user)

	utils.JsonSuccessResponse(c, nil)
}

type enterdata struct {
	Teamname string `json:"team_name" binding:"required"`
	Account  string `json:"account" binding:"required"`
}

func Enter(c *gin.Context) {
	var data enterdata
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}

}
