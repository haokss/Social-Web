package sse

import (
	"container/heap"
	"encoding/json"
	"sync"
	"time"
	"todo_list/model"
)

// Message 消息结构
type Message struct {
	Event      string
	Data       interface{}
	TargetIDs  []uint // 空切片表示广播
	ScheduleID string // 可选：用于取消定时任务
}

// Client 客户端连接
type Client struct {
	UserID  uint
	Message chan Message
}

// delayedMessage 延迟消息结构
type delayedMessage struct {
	msg       Message
	timestamp time.Time // 绝对时间
	index     int       // 堆索引
}

// priorityQueue 优先级队列（最小堆）
type priorityQueue []*delayedMessage

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].timestamp.Before(pq[j].timestamp) }
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*delayedMessage)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // 避免内存泄漏
	item.index = -1 // 标记已移除
	*pq = old[0 : n-1]
	return item
}

// Broker 消息代理中心
type Broker struct {
	clients    map[uint][]*Client // 用户ID到客户端的映射
	pq         *priorityQueue     // 定时任务队列
	muClients  sync.RWMutex       // 客户端映射锁
	muQueue    sync.Mutex         // 队列操作锁
	queueCond  *sync.Cond         // 队列条件变量
	shutdownCh chan struct{}      // 关闭通道
}

func NewBroker() *Broker {
	pq := make(priorityQueue, 0)
	b := &Broker{
		clients:    make(map[uint][]*Client),
		pq:         &pq,
		shutdownCh: make(chan struct{}),
	}
	b.queueCond = sync.NewCond(&b.muQueue)
	heap.Init(b.pq)
	go b.scheduleWorker()
	return b
}

// 定时任务工作协程
func (b *Broker) scheduleWorker() {
	for {
		select {
		case <-b.shutdownCh:
			return
		default:
			b.muQueue.Lock()

			// 等待队列中有任务
			for b.pq.Len() == 0 {
				b.queueCond.Wait()
			}

			// 获取最近的任务
			next := (*b.pq)[0]
			now := time.Now()

			if next.timestamp.After(now) {
				// 计算等待时间
				waitDuration := next.timestamp.Sub(now)
				b.muQueue.Unlock()

				// 超时等待
				select {
				case <-b.shutdownCh: // 在等待过程中检查是否关闭
					return
				case <-time.After(waitDuration):
				}

				// 重新进入循环，确保任务到期
				continue
			}

			// 取出任务并处理
			item := heap.Pop(b.pq).(*delayedMessage)
			b.muQueue.Unlock()

			b.Notify(item.msg)
		}
	}
}

// Register 注册客户端
func (b *Broker) Register(client *Client) {
	b.muClients.Lock()
	defer b.muClients.Unlock()
	b.clients[client.UserID] = append(b.clients[client.UserID], client)
}

// Unregister 注销客户端
func (b *Broker) Unregister(userID uint, client *Client) {
	b.muClients.Lock()
	defer b.muClients.Unlock()

	clients := b.clients[userID]
	for i, c := range clients {
		if c == client {
			b.clients[userID] = append(clients[:i], clients[i+1:]...)
			close(c.Message) // 关闭通道
			break
		}
	}
}

// 立即通知接口
func (b *Broker) Notify(msg Message) {
	b.muClients.RLock()
	defer b.muClients.RUnlock()

	saveToDB := func(userID uint) {
		dataBytes, err := json.Marshal(msg.Data)
		if err != nil {
			return
		}
		notification := model.Notification{
			UserID:     userID,
			Event:      msg.Event,
			Data:       string(dataBytes),
			ScheduleID: msg.ScheduleID,
		}
		_ = model.DB.Create(&notification).Error // 可加日志
	}

	send := func(userID uint, client *Client) {
		select {
		case client.Message <- msg:
		default:
			// 忽略失败的发送
		}
	}

	if len(msg.TargetIDs) == 0 {
		// 广播：所有用户都保存 + 推送
		processed := make(map[uint]bool)
		for userID, clients := range b.clients {
			if !processed[userID] {
				saveToDB(userID)
				processed[userID] = true
			}
			for _, client := range clients {
				send(userID, client)
			}
		}
	} else {
		// 先保存所有目标用户的通知记录
		for _, userID := range msg.TargetIDs {
			saveToDB(userID)
		}
		// 然后给在线用户发送
		for _, userID := range msg.TargetIDs {
			if clients, exists := b.clients[userID]; exists {
				for _, client := range clients {
					send(userID, client)
				}
			}
		}
	}
}

// 定时通知接口
func (b *Broker) ScheduleNotify(msg Message, notifyTime time.Time) {
	b.muQueue.Lock()
	defer b.muQueue.Unlock()

	heap.Push(b.pq, &delayedMessage{
		msg:       msg,
		timestamp: notifyTime.UTC(), // 使用UTC时间避免时区问题
	})
	b.queueCond.Signal() // 唤醒工作协程
}

// Shutdown 优雅关闭
func (b *Broker) Shutdown() {
	close(b.shutdownCh)
}

// 发送通知 实例结构
// broker := c.MustGet("sseBroker").(*sse.Broker)
//
//	broker.Notify(sse.Message{
//		Event: "instant_notification",
//		Data:  map[string]interface{}{"title": "提醒", "content": "Hello"},
//	})
//
