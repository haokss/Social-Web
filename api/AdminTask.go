package api

import (
	"todo_list/package/utils"
	"todo_list/service"

	"github.com/gin-gonic/gin"
)

func AdminGetAllUserTasks(c *gin.Context) {
	// 解析token并验证权限
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if claim.Role != "admin" {
		c.JSON(403, gin.H{"status": 403, "msg": "无权限"})
		return
	}

	// 创建 service 实例，获取所有任务
	var service service.AdminTaskListService
	res := service.GetAllUserTasks()

	// 返回响应
	c.JSON(res.Status, res)
}

func AdminBatchAudit(c *gin.Context) {
	// 解析 token
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if claim.Role != "admin" {
		c.JSON(403, gin.H{"status": 403, "msg": "无权限"})
		return
	}

	var service service.AdminBatchAuditService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(400, gin.H{"status": 400, "msg": "请求参数错误"})
		return
	}

	res := service.BatchAudit(c)
	c.JSON(res.Status, res)
}
