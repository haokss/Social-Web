package serializer

import "todo_list/model"

type User struct {
	ID       uint   `json:"id" form:"id" example:"1"`
	UserName string `json:"user_name" form:"user_name" example:"haoks"`
	CreateAt int64  `json:"create_at" form:"create_at"`
}

// 用户基础信息
type UserBaseInfo struct {
	ID       uint   `json:"id" form:"id" example:"1"`
	UserName string `json:"user_name" form:"user_name"`
	Phone    string `json:"phone" form:"phone"`
	Email    string `json:"email" form:"email"`
	AvaUrl   string `json:"ava_url" form:"ava_url"`
}

// 序列化返回数据
func BuildUser(user model.User) User {
	return User{
		ID:       user.ID,
		UserName: user.UserName,
		CreateAt: user.CreatedAt.Unix(),
	}
}

func BuildUserBaseInfo(user model.UserBaseInfo) UserBaseInfo {
	return UserBaseInfo{
		ID:       user.ID,
		UserName: user.UserName,
		Phone:    user.Phone,
		Email:    user.Email,
		AvaUrl:   user.AvaUrl,
	}

}
