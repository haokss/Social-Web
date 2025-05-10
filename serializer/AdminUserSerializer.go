package serializer

import "todo_list/model"

// 单个用户格式化
func BuildAdminUser(user model.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		UserName: user.UserName,
		Phone:    user.Phone,
		Email:    user.Email,
		AvaUrl:   user.AvaUrl,
		Role:     user.Role,
	}
}

// 批量用户格式化
func BuildAdminUsers(users []model.User) []UserResponse {
	var res []UserResponse
	for _, user := range users {
		res = append(res, BuildAdminUser(user))
	}
	return res
}

// 前端需要的返回结构
type UserResponse struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	AvaUrl   string `json:"ava_url"`
	Role     int    `json:"role"`
}
