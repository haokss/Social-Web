package service

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"todo_list/cache"
	"todo_list/model"
	"todo_list/package/report"
	"todo_list/serializer"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ImportBillService struct {
	Type string                `form:"type" binding:"required"` // "wechat" 或 "alipay"
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type ListBillService struct {
	Name            string  `form:"name"`            // 交易对方（模糊搜索）
	TransactionType string  `form:"transactionType"` // 类型筛选
	StartDate       string  `form:"startDate"`       // 开始日期
	EndDate         string  `form:"endDate"`         // 结束日期
	MinAmount       float64 `form:"minAmount"`
	MaxAmount       float64 `form:"maxAmount"`
	Page            int     `form:"page" binding:"required"`
	PageSize        int     `form:"pageSize" binding:"required"`
}

func (service *ImportBillService) ImportBill(c *gin.Context, id uint) serializer.Response {
	// 保存上传文件到临时目录
	filename := service.File.Filename
	tempPath := filepath.Join(os.TempDir(), filename)
	if err := c.SaveUploadedFile(service.File, tempPath); err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "上传文件失败: " + err.Error(),
		}
	}

	defer os.Remove(tempPath)

	var bills []model.Bill
	var err error

	switch service.Type {
	case "wechat":
		bills, err = report.ParseWeChatXLSX(tempPath)
	case "alipay":
		bills, err = report.ParseAlipayCSV(tempPath)
	default:
		return serializer.Response{Status: 400, Msg: "账单类型错误"}
	}

	if err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "解析账单失败: " + err.Error(),
		}
	}

	// 开启事务（添加调试日志）
	fmt.Println("开始事务，准备插入", len(bills), "条记录")
	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("事务回滚:", r)
		}
	}()

	// 插入
	successCount := 0
	for i, bill := range bills {
		bill.Uid = id
		if err := tx.Create(&bill).Error; err != nil {
			tx.Rollback()
			return serializer.Response{
				Status: 500,
				Msg:    fmt.Sprintf("第%d条记录插入失败: %v\n数据: %+v", i+1, err, bill),
			}
		}
		successCount++
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "事务提交失败: " + err.Error(),
		}
	}

	// var userBills []model.Bill
	// model.DB.Where("uid = ?", id).Find(&userBills)

	// 更新缓存
	if len(bills) == 0 {
		log.Warn("查询到 0 条账单记录，不更新缓存")
	} else {
		cache.AddUserBills(id, bills)
	}

	return serializer.Response{
		Status: 200,
		Msg:    fmt.Sprintf("成功导入 %d/%d 条记录", successCount, len(bills)),
	}
}

func (service *ListBillService) GetBillList(id uint) serializer.Response {
	// 获取当前用户的账单（全部）
	bills := cache.GetUserBills(id)

	// 筛选后的结果
	var filtered []model.Bill

	for _, bill := range bills {
		if service.Name != "" && !strings.Contains(bill.Counterparty, service.Name) {
			continue
		}
		if service.TransactionType != "" && bill.TransactionType != service.TransactionType {
			continue
		}
		if service.StartDate != "" && bill.TransactionTime < service.StartDate {
			continue
		}
		if service.EndDate != "" && bill.TransactionTime > service.EndDate {
			continue
		}
		if service.MinAmount > 0 && bill.Amount < service.MinAmount {
			continue
		}
		if service.MaxAmount > 0 && bill.Amount > service.MaxAmount {
			continue
		}
		filtered = append(filtered, bill)
	}

	total := len(filtered)

	// 按交易时间倒序排序
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].TransactionTime > filtered[j].TransactionTime
	})

	// 分页
	start := (service.Page - 1) * service.PageSize
	end := start + service.PageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	paged := filtered[start:end]

	return serializer.Response{
		Status: 200,
		Msg:    "获取成功",
		Data: map[string]interface{}{
			"list":  paged,
			"total": total,
		},
	}
}
