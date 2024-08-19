package workerpool

import (
	"errors"
	"fmt"
	"sync"
)

var (
	defaultCapacity    = 10
	maxCapacity        = 100
	ErrWorkerPoolFreed = errors.New("workerpool freed")
)

type Pool struct {
	capacity int
	active   chan struct{}
	tasks    chan Task
	wg       sync.WaitGroup //用于在pool销毁时等待所有worker退出
	quit     chan struct{}  // 用于通知各个worker退出的信号channel
}

type Task func()

func New(capacity int) *Pool {
	if capacity <= 0 {
		capacity = defaultCapacity
	}
	if capacity > maxCapacity {
		capacity = maxCapacity
	}

	p := &Pool{
		capacity: capacity,
		tasks:    make(chan Task),
		quit:     make(chan struct{}),
		active:   make(chan struct{}, capacity),
	}

	go p.run()
	return p
}

func (p *Pool) run() {
	idx := 0
	for {
		select {
		case <-p.quit:
			return
		case p.active <- struct{}{}:
			idx++
			fmt.Printf("worker[%03d]:active\n", idx)
			p.newWorker(idx)

		}
	}
}

func (p *Pool) newWorker(i int) {
	p.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]:recover panic[%s] and exit\n", i, err)
				<-p.active
			}
			p.wg.Done()
		}()

		fmt.Printf("worker[%03d]:start\n", i)
		for {
			select {
			case <-p.quit:
				fmt.Printf("worker[%03d]:exit\n", i)
				<-p.active
				return
			case t := <-p.tasks:
				fmt.Printf("worker[%03d]:receive a task\n", i)
				t()
			}

		}
	}()
}

func (p *Pool) Schedule(t Task) error {
	select {
	case <-p.quit:
		return ErrWorkerPoolFreed
	case p.tasks <- t:
		return nil
	}
}

func (p *Pool) Free() {
	close(p.quit)
	p.wg.Wait()
	fmt.Printf("workerpool freed!\n")
}
