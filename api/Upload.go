package api

import (
	sse "todo_list/package/SSE"
	"todo_list/package/utils"
	"todo_list/service"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

// 上传头像
func UploadAva(c *gin.Context) {
	var upload_ava service.UploadService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	broker := c.MustGet("sseBroker").(*sse.Broker)
	broker.Notify(sse.Message{
		Event: "instant_notification",
		Data:  map[string]interface{}{"title": "提醒", "content": "Hello"},
	})

	if err := c.ShouldBind(&upload_ava); err == nil {
		res := upload_ava.UploadAva(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

// 上传图片
func UploadImage(c *gin.Context) {
	var upload_ava service.UploadService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&upload_ava); err == nil {
		res := upload_ava.UploadImage(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

// 获取图片
func GetImage(c *gin.Context) {
	var upload_ava service.UploadService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	// if err := c.ShouldBind(&upload_ava); err == nil {
	res := upload_ava.GetImage(claim.Id)
	c.JSON(200, res)
	// } else {
	// 	logging.Error(err)
	// 	c.JSON(400, err)
	// }
}
