package routine

import (
	"time"

	h "github.com/McaxDev/Back/handler"
	"github.com/McaxDev/Back/util"
)

var ScheduleList = []func(){
	util.ClearExpDefault(h.Challenges),
	util.ClearExpired(h.Mailsent, func(s h.MailStruct) time.Time {
		return s.Expiry
	}),
	util.ClearExpDefault(h.IpTimeMap),
}

func Schedule(minute int, operations ...func()) {
	ticker := time.NewTicker(time.Duration(minute) * time.Minute)
	for {
		<-ticker.C
		for _, operation := range operations {
			operation()
		}
	}
}
