package serializer

type TaskResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Status   int    `json:"status"`
	EndTime  int64  `json:"end_time"` // 前端用时间戳
	Priority int    `json:"priority"`
	Type     int    `json:"type"`
}
