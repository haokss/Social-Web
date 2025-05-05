package api

import (
	"todo_list/package/utils"
	"todo_list/service"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

// 创建一个亲属关系
func CreateRelative(c *gin.Context) {
	var createRelative service.CreateRelativeService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&createRelative); err == nil {
		res := createRelative.CreateRelative(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

// 获取所有亲属关系信息
func GetAllRelatives(c *gin.Context) {
	var listRelative service.ListRelativeService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&listRelative); err == nil {
		res := listRelative.GetRelativeList(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

// 更新亲属关系
func UpdateRelative(c *gin.Context) {
	var updateService service.UpdateRelativeService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&updateService); err == nil {
		res := updateService.Update(claim.Id, c.Param("id"))
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

// 删除亲属关系
func DeleteRelative(c *gin.Context) {
	var deleteService service.DeleteRelativeService
	// claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&deleteService); err == nil {
		res := deleteService.DeleteRelative()
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}
