package service

import (
	"todo_list/model"
	"todo_list/package/utils"
	"todo_list/serializer"

	"github.com/jinzhu/gorm"
)

type UserService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=5,max=15"`
	Password string `form:"password" json:"password" binding:"required,min=5,max=16"`
}

type UserBaseInfoService struct {
	Email    string `form:"email" json:"email"`
	Phone    string `form:"phone" json:"phone"`
	UserName string `form:"user_name" json:"user_name"`
}

// 用户注册
func (service *UserService) Register() serializer.Response {
	var user model.User
	var count int
	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).First(&user).Count(&count)
	if count == 1 {
		return serializer.Response{
			Status: 400,
			Msg:    "User has Registered!",
		}
	}
	// 注册校验
	user.UserName = service.UserName
	if err := user.SetPassWord(service.Password); err != nil {
		return serializer.Response{
			Status: 400,
			Msg:    err.Error(),
		}
	}
	// 创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "Register Failed!",
		}
	}
	return serializer.Response{
		Status: 200,
		Msg:    "Register Success!",
	}
}

// 用户登录
func (service *UserService) Login() serializer.Response {
	var user model.User
	// 查询账号信息
	if err := model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return serializer.Response{
				Status: 400,
				Msg:    "User don't exist!",
			}
		}
		// 其他错误
		return serializer.Response{
			Status: 500,
			Msg:    "Login Error!",
		}
	}
	// 对比账号信息
	if user.CheckPassword(service.Password) == false {
		return serializer.Response{
			Status: 400,
			Msg:    "Password Error!",
		}
	}
	// 登录成功生成Token

	// 校验用户身份
	role := "user"
	if user.Role == 0 {
		role = "admin"
	}

	token, err := utils.GenerateToken(user.ID, service.UserName, service.Password, role)
	if err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "Token Error!",
		}
	}
	return serializer.Response{
		Status: 200,
		Data:   serializer.TokenData{User: serializer.BuildUser(user), Token: token},
		Msg:    "Login Success!",
	}
}

// 用户信息
func (service *UserBaseInfoService) UserInfo(uid uint) serializer.Response {
	var user_info model.UserBaseInfo
	code := 200
	err := model.DB.First(&user_info, uid).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "search error!",
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildUserBaseInfo(user_info),
		Msg:    "user info success!",
	}
}

func (service *UserBaseInfoService) UpdateUserInfo(uid uint) serializer.Response {
	var user_info model.UserBaseInfo
	code := 200
	err := model.DB.First(&user_info, uid).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "search error!",
		}
	}
	user_info = model.UserBaseInfo{
		Email:    service.Email,
		Phone:    service.Phone,
		UserName: service.UserName,
	}

	if err := model.DB.Model(&user_info).Where("id = ?", uid).Updates(service).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "update user_info failed!",
		}
	}

	return serializer.Response{
		Status: code,
		// Data:   serializer.BuildUserBaseInfo(user_info),
		Msg: "update user info success!",
	}
}
