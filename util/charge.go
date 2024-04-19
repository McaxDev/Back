package util

import (
	"errors"

	co "github.com/McaxDev/Back/config"
)

// 对特定的操作进行扣费的函数
func Charge(userID uint, cost int) error {

	// 开启事务
	tx := co.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 将用户的余额映射到结构体对象
	var balance co.AxolotlCoin
	if err := co.DB.First(&balance, "user_id = ?", userID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 检查用户币的数量
	left := balance.Azure + balance.Pearl - cost
	if left < 0 {
		tx.Rollback()
		return errors.New("你钱不够")
	}

	// 执行扣费逻辑
	if balance.Pearl >= cost {
		balance.Pearl -= cost
	} else {
		cost -= balance.Pearl
		balance.Pearl = 0
		balance.Azure -= cost
	}

	// 更新用户的币
	if err := tx.Save(&balance).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	// 发送发布订阅事件
	co.CoinPubSub.Pub(userID, left)
	return nil
}
