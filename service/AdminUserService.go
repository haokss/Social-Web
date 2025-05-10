package service

import (
	"todo_list/model"
	"todo_list/serializer"
)

type AdminUserListService struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type AdminEditUserService struct {
	UserName string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Role     int    `json:"role"`
}

type AdminResetPasswordService struct{}

func (service *AdminUserListService) ListUsers() serializer.Response {
	var users []model.User
	var total int64

	if service.Page < 1 {
		service.Page = 1
	}
	if service.PageSize <= 0 {
		service.PageSize = 10
	}

	offset := (service.Page - 1) * service.PageSize

	// 查询总数
	if err := model.DB.Model(&model.User{}).Count(&total).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "获取用户总数失败",
		}
	}

	// 分页查询用户
	if err := model.DB.Limit(service.PageSize).Offset(offset).Find(&users).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "获取用户列表失败",
		}
	}

	return serializer.Response{
		Status: 200,
		Data: map[string]interface{}{
			"list":  serializer.BuildAdminUsers(users),
			"total": total,
		},
		Msg: "获取用户列表成功",
	}
}

func (service *AdminEditUserService) EditUser(userID string) serializer.Response {
	var user model.User

	// 查询用户是否存在
	if err := model.DB.First(&user, userID).Error; err != nil {
		return serializer.Response{
			Status: 404,
			Msg:    "用户不存在",
		}
	}

	// 更新用户信息
	user.UserName = service.UserName
	user.Phone = service.Phone
	user.Email = service.Email
	user.Role = service.Role

	// 保存更新
	if err := model.DB.Save(&user).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "更新用户信息失败",
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "用户信息更新成功",
	}
}

func (service *AdminResetPasswordService) ResetPassword(userID string) serializer.Response {
	var user model.User

	// 查询用户是否存在
	if err := model.DB.First(&user, userID).Error; err != nil {
		return serializer.Response{
			Status: 404,
			Msg:    "用户不存在",
		}
	}

	// 设置新密码
	user.SetPassWord("123456")

	// 保存用户的更新信息
	if err := model.DB.Save(&user).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "更新密码失败",
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "密码重置成功",
	}
}
