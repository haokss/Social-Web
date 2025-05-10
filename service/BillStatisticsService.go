package service

import (
	"sort"
	"time"
	"todo_list/cache"
	"todo_list/serializer"
)

// StatisticsService 提供统计服务
type BillStatisticsService struct {
	Range string `form:"range"` // 7d, 30d, 90d
}

func (s *BillStatisticsService) Stats(uid uint) serializer.Response {
	bills := cache.GetUserBills(uid)
	now := time.Now()
	var days int

	// 根据时间范围设置天数
	switch s.Range {
	case "7d":
		days = 7
	case "90d":
		days = 90
	default:
		days = 30
	}

	startTime := now.AddDate(0, 0, -days)

	// 统计数据
	var totalIncome, totalExpense, totalTransactions int
	var incomeMap = make(map[string]float64)
	var expenseMap = make(map[string]float64)

	var trendDates []string
	var trendIncome []float64
	var trendExpense []float64

	// 初始化日期映射
	for i := 0; i < days; i++ {
		day := startTime.AddDate(0, 0, i).Format("2006-01-02")
		trendDates = append(trendDates, day)
		incomeMap[day] = 0
		expenseMap[day] = 0
	}

	// 统计本期和上期
	var prevIncome, prevExpense, prevTransactions float64

	for _, b := range bills {
		createdAt := b.CreatedAt
		if createdAt.After(startTime) && createdAt.Before(now) {
			totalTransactions++
			if b.IncomeExpense == "收入" {
				totalIncome += int(b.Amount)
			} else if b.IncomeExpense == "支出" {
				totalExpense += int(b.Amount)
			}

			dayKey := createdAt.Format("2006-01-02")
			if b.IncomeExpense == "收入" {
				incomeMap[dayKey] += b.Amount
			} else if b.IncomeExpense == "支出" {
				expenseMap[dayKey] += b.Amount
			}
		}

		// 上一周期
		prevStart := startTime.AddDate(0, 0, -days)
		if createdAt.After(prevStart) && createdAt.Before(startTime) {
			if b.IncomeExpense == "收入" {
				prevIncome += b.Amount
			} else if b.IncomeExpense == "支出" {
				prevExpense += b.Amount
			}
			prevTransactions++
		}
	}

	// 汇总趋势数据
	for _, date := range trendDates {
		trendIncome = append(trendIncome, incomeMap[date])
		trendExpense = append(trendExpense, expenseMap[date])
	}

	// 环比变化计算
	incomeChange := percentChange(float64(totalIncome), prevIncome)
	expenseChange := percentChange(float64(totalExpense), prevExpense)
	transactionChange := percentChange(float64(totalTransactions), prevTransactions)
	netIncome := float64(totalIncome) - float64(totalExpense)
	netChange := percentChange(netIncome, prevIncome-prevExpense)

	// 分类统计
	incomeCategories := make(map[string]float64)
	expenseCategories := make(map[string]float64)
	for _, b := range bills {
		if b.CreatedAt.After(startTime) && b.CreatedAt.Before(now) {
			if b.IncomeExpense == "收入" {
				incomeCategories[b.TransactionType] += b.Amount
			} else if b.IncomeExpense == "支出" {
				expenseCategories[b.TransactionType] += b.Amount
			}
		}
	}

	// TOP 10 大额交易
	sort.Slice(bills, func(i, j int) bool {
		return bills[i].Amount > bills[j].Amount
	})
	topTransactions := []serializer.TransactionRank{}
	for _, b := range bills {
		if len(topTransactions) >= 10 {
			break
		}
		topTransactions = append(topTransactions, serializer.TransactionRank{
			TransactionTime: b.CreatedAt.Format("2006-01-02 15:04:05"),
			Counterparty:    b.Counterparty,
			Product:         b.Product,
			IncomeExpense:   b.IncomeExpense,
			Amount:          b.Amount,
			PaymentMethod:   b.PaymentMethod,
		})
	}

	// 返回数据
	return serializer.Response{
		Status: 200,
		Data: map[string]interface{}{
			"stats": map[string]interface{}{
				"totalIncome":       totalIncome,
				"totalExpense":      totalExpense,
				"totalTransactions": totalTransactions,
				"netIncome":         netIncome,
				"incomeChange":      incomeChange,
				"expenseChange":     expenseChange,
				"transactionChange": transactionChange,
				"netChange":         netChange,
			},
			"trend": map[string]interface{}{
				"dates":   trendDates,
				"income":  trendIncome,
				"expense": trendExpense,
			},
			"incomeCategories":  categoryList(incomeCategories),
			"expenseCategories": categoryList(expenseCategories),
			"topTransactions":   topTransactions,
		},
	}
}

func percentChange(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return ((current - previous) / previous) * 100
}

func categoryList(categoryMap map[string]float64) []map[string]interface{} {
	list := []map[string]interface{}{}
	for k, v := range categoryMap {
		list = append(list, map[string]interface{}{
			"category": k,
			"amount":   v,
		})
	}
	return list
}
