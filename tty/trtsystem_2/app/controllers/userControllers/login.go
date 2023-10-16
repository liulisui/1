package userControllers

import (
	"CMS/app/midware"
	"CMS/app/models"
	"CMS/app/service/codeservice"
	userService "CMS/app/service/userservice"
	"CMS/app/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginData struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录

func Login(c *gin.Context) {
	r := gin.Default()
	store := cookie.NewStore([]byte("ilovejinghong"))
	r.Use(sessions.Sessions("mysession", store))

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
	session := sessions.Default(c)
	session.Set("username", user.UserName)
	session.Save()

	c.Redirect(200, "/home")

	utils.JsonSuccessResponse(c, user)

	r.GET("/dashboard", midware.AuthMiddleware, func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")

		c.String(200201, "欢迎访问仪表盘，%s!", username)
	})

	// 注销路由
	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("username")
		session.Save()

		c.Redirect(200202, "/login")
	})

	r.Run(":8080")
}
