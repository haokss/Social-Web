package api

import (
	"todo_list/package/utils"
	"todo_list/service"

	"github.com/gin-gonic/gin"
)

func AdminGetAllTimingTasks(c *gin.Context) {
	// 解析token并验证权限
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if claim.Role != "admin" {
		c.JSON(403, gin.H{"status": 403, "msg": "无权限"})
		return
	}

	// 创建 service 实例，获取所有定时任务
	var service service.AdminTimingTaskListService
	res := service.GetAllTimingTasks()

	// 返回响应
	c.JSON(res.Status, res)
}

func AdminBatchAuditTimingTask(c *gin.Context) {
	// 解析 token
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if claim.Role != "admin" {
		c.JSON(403, gin.H{"status": 403, "msg": "无权限"})
		return
	}

	var service service.AdminBatchAuditTimingTaskService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(400, gin.H{"status": 400, "msg": "请求参数错误"})
		return
	}

	res := service.BatchAuditTimingTask(c)
	c.JSON(res.Status, res)
}
