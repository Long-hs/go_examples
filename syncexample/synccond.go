package syncexample

import (
	"log"
	"sync"
	"time"
)

var (
	mutex sync.Mutex
	cond  = sync.NewCond(&mutex)
)

func workCond(i int) {
	cond.L.Lock()
	defer cond.L.Unlock()
	log.Printf("%d正在等待\n", i)
	cond.Wait()
	log.Printf("我是:%d\n", i)
}

func syncCond() {
	for i := 0; i < 4; i++ {
		go workCond(i)
	}
	time.Sleep(time.Second * 2)
	//只通知一个协程(等待队列的第一个)
	signal()
	//通知所有协程
	broadcast()
	time.Sleep(time.Second * 10)
}

func broadcast() {
	cond.Broadcast()
}

func signal() {
	cond.Signal()
}
