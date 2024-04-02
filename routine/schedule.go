package routine

import "time"

func Schedule(minute int, operations ...func()) {
	ticker := time.NewTicker(time.Duration(minute) * time.Minute)
	for {
		<-ticker.C
		for _, operation := range operations {
			operation()
		}
	}
}
