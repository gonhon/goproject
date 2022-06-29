package waitgroup

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	id   int
	rand int
}

type Result struct {
	job Job
	sum int
}

var (
	jobs = make(chan Job, 10)
	res  = make(chan Result, 10)
)

//计算一个数每一位加起来的和
func dist(num int) int {
	sum := 0
	no := num
	for no != 0 {
		temp := no % 10
		sum += temp
		no /= 10
	}
	return sum
}

//根据Job生成Res
func worker(wait *sync.WaitGroup) {
	for job := range jobs {
		res <- Result{job, dist(job.rand)}
	}
	wait.Done()
}

//创建工作组
func createWorkerPool(num int) {

	var wait sync.WaitGroup
	for i := 0; i < num; i++ {
		wait.Add(1)
		go worker(&wait)
	}
	wait.Wait()
	close(res)
}

//创建Job
func createJob(num int) {
	for i := 0; i < num; i++ {
		jobs <- Job{i, rand.Intn(1000)}
	}
	close(jobs)
}

func resPrintf(flag chan bool) {
	for r := range res {
		fmt.Printf("Job id %d, input random no %d , sum of digits %d\n", r.job.id, r.job.rand, r.sum)
	}
	flag <- true
}

func Run() {

	startTime := time.Now()
	noJobs := 10
	go createJob(noJobs)
	flag := make(chan bool)
	go resPrintf(flag)
	noWorkers := 10
	createWorkerPool(noWorkers)
	<-flag
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("耗时:", diff)

}
