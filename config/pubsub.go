package config

import (
	"sync"
)

// 表示发布订阅模式的结构体
type CoinPubSubTemplate struct {
	Mu   sync.RWMutex
	User map[uint][]*CoinType
}

// 表示两种币的结构体
type CoinType struct {
	Pearl chan int
	Azure chan int
}

// 创建上面结构体的结构体对象
var CoinPubSub = &CoinPubSubTemplate{User: make(map[uint][]*CoinType)}

// 删除一个正在订阅的设备
func (ps *CoinPubSubTemplate) RemoveDevice(userID uint, device *CoinType) {

	ps.Mu.Lock()
	ps.Mu.Unlock()

	// 查找并删除指定的设备
	devices := ps.User[userID]
	for key, value := range devices {
		if value == device {
			devices = append(devices[:key], devices[key+1:]...)
			break
		}
	}
}

// 为结构体注册发布余额的方法
func (ps *CoinPubSubTemplate) Pub(userID uint, azure, pearl int) {
	ps.Mu.RLock()
	defer ps.Mu.RUnlock()

	// 检查通道是否存在并发送余额
	for _, device := range ps.User[userID] {
		select {
		case device.Azure <- azure:
		default:
			<-device.Azure
			device.Azure <- azure
		}
		select {
		case device.Pearl <- pearl:
		default:
			<-device.Pearl
			device.Pearl <- pearl
		}
	}
}
