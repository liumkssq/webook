package timer

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"testing"
	"time"
)

func TestTimerTicker(t *testing.T) {
	tm := time.NewTicker(time.Second)
	defer tm.Stop()
	for now := range tm.C {
		t.Log(now)
	}
}

func TestCron(t *testing.T) {
	expr := cron.New(cron.WithSeconds())
	expr.AddJob("@every 1s", myJob{})
	expr.AddFunc("@every 1s", func() {
		t.Log("start")
		time.Sleep(10 * time.Second)
		t.Log("end")
	})
	expr.Start()
	time.Sleep(5 * time.Second)
	stop := expr.Stop()
	t.Log("send stop signal to cron")
	<-stop.Done()
	t.Log("cron stopped")
}

type myJob struct {
}

func (j myJob) Run() {
	fmt.Println("hello world")
}
