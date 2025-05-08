package cache

import (
	"sync"
	"todo_list/model"
)

var (
	BillCache     = make(map[uint][]model.Bill) // key: user ID → 用户账单列表
	BillCacheLock sync.RWMutex
)

// InitBillCache 在服务启动时调用，加载所有账单到内存（如果需要全量加载）
func InitBillCache() error {
	var bills []model.Bill
	if err := model.DB.Find(&bills).Error; err != nil {
		return err
	}
	BillCacheLock.Lock()
	defer BillCacheLock.Unlock()
	for _, bill := range bills {
		uid := bill.Uid // 假设你有 CreatedBy 字段或类似的用户 ID
		BillCache[uid] = append(BillCache[uid], bill)
	}
	return nil
}

// GetUserBills 从缓存获取指定用户的账单列表
func GetUserBills(uid uint) []model.Bill {
	BillCacheLock.RLock()
	defer BillCacheLock.RUnlock()
	return BillCache[uid]
}

// SetUserBills 更新指定用户的账单列表到缓存
func SetUserBills(uid uint, bills []model.Bill) {
	BillCacheLock.Lock()
	defer BillCacheLock.Unlock()
	BillCache[uid] = bills
}

// AddBill 把一条账单加进缓存
func AddBill(uid uint, bill model.Bill) {
	BillCacheLock.Lock()
	defer BillCacheLock.Unlock()
	BillCache[uid] = append(BillCache[uid], bill)
}

// DeleteUserBills 清空某个用户的账单缓存
func DeleteUserBills(uid uint) {
	BillCacheLock.Lock()
	defer BillCacheLock.Unlock()
	delete(BillCache, uid)
}
