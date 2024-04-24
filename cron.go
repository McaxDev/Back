package main

import (
	"time"

	co "github.com/McaxDev/Back/config"
	cr "github.com/McaxDev/Back/cron"
	h "github.com/McaxDev/Back/handler"
	"github.com/McaxDev/Back/util"
	ut "github.com/McaxDev/Back/util"
	"github.com/robfig/cron/v3"
)

// 每天执行定时任务的协程
func Cron() {

	// 创建一个新的cron实例
	c := cron.New(cron.WithSeconds())

	// 每小时一次的定时任务
	if _, err := c.AddFunc("0 0 * * * *", func() {

		// 每小时给互通服玩家发金币
		if mes, err := ut.Rcon("main", "money give @a 1000"); err != nil {
			co.SysLog("ERROR", err.Error())
		} else {
			co.SysLog("INFO", mes)
		}

		// 缓存玩家的UUID
		if err := util.ForEach(cr.CachePlayerUUID, "main", "sc"); err != nil {
			co.SysLog("ERROR", "无法缓存玩家的UUID"+err.Error())
		}

		// 缓存玩家游戏数据
		if err := util.ForEach(cr.CacheData, "main", "sc"); err != nil {
			co.SysLog("ERROR", "无法缓存玩家的游戏数据"+err.Error())
		}

	}); err != nil {
		co.SysLog("ERROR", "发金币的定时任务添加失败"+err.Error())
		return
	}

	// 每天执行一次的定时任务
	if _, err := c.AddFunc("0 0 0 * * *", func() {
		err := co.DB.Model(&co.AxolotlCoin{}).Update("Pearl", 10).Error
		if err != nil {
			co.SysLog("ERROR", err.Error())
		} else {
			co.SysLog("INFO", "已将所有玩家的PearlCoin重置为50")
		}
	}); err != nil {
		co.SysLog("ERROR", "每日重置PearlCoin的定时任务添加失败")
		return
	}

	// 每一分钟执行一次的定时任务
	if _, err := c.AddFunc("0 */10 * * * *", func() {

		// 清理内存
		ut.ClearExpDefault(h.Challenges, h.IpTimeMap)
		ut.ClearExpired(func(s h.MailStruct) time.Time { return s.Expiry }, h.Mailsent)

	}); err != nil {
		co.SysLog("ERROR", "添加自动清理内存的定时任务失败")
		return
	}

	// 启动cron实例
	c.Start()
}
