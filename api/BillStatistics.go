package api

import (
	"todo_list/package/utils"
	"todo_list/service"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

func GetBillStats(c *gin.Context) {
	var service service.BillStatisticsService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.Stats(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}
