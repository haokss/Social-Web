package api

import (
	"todo_list/package/utils"
	"todo_list/service"

	"github.com/gin-gonic/gin"
)

func AdminListUsers(c *gin.Context) {
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if claim.Role != "admin" {
		c.JSON(403, gin.H{"status": 403, "msg": "无权限"})
		return
	}

	var service service.AdminUserListService
	if err := c.ShouldBindQuery(&service); err == nil {
		res := service.ListUsers()
		c.JSON(200, res)
	} else {
		c.JSON(400, gin.H{"status": 400, "msg": "参数绑定失败"})
	}
}

func AdminEditUser(c *gin.Context) {
	// 解析token并验证权限
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if claim.Role != "admin" {
		c.JSON(403, gin.H{"status": 403, "msg": "无权限"})
		return
	}

	// 获取用户ID
	userID := c.Param("id")

	// 绑定请求体到结构体
	var service service.AdminEditUserService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(400, gin.H{"status": 400, "msg": "参数绑定失败"})
		return
	}

	// 调用服务层处理编辑用户逻辑
	res := service.EditUser(userID)
	c.JSON(res.Status, res)
}

func AdminResetPassword(c *gin.Context) {
	// 解析token并验证权限
	var service service.AdminResetPasswordService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if claim.Role != "admin" {
		c.JSON(403, gin.H{"status": 403, "msg": "无权限"})
		return
	}

	// 获取用户ID
	userID := c.Param("id")

	// 调用重置密码服务
	res := service.ResetPassword(userID)
	c.JSON(res.Status, res)
}
