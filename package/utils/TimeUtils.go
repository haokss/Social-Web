package utils

import "time"

func parseTime(timeStr string) time.Time {
	t, err := time.Parse("2006-01-02 15:04", timeStr) // 根据实际格式调整
	if err != nil {
		return time.Now() // 或处理为默认值
	}
	return t
}
