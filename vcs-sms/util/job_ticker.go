package util

import (
	"time"
)

type JobTicker struct {
	INTERVAL_PERIOD time.Duration
	HOUR_TO_TICK    int
	MINUTE_TO_TICK  int
	SECOND_TO_TICK  int
	timer           *time.Timer
}

func (t *JobTicker) DoPeriodicTask(callback func()) {
	go t.runningRoutine(callback)
}

func (t *JobTicker) runningRoutine(callback func()) {
	t.updateTimer()
	for {
		<-t.timer.C
		callback()
		t.updateTimer()
	}
}

func (t *JobTicker) updateTimer() {
	nextTick := time.Date(time.Now().Year(), time.Now().Month(),
		time.Now().Day(), t.HOUR_TO_TICK, t.MINUTE_TO_TICK, t.SECOND_TO_TICK, 0, time.Local)
	if !nextTick.After(time.Now()) {
		nextTick = nextTick.Add(t.INTERVAL_PERIOD)
	}
	diff := nextTick.Sub(time.Now())
	if t.timer == nil {
		t.timer = time.NewTimer(diff)
	} else {
		t.timer.Reset(diff)
	}
}
