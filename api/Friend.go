package api

import (
	"todo_list/package/utils"
	"todo_list/serializer"
	"todo_list/service"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

func CreateFriend(c *gin.Context) {
	var service service.CreateFriendService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	if err := c.ShouldBind(&service); err == nil {
		res := service.Create(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, serializer.Response{Status: 400, Msg: "参数绑定失败"})
	}
}

func ListFriend(c *gin.Context) {
	var service service.ListFriendService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	if err := c.ShouldBindQuery(&service); err == nil {
		res := service.List(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

func UpdateFriend(c *gin.Context) {
	var service service.UpdateFriendService
	id := c.Param("id")

	if err := c.ShouldBind(&service); err == nil {
		res := service.Update(id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, serializer.Response{Status: 400, Msg: "参数绑定失败"})
	}
}

func BatchDeleteFriends(c *gin.Context) {
	var service service.BatchDeleteFriendService

	if err := c.ShouldBindJSON(&service); err == nil {
		res := service.BatchDelete()
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, serializer.Response{Status: 400, Msg: "参数错误"})
	}
}

func GetUnsetMapFriends(c *gin.Context) {
	var listFriend service.ListFriendService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&listFriend); err == nil {
		res := listFriend.GetUnsetMapFriends(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}
