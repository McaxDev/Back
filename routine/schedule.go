package routine

import "time"

func Schedule(minute int, operation func()) {
	ticker := time.NewTicker(time.Duration(minute) * time.Minute)
	for {
		<-ticker.C
		operation()
	}
}
