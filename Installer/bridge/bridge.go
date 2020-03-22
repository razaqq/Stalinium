package bridge

import (
	"github.com/therecipe/qt/core"
	"time"
)

type AppBridge struct {
	core.QObject

	Total     uint64
	TotalDone uint64
	StartTime time.Time

	_ func(percentDone, totalDone, total, speed float64) `signal:"progress"`
	_ func(msg string)                                   `signal:"error"`
	_ func()                                             `signal:"success"`
}

func (ab *AppBridge) Write(p []byte) (int, error) {
	n := len(p)
	ab.TotalDone += uint64(n)
	ab.update()
	return n, nil
}

func (ab *AppBridge) update() {
	totalTime := time.Now().Sub(ab.StartTime).Seconds()
	percentDone := float64(ab.TotalDone) * 100 / float64(ab.Total)
	totalDone := float64(ab.TotalDone) / 1024 / 1024
	total := float64(ab.Total) / 1024 / 1024
	speed := float64(ab.TotalDone) / 1024 / 1024 / totalTime
	ab.Progress(percentDone, totalDone, total, speed)
}
