package userControllers

import (
	"CMS/app/models"
	"CMS/app/service/codeservice"
	userservice "CMS/app/service/userservice"
	"CMS/app/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterData struct {
	UserName    string `json:"name" binding:"required"`
	Account     string `json:"account" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Sex         int    `json:"sex" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Phonenumber string `json:"phonenumber" binding:"required"`
	Type        int    `json:"type" binding:"required"`
}

// 注册
func Register(c *gin.Context) {
	var data RegisterData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}

	//判断账号合法性
	flag := userservice.CheckAccountlegitimacy(data.Account)
	if !flag {
		utils.JsonErrorResponse(c, 1, "账号只能由数字组成!")
		return
	}

	//判断邮箱合法性
	flag = userservice.CheckEmaillegitimacy(data.Email)
	if !flag {
		utils.JsonErrorResponse(c, 1, "邮箱格式不正确！")
		return
	}
	//判断手机号合法性
	flag = userservice.CheckPhonenumberlegitimacy(data.Phonenumber)
	if !flag {
		utils.JsonErrorResponse(c, 1, "手机号格式不正确！")
		return
	}
	//判断密码长度大于8位
	flag = userservice.CheckPasswordLength(data.Password)
	if flag {
		utils.JsonErrorResponse(c, 1, "密码长度最少为8位！")
		return
	}
	//判断密码格式是否正确
	flag = userservice.CheckPasswordFormat(data.Password)
	if !flag {
		utils.JsonErrorResponse(c, 1, "密码需同时包含大写字母、小写字母和数字,并且不能包含空格！")
		return
	}
	// 判断账号是否已经注册
	err = userservice.CheckUserExistByAccount(data.Account)
	if err == nil {
		utils.JsonErrorResponse(c, 4, "账号已经被注册！")
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		utils.JsonInternalServerErrorResponse(c)
		return
	}

	//判断用户名是否已经被使用
	err = userservice.CheckUserExistByName(data.UserName)
	if err == nil {
		utils.JsonErrorResponse(c, 4, "该用户名已经被使用！")
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		utils.JsonInternalServerErrorResponse(c)
		return
	}

	//判断邮箱是否已经绑定
	err = userservice.CheckUserExistByEmail(data.Email)
	if err == nil {
		utils.JsonErrorResponse(c, 4, "邮箱已经被绑定！")
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		utils.JsonInternalServerErrorResponse(c)
		return
	}
	//判断手机号是否已经绑定
	err = userservice.CheckUserExistByPhonenumber(data.Phonenumber)
	if err == nil {
		utils.JsonErrorResponse(c, 4, "手机号已经被绑定！")
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		utils.JsonInternalServerErrorResponse(c)
		return
	}

	flag = userservice.CheckTypelegitimacy(data.Type)
	if flag {
		utils.JsonErrorResponse(c, 5, "用户类型只能为管理员或普通用户")
		return
	}

	// 注册用户
	err = userservice.Register(models.User{
		UserName:      data.UserName,
		Account:       data.Account,
		Password:      codeservice.AesEncryptCBC([]byte(data.Password)),
		Sex:           data.Sex,
		Email:         data.Email,
		Type:          data.Type,
		Phonenumber:   data.Phonenumber,
		Teambelonging: "null",
	})
	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
