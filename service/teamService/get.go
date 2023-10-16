package teamService

import (
	"walk-server/global"
	"walk-server/model"
)

func GetTeamByID(id uint) (*model.Team, error) {
	user := model.Team{}
	result := global.DB.Where(
		&model.Team{
			ID: id,
		},
	).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
