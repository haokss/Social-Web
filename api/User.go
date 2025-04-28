package api

// 接口注册层
import (
	"todo_list/package/utils"
	"todo_list/service"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

// 用户注册接口注册
func UserRegister(c *gin.Context) {
	var userRegister service.UserService
	if err := c.ShouldBind(&userRegister); err == nil {
		res := userRegister.Register()
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 用户登录接口注册
func UserLogin(c *gin.Context) {
	var userLogin service.UserService
	if err := c.ShouldBind(&userLogin); err == nil {
		res := userLogin.Login()
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 用户基础信息
func UserInfo(c *gin.Context) {
	var user_info service.UserBaseInfoService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&user_info); err == nil {
		res := user_info.UserInfo(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

func UpdateUserInfo(c *gin.Context) {
	var user_info service.UserBaseInfoService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	// if err := c.ShouldBind(&user_info); err == nil {
	res := user_info.UpdateUserInfo(claim.Id)
	c.JSON(200, res)
	// } else {
	// 	logging.Error(err)
	// 	c.JSON(400, err)
	// }
}
