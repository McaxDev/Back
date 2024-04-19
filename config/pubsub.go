package config

import "sync"

// 表示发布订阅模式的结构体
type CoinPubSubTemplate struct {
	Mu   sync.RWMutex
	User map[uint]chan int
}

// 创建上面结构体的结构体对象
var CoinPubSub = &CoinPubSubTemplate{User: make(map[uint]chan int)}

// 为结构体注册订阅余额的方法
func (ps *CoinPubSubTemplate) Sub(userID uint) (<-chan int, error) {

	ps.Mu.Lock()
	defer ps.Mu.Unlock()

	// 检查这个用户是否存在新的余额通道，不存在就创建
	ch, found := ps.User[userID]
	if !found {

		// 创建空通道
		ch = make(chan int, 1)
		ps.User[userID] = ch

		// 向通道里加入用户的余额作为初始值
		var coin AxolotlCoin
		if err := DB.First(&coin, "user_id = ?", userID).Error; err != nil {
			return nil, err
		}
		ch <- coin.Azure + coin.Pearl
	}
	return ch, nil
}

// 为结构体注册发布余额的方法
func (ps *CoinPubSubTemplate) Pub(userID uint, balance int) {

	ps.Mu.RLock()
	defer ps.Mu.RUnlock()

	// 检查通道是否存在并发送余额
	if ch, found := ps.User[userID]; found {
		select {
		case ch <- balance:
		default:
			<-ch
			ch <- balance
		}
	}
}

// 取消订阅一个通道的函数
func (ps *CoinPubSubTemplate) UnSub(userID uint) {

	ps.Mu.Lock()
	defer ps.Mu.Unlock()

	// 检查通道是否存在并删除
	if ch, found := CoinPubSub.User[userID]; found {
		close(ch)
		delete(CoinPubSub.User, userID)
	}
}
