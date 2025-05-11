package api

import (
	"todo_list/package/utils"
	"todo_list/service"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

func ListNotification(c *gin.Context) {
	var list service.NotificationListService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	if err := c.ShouldBind(&list); err == nil {
		res := list.List(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

// 标记单个已读
func MarkNotificationRead(c *gin.Context) {
	var service service.NotificationListService

	id := c.Param("id")                                        // 获取通知ID
	claim, _ := utils.ParseToken(c.GetHeader("Authorization")) // 获取用户ID

	// 调用服务方法标记为已读
	res := service.MarkNotificationRead(claim.Id, id)
	c.JSON(res.Status, res)
}

func MarkAllNotificationsRead(c *gin.Context) {
	var service service.NotificationListService

	claim, _ := utils.ParseToken(c.GetHeader("Authorization")) // 获取用户ID

	// 调用服务方法标记所有未读消息为已读
	res := service.MarkAllNotificationsRead(claim.Id)
	c.JSON(res.Status, res)
}
