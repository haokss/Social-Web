package api

import (
	"todo_list/model"
	"todo_list/package/utils"
	"todo_list/serializer"
	"todo_list/service"

	"github.com/gin-gonic/gin"
)

func CreatePoint(c *gin.Context) {
	var createPointService service.CreatePointService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	// 绑定请求体数据到 CreatePointService
	if err := c.ShouldBind(&createPointService); err == nil {
		// 调用服务层的 Create 方法，传递用户 ID 和请求数据
		res := createPointService.Create(claim.Id)
		c.JSON(200, res)
	} else {
		// 处理请求绑定失败的情况
		c.JSON(400, serializer.Response{
			Status: 400,
			Msg:    err.Error(),
		})
	}
}

func ListPoints(c *gin.Context) {
	var listPointService service.ListPointService
	claim, err := utils.ParseToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(401, serializer.Response{
			Status: 401,
			Msg:    "invalid token",
		})
		return
	}

	res := listPointService.List(claim.Id)
	c.JSON(200, res)
}

// UpdatePoint 修改点位
func UpdatePoint(c *gin.Context) {
	var listPointService service.ListPointService

	// 获取用户信息
	claim, err := utils.ParseToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(401, serializer.Response{
			Status: 401,
			Msg:    "invalid token",
		})
		return
	}

	// 解析请求体中的点位数据
	var point model.Point
	if err := c.ShouldBindJSON(&point); err != nil {
		c.JSON(400, serializer.Response{
			Status: 400,
			Msg:    "请求参数错误",
		})
		return
	}

	// 更新点位数据
	res := listPointService.Update(point, claim.Id)
	c.JSON(200, res)
}

// DeletePoint 删除点位
func DeletePoint(c *gin.Context) {
	var listPointService service.ListPointService

	// 获取用户信息
	claim, err := utils.ParseToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(401, serializer.Response{
			Status: 401,
			Msg:    "invalid token",
		})
		return
	}

	// 获取点位 ID
	pointID := c.Param("id")
	if pointID == "" {
		c.JSON(400, serializer.Response{
			Status: 400,
			Msg:    "点位ID不能为空",
		})
		return
	}

	// 删除点位
	res := listPointService.Delete(pointID, claim.Id)
	c.JSON(200, res)
}
