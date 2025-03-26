package syncexample

import (
	"log"
	"sync"
	"time"
)

var group sync.WaitGroup

func workWaitGroup(i int) {
	log.Printf("我是%d\n", i)
	time.Sleep(time.Second * time.Duration(i))
	log.Printf("%d 完成工作\n", i)
	group.Done()
}

func syncWaitGroup() {
	for i := 0; i < 4; i++ {
		go workWaitGroup(i)
		group.Add(1)
	}
	group.Wait()
	log.Println("所有 workWaitGroup 完成工作")
}
