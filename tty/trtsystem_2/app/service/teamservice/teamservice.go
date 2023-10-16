package classservice

import (
	"CMS/app/models"
	"CMS/config/database"
)

func CheckTeamExistByName(TeamName string) error {
	result := database.DB.Where("team_name = ?", TeamName).First(&models.Team{})
	return result.Error
}
func CreateTeam(Team models.Team) error {
	result := database.DB.Create(&Team)
	return result.Error
}
func UpdateTeam(Team models.Team) error {
	result := database.DB.Save(&Team)
	return result.Error
}
func DeleteTeam(user_name string, team_name string) error {
	result := database.DB.Where("TeamLeader = ? AND TeamName=? ", user_name, team_name).Delete(&models.Team{})
	return result.Error
}
func GetTeamMember(team_name string) ([]models.User, error) {
	var memberList []models.User
	result := database.DB.Table("users").Where("teambelonging = ?", team_name).Find(&memberList)
	if result.Error != nil {
		return nil, result.Error
	}
	return memberList, nil
}
func CheckTeams(team_name string) (*models.Team, error) {
	var team models.Team
	result := database.DB.Table("teams").Where("team_name = ?", team_name).First(&team)
	if result.Error != nil {
		return nil, result.Error
	}
	return &team, nil
}
func CompareTeamPwd(pwd1 string, pwd2 string) bool {
	return pwd1 == pwd2
}
func JoinTeam(team_name string, user_account string) error {
	var user models.User
	result := database.DB.Where("Account = ?", user_account).First(&user)
	if result.Error != nil {
		return result.Error
	}
	user.Teambelonging = team_name
	database.DB.Save(&user)
	return nil
}
