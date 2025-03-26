package goroutine_example

import (
	"log"
	"sync"
)

type ITask interface {
	Run()
}
type IGoPool interface {
	IStart()
	ISchedule(task ITask)
	IClose()
}
type G struct {
	id    int //协程id
	count int //消费任务数量
}

func NewG(id int) *G {
	return &G{
		id:    id,
		count: 0,
	}
}

type gPool struct {
	workers int        //并发协程数
	taskCh  chan ITask //任务管道
	wg      sync.WaitGroup
	Gs      map[int]*G
}

func NewPool(workers int) IGoPool {
	Gs := make(map[int]*G)
	for i := 0; i < workers; i++ {
		g := NewG(i)
		Gs[i] = g
	}
	return &gPool{
		workers: workers,
		taskCh:  make(chan ITask),
		wg:      sync.WaitGroup{},
		Gs:      Gs,
	}
}

func (g *gPool) IStart() {
	for _, g2 := range g.Gs {
		g.wg.Add(1)
		go func(g2 *G) {
			defer g.wg.Done()
			for task := range g.taskCh {
				task.Run()
				g2.count++
			}
		}(g2)
	}
}

func (g *gPool) ISchedule(task ITask) {
	g.taskCh <- task
}

func (g *gPool) IClose() {
	close(g.taskCh)
	g.wg.Wait()
	log.Printf("所有协程任务完成\n")
	for _, g2 := range g.Gs {
		log.Printf("协程 %d 消费了 %d 个任务\n", g2.id, g2.count)
	}
}

type taskStrcut struct {
	id int
}

func (s *taskStrcut) Run() {
	log.Printf("我是%d\n", s.id)
}

func NewTask(id int) ITask {
	return &taskStrcut{id: id}
}

func RunExample() {
	pool := NewPool(5)
	pool.IStart()
	for i := 0; i < 10000; i++ {
		newTask := NewTask(i)
		pool.ISchedule(newTask)
	}
	pool.IClose()
}
