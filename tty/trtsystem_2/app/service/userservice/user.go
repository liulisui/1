package userService

import (
	"CMS/app/models"
	"CMS/config/database"
	"regexp"
	"strings"
)

func CheckUserExistByAccount(account string) error {
	result := database.DB.Where("account = ?", account).First(&models.User{})
	return result.Error
}

func CheckPasswordFormat(Password string) bool {
	pd1 := 0
	pd2 := 0
	pd3 := 0
	for i := 0; i < len(Password); i++ {
		if Password[i] == ' ' {
			return false
		}
		if Password[i] <= '9' && Password[i] >= '0' {
			pd1 = 1
		}
		if Password[i] <= 'Z' && Password[i] >= 'A' {
			pd2 = 1
		}
		if Password[i] <= 'z' && Password[i] >= 'a' {
			pd3 = 1
		}
	}
	if pd1 == 1 && pd2 == 1 && pd3 == 1 {
		return true
	}
	return false

}

func CheckUserExistByPhonenumber(phonenumber string) error {
	result := database.DB.Table("users").Where("phonenumber = ?", phonenumber).First(&models.User{})
	return result.Error
}

func CheckUserExistByEmail(email string) error {
	result := database.DB.Table("users").Where("email = ?", email).First(&models.User{})
	return result.Error
}

func CheckUserExistByName(name string) error {
	result := database.DB.Table("users").Where("user_name = ?", name).First(&models.User{})
	return result.Error
}

func GetUserByAccount(account string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("account = ?", account).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByName(name string) (*models.User, error) {
	var user models.User
	result := database.DB.Table("users").Where("user_name = ?", name).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func ComparePwd(pwd1 string, pwd2 string) bool {
	return pwd1 == pwd2
}
func Register(user models.User) error {
	result := database.DB.Create(&user)
	return result.Error
}
func CheckAccountlegitimacy(account string) bool {
	match, _ := regexp.MatchString("^[0-9]+$", account)
	return match
}
func CheckEmaillegitimacy(email string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$", email)
	return match
}
func CheckPhonenumberlegitimacy(phonenumber string) bool {
	match, _ := regexp.MatchString("^1[3-9][0-9]{9}$", phonenumber)
	return match
}
func CheckTypelegitimacy(Type int) bool {
	if (Type > 2) || (Type < 1) {
		return true
	} else {
		return false
	}
}
func CheckPasswordLength(Password string) bool {
	PawLength := strings.Count(Password, "")
	if PawLength < 8 {
		return true
	} else {
		return false
	}
}

func CheckUser(UserId int) (*models.User, error) {
	var user models.User
	result := database.DB.Where("id = ?", UserId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func UpdateUser(user models.User) error {
	result := database.DB.Save(&user)
	return result.Error
}
