package concurrence

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

type Score struct {
	Num int
}

func (s *Score) Do() error {
	fmt.Println("num:", s.Num)
	time.Sleep(1 * 1 * time.Second)

	return nil
}

func TestConsurrence(t *testing.T) {
	num := 100 * 100
	p := NewWorkerPool(num)
	p.Run()
	datanum := 100 * 100 * 100
	go func() {
		for i := 1; i <= datanum; i++ {
			sc := &Score{Num: i}
			p.JobQueue <- sc
		}
	}()

	for {
		fmt.Println("runtime.NumGoroutine() :", runtime.NumGoroutine())
		time.Sleep(2 * time.Second)
	}

}
