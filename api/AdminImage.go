package api

import (
	"todo_list/package/utils"
	"todo_list/serializer"
	"todo_list/service"

	"github.com/gin-gonic/gin"
)

func AdminGetAllImages(c *gin.Context) {
	var uploadService service.UploadService

	// 解析 token
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	// 检查是否是管理员（假设 claim.Role 存储角色）
	if claim.Role != "admin" {
		c.JSON(403, serializer.Response{
			Status: 403,
			Msg:    "无权限，只有管理员可以查看所有图片",
		})
		return
	}

	res := uploadService.GetAllImages()
	c.JSON(200, res)
}
