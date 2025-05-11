package api

import (
	"strconv"
	"todo_list/package/utils"
	"todo_list/serializer"
	"todo_list/service"

	"github.com/gin-gonic/gin"
)

func AdminGetAllImages(c *gin.Context) {
	var uploadService service.UploadService

	// 解析 token
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if claim.Role != "admin" {
		c.JSON(403, serializer.Response{
			Status: 403,
			Msg:    "无权限，只有管理员可以查看所有图片",
		})
		return
	}

	// 分页参数（默认 page=1, page_size=10）
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword") // 可选：搜索关键字（模糊匹配名称或描述）

	res := uploadService.GetAllImages(page, pageSize, keyword)
	c.JSON(200, res)
}

func AdminBatchAuditImages(c *gin.Context) {
	// 解析 token
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if claim.Role != "admin" {
		c.JSON(403, serializer.Response{
			Status: 403,
			Msg:    "无权限，只有管理员可以审核图片",
		})
		return
	}

	var service service.AdminBatchAuditImageService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(400, serializer.Response{
			Status: 400,
			Msg:    "请求参数错误",
		})
		return
	}

	res := service.BatchAudit()
	c.JSON(res.Status, res)
}
