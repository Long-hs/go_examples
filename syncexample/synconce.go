package syncexample

import (
	"log"
	"sync"
)

var (
	once sync.Once
)

func f() {
	log.Println("hahaha")
}

func syncOnce() {
	var group sync.WaitGroup
	for i := 0; i < 4; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			once.Do(f)
		}()
	}
	group.Wait()
	once.Do(f)
}
