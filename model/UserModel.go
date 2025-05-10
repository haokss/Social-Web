package model

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserName       string `gorm:"unique"`
	PasswordDigest string
	Phone          string
	Email          string
	AvaUrl         string
	Role           int `gorm:"default:1"`
	gorm.Model
}

// 用户基础信息
type UserBaseInfo struct {
	UserName string
	Phone    string
	Email    string
	AvaUrl   string
	gorm.Model
}

func (UserBaseInfo) TableName() string {
	return "user"
}

// 加密
func (user *User) SetPassWord(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// 验证
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}
