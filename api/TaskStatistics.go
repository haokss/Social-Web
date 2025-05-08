package api

import (
	"todo_list/package/utils"
	"todo_list/service"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

func GetTaskStats(c *gin.Context) {
	var service service.StatisticsService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.Stats(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

func GetTaskTrend(c *gin.Context) {
	var service service.TrendService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.Trend(claim.Id, c.Query("range"))
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

func GetTaskTypeDistribution(c *gin.Context) {
	var service service.StatisticsService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	res := service.TypeDistribution(claim.Id)
	c.JSON(200, res)
}
