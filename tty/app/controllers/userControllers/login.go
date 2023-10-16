package userControllers

import (
	"CMS/app/models"
	"CMS/app/service/codeservice"
	userService "CMS/app/service/userservice"
	"CMS/app/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginData struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录
func Login(c *gin.Context) {
	var data LoginData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}
	//判断是否存在
	err = userService.CheckUserExistByAccount(data.Account)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, 4, "用户不存在!")
		} else {
			utils.JsonInternalServerErrorResponse(c)
		}
		return
	}
	//获取用户信息
	var user *models.User
	user, err = userService.GetUserByAccount(data.Account)
	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}
	//判断密码正确
	flag := userService.ComparePwd(data.Password, string(codeservice.AesDecryptCBC(user.Password)))
	if !flag {
		utils.JsonErrorResponse(c, 2, "密码错误!")
		return
	}

	utils.JsonSuccessResponse(c, user)
}
