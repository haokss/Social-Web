package cache

import (
	"sync"
	"todo_list/model"

	log "github.com/sirupsen/logrus"
)

var (
	BillCache     = make(map[uint][]model.Bill) // key: 用户ID, value: 该用户的账单列表
	BillCacheLock sync.RWMutex
)

// InitBillCache 启动时加载所有账单
func InitBillCache() error {
	var bills []model.Bill
	if err := model.DB.Find(&bills).Error; err != nil {
		return err
	}
	BillCacheLock.Lock()
	defer BillCacheLock.Unlock()
	for _, bill := range bills {
		BillCache[bill.Uid] = append(BillCache[bill.Uid], bill)
	}
	log.Info(len(bills))
	return nil
}

// UpdateUserBills 更新某个用户的账单列表
func UpdateUserBills(uid uint, bills []model.Bill) {
	BillCacheLock.Lock()
	defer BillCacheLock.Unlock()
	BillCache[uid] = bills
}

// GetUserBills 获取某个用户的账单列表
func GetUserBills(uid uint) []model.Bill {
	BillCacheLock.RLock()
	defer BillCacheLock.RUnlock()
	return BillCache[uid]
}

func AddUserBills(uid uint, newBills []model.Bill) {
	BillCacheLock.Lock()
	defer BillCacheLock.Unlock()

	// 取出原有缓存
	existing, ok := BillCache[uid]

	// 构造一个 map 存已存在的 TransactionID，加速查重
	existingIDs := make(map[string]struct{})
	if ok {
		for _, bill := range existing {
			existingIDs[bill.TransactionID] = struct{}{}
		}
	}

	// 过滤出不重复的新账单
	var uniqueNewBills []model.Bill
	for _, bill := range newBills {
		if _, found := existingIDs[bill.TransactionID]; !found {
			uniqueNewBills = append(uniqueNewBills, bill)
			// 加入 map，避免后续重复（如果 newBills 本身有重复）
			existingIDs[bill.TransactionID] = struct{}{}
		}
	}

	if ok {
		// 合并唯一新账单
		BillCache[uid] = append(existing, uniqueNewBills...)
	} else {
		// 没有旧缓存，直接用唯一新账单
		BillCache[uid] = uniqueNewBills
	}
}
