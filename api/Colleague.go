package api

import (
	"todo_list/package/utils"
	"todo_list/serializer"
	"todo_list/service"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

func CreateColleague(c *gin.Context) {
	var service service.CreateColleagueService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	if err := c.ShouldBind(&service); err == nil {
		res := service.Create(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, serializer.Response{
			Status: 400,
			Msg:    "参数绑定失败",
		})
	}
}

func ListTimingTask(c *gin.Context) {
	var service service.ListColleagueService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	if err := c.ShouldBindQuery(&service); err == nil {
		res := service.List(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

func UpdateColleague(c *gin.Context) {
	var service service.UpdateColleagueService
	id := c.Param("id")

	if err := c.ShouldBind(&service); err == nil {
		res := service.Update(id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, serializer.Response{
			Status: 400,
			Msg:    "参数绑定失败",
		})
	}
}

func BatchDeleteColleagues(c *gin.Context) {
	var service service.BatchDeleteColleagueService
	if err := c.ShouldBindJSON(&service); err == nil {
		res := service.BatchDelete()
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, serializer.Response{
			Status: 400,
			Msg:    "参数错误",
		})
	}
}

func GetUnsetMapColleagues(c *gin.Context) {
	var listColleague service.ListColleagueService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&listColleague); err == nil {
		res := listColleague.GetUnsetMapColleagues(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}
