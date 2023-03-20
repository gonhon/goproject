package current

import (
	"sync"
	"testing"
	"time"
)

var (
	group sync.WaitGroup
)

func TestLock(t *testing.T) {
	//生产者

	for i := 1; i <= MAX; i++ {
		go func() {
			group.Add(1)
			defer group.Done()
			time.Sleep(time.Millisecond * 200)
			productionLocal()
		}()
	}

	//消费者

	for j := 1; j <= MAX; j++ {
		go func() {
			group.Add(1)
			defer group.Done()
			time.Sleep(time.Millisecond * 200)
			consumeLock()
		}()
	}

	group.Wait()
}
