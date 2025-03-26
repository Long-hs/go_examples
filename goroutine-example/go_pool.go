package goroutine_example

import (
	"log"
	"sync"
)

func task(id int) {
	log.Printf("我是：%d\n", id)
}

func goroutineExample() {
	workers, tasks := 5, 10000

	var wg sync.WaitGroup

	taskCh := make(chan int)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			for i := range taskCh {
				task(i)
			}
			defer wg.Done()
		}()
	}

	for i := 0; i < tasks; i++ {
		taskCh <- i
	}
	close(taskCh)
	wg.Wait()
	log.Printf("所有协程任务完成\n")
}
