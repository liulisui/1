package router

import (
	"CMS/app/controllers/teamControllers"
	userControllers "CMS/app/controllers/userControllers"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	const pre = "/api"
	api := r.Group(pre)
	{
		api.POST("/login", userControllers.Login)
		api.POST("/register", userControllers.Register)
		api.GET("/personalinformation", userControllers.InformationGet)
		api.PUT("/personalinformation", userControllers.InformationUpdata)
		api.POST("/createteam", teamControllers.CreateTeam)
		api.POST("/jointeam", teamControllers.Join)
		api.GET("/teaminformation", teamControllers.GetTeam)
		team := api.Group("/team")
		{
			team.PUT("/update", teamControllers.UpdateTeam_)
			team.DELETE("/deleteteam", teamControllers.DeleteTeam)
			team.DELETE("/deleteteammate", teamControllers.DeleteTeammate)
			team.POST("/sign", teamControllers.Sign)
			team.POST("/cancel", teamControllers.Cancel)
			team.POST("/enter", teamControllers.Enter)
		}
	}

}
