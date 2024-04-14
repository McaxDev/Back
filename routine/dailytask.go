package routine

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/robfig/cron/v3"
)

// 每天执行定时任务的协程
func DailyTask() {

	// 创建一个新的cron实例
	c := cron.New(cron.WithSeconds())

	// 添加每小时发金币的定时任务
	_, err := c.AddFunc("0 0 * * * *", func() {
		mes, err := util.Rcon("main", "money give @a 1000")
		if err != nil {
			co.SysLog("ERROR", err.Error())
		} else {
			co.SysLog("INFO", mes)
		}
	})
	if err != nil {
		co.SysLog("ERROR", "发金币的定时任务添加失败")
		return
	}

	// 添加每日重置PearlCoin的定时任务
	_, err = c.AddFunc("0 0 0 * * *", func() {
		err := co.DB.Model(&co.AxolotlCoin{}).Updates(map[string]interface{}{"Pearl": 50}).Error
		if err != nil {
			co.SysLog("ERROR", err.Error())
		} else {
			co.SysLog("INFO", "已将所有玩家的PearlCoin重置为50")
		}
	})
	if err != nil {
		co.SysLog("ERROR", "每日重置PearlCoin的定时任务添加失败")
		return
	}

	// 启动cron实例
	c.Start()
}
