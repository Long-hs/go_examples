package syncexample

import (
	"log"
	"sync"
	"time"
)

var pool = sync.Pool{
	New: func() interface{} {
		return make([]int, 0)
	},
}

func workPool(i int) {
	// 从池中获取对象
	slice := pool.Get().([]int)
	log.Println("从池中获取：", slice)
	// 使用对象
	slice = append(slice, i)
	log.Println("使用：", slice)
	// 将对象放回池中
	pool.Put(slice)
	log.Println("放回池中")
}

func syncPool() {
	var group sync.WaitGroup
	for i := 0; i < 4; i++ {
		time.Sleep(time.Second * time.Duration(i))
		group.Add(1)
		index := i
		go func() {
			workPool(index)
			defer group.Done()
		}()
	}
	group.Wait()
}
