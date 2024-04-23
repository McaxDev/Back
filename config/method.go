package config

import "errors"

// 修改一个用户的余额，负数扣费，正数加钱
func (user *User) Transact(amount int) error {

	// 开启事务
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查用户币的数量
	left := user.Balance.Azure + user.Balance.Pearl + amount
	if left < 0 {
		tx.Rollback()
		return errors.New("你钱不够")
	}

	if amount < 0 { // 如果是扣费操作，先扣Pearl再扣Azure
		cost := -amount
		if user.Balance.Pearl >= cost {
			user.Balance.Pearl -= cost
		} else {
			cost -= user.Balance.Pearl
			user.Balance.Pearl = 0
			user.Balance.Azure -= cost
		}
	} else { // 如果是加钱操作，只加Azure
		user.Balance.Azure += amount
	}

	// 更新用户的币
	if err := tx.Save(&user.Balance).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	// 发送发布订阅事件
	CoinPubSub.Pub(user.ID, user.Balance.Azure, user.Balance.Pearl)
	return nil
}
