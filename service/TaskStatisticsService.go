package service

import (
	"sort"
	"time"
	"todo_list/cache"
	"todo_list/model"
	"todo_list/serializer"
)

type StatisticsService struct{}

type TrendService struct{}

type TypeDistributionService struct{}

type UpcomingService struct{}

type HighPriorityService struct{}

// /stats: 统计 Status 数量和时间范围
func (s *StatisticsService) Stats(uid uint) serializer.Response {
	tasks := cache.GetUserTasks(uid)
	var total, completed, inProgress, overdue int
	now := time.Now().Unix()

	// 当前时间点的统计
	for _, t := range tasks {
		total++
		if t.Status == 1 {
			completed++
		} else {
			if t.EndTime < now {
				overdue++
			} else {
				inProgress++
			}
		}
	}

	// 计算完成率
	var completionRate float64
	if total > 0 {
		completionRate = float64(completed) / float64(total) * 100
	}

	// 环比变化（这里简单用数量比，不考虑百分比增减的方向）
	lastMonthCount := 0
	lastWeekInProgress := 0
	lastWeekOverdue := 0

	monthAgo := now - 30*24*3600
	weekAgo := now - 7*24*3600

	for _, t := range tasks {
		if t.CreatedAt.Unix() >= monthAgo && t.CreatedAt.Unix() < now {
			lastMonthCount++
		}
		if t.UpdatedAt.Unix() >= weekAgo && t.UpdatedAt.Unix() < now {
			if t.Status == 0 && t.EndTime >= now {
				lastWeekInProgress++
			}
			if t.Status == 0 && t.EndTime < now {
				lastWeekOverdue++
			}
		}
	}

	// 避免除零
	taskChange := 0.0
	if lastMonthCount > 0 {
		taskChange = (float64(total - lastMonthCount)) / float64(lastMonthCount) * 100
	}

	progressChange := 0.0
	if lastWeekInProgress > 0 {
		progressChange = (float64(inProgress - lastWeekInProgress)) / float64(lastWeekInProgress) * 100
	}

	overdueChange := 0.0
	if lastWeekOverdue > 0 {
		overdueChange = (float64(overdue - lastWeekOverdue)) / float64(lastWeekOverdue) * 100
	}

	data := map[string]interface{}{
		"totalTasks":      total,
		"completedTasks":  completed,
		"inProgressTasks": inProgress,
		"overdueTasks":    overdue,
		"completionRate":  int(completionRate),
		"taskChange":      int(taskChange),
		"progressChange":  int(progressChange),
		"overdueChange":   int(overdueChange),
	}

	return serializer.Response{Status: 200, Data: data}
}

func (s *TrendService) Trend(uid uint, rangeParam string) serializer.Response {
	tasks := cache.GetUserTasks(uid)
	var periodDays int64
	switch rangeParam {
	case "7d":
		periodDays = 7
	case "30d":
		periodDays = 30
	case "90d":
		periodDays = 90
	default:
		periodDays = 7 // 默认 7 天
	}

	now := time.Now()
	dateCounts := make(map[string]map[string]int) // date -> {"created": x, "completed": y}

	// 初始化每天的计数
	for i := int64(0); i < periodDays; i++ {
		date := now.AddDate(0, 0, -int(i)).Format("2006-01-02")
		dateCounts[date] = map[string]int{"created": 0, "completed": 0}
	}

	for _, t := range tasks {
		createdDate := t.CreatedAt.Format("2006-01-02")
		if _, ok := dateCounts[createdDate]; ok {
			dateCounts[createdDate]["created"]++
		}
		if t.Status == 1 {
			completedDate := t.UpdatedAt.Format("2006-01-02")
			if _, ok := dateCounts[completedDate]; ok {
				dateCounts[completedDate]["completed"]++
			}
		}
	}

	// 构造返回数组（按日期升序）
	var dates []string
	var createdCounts []int
	var completedCounts []int

	for i := int64(periodDays - 1); i >= 0; i-- {
		date := now.AddDate(0, 0, -int(i)).Format("2006-01-02")
		dates = append(dates, date)
		createdCounts = append(createdCounts, dateCounts[date]["created"])
		completedCounts = append(completedCounts, dateCounts[date]["completed"])
	}

	data := map[string]interface{}{
		"dates":     dates,
		"created":   createdCounts,
		"completed": completedCounts,
	}

	return serializer.Response{
		Status: 200,
		Data:   data,
	}
}

func (s *StatisticsService) TypeDistribution(uid uint) serializer.Response {
	var tasks []model.TimingTask

	for taskID := range cache.TimingTaskCache {
		// 使用 GetTimingTask 获取单个任务
		if task, exists := cache.GetTimingTask(taskID); exists {
			// 检查任务是否属于当前用户
			if task.Uid == uid {
				tasks = append(tasks, task)
			}
		}
	}

	// 根据任务的类型进行统计
	typeCounts := make(map[int]int)
	for _, t := range tasks {
		typeCounts[t.Type]++
	}

	var types []int
	for t := range typeCounts {
		types = append(types, t)
	}
	sort.Ints(types)

	var counts []int
	for _, t := range types {
		counts = append(counts, typeCounts[t])
	}

	// 返回数据
	data := map[string]interface{}{
		"types":  types,
		"counts": counts,
	}

	return serializer.Response{Status: 200, Data: data}
}
