package userControllers

import (
	"CMS/app/models"
	"CMS/app/service/codeservice"
	userService "CMS/app/service/userservice"
	"CMS/app/utils"

	"github.com/gin-gonic/gin"
)

type PersonalData struct {
	ID          uint   `json:"id"`
	UserName    string `json:"name" binding:"required"`
	Account     string `json:"account" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Repassword  string `json:"repassword" binding:"required"`
	Sex         int    `json:"sex" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Phonenumber string `json:"phonenumber" binding:"required"`
}

func InformationUpdata(c *gin.Context) {
	var data PersonalData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}

	//判断邮箱合法性
	flag := userService.CheckEmaillegitimacy(data.Email)
	if !flag {
		utils.JsonErrorResponse(c, 400, "邮箱格式不正确！")
		return
	}
	//判断手机号合法性
	flag = userService.CheckPhonenumberlegitimacy(data.Phonenumber)
	if !flag {
		utils.JsonErrorResponse(c, 400, "手机号格式不正确！")
		return
	}
	//判断密码长度大于8位
	flag = userService.CheckPasswordLength(data.Password)
	if flag {
		utils.JsonErrorResponse(c, 400, "密码长度最少为8位!")
		return
	}
	//判断密码格式是否正确
	flag = userService.CheckPasswordFormat(data.Password)
	if !flag {
		utils.JsonErrorResponse(c, 400, "密码需同时包含大写字母、小写字母和数字,并且不能包含空格！")
		return
	}
	//判断两次输入密码是否一致
	flag = userService.ComparePwd(data.Password, data.Repassword)

	if !flag {
		utils.JsonErrorResponse(c, 400, "两次输入的密码不相同！")
		return
	}

	//修改信息
	DataOri, _ := userService.GetUserByAccount(data.Account)
	err = userService.UpdateUser(models.User{
		ID:            DataOri.ID,
		UserName:      data.UserName,
		Account:       data.Account,
		Password:      codeservice.AesEncryptCBC([]byte(data.Password)),
		Sex:           data.Sex,
		Email:         data.Email,
		Phonenumber:   data.Phonenumber,
		Type:          DataOri.Type,
		Teambelonging: DataOri.Teambelonging,
	})
	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}

	utils.JsonSuccessResponse(c, nil)

}

type User struct {
	Account string `json:"account"`
}

func InformationGet(c *gin.Context) {
	// var data User
	var user *models.User
	account := c.Query("account")
	var err error
	user, err = userService.GetUserByAccount(account)
	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}

	utils.JsonSuccessResponse(c, user)
}
