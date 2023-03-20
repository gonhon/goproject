package current

import (
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	var prodconsLock = ProdconsLock{}
	Exec(prodconsLock, time.Millisecond*200, time.Millisecond*200)

}

func TestChan(t *testing.T) {
	var prodconsChan = ProdconsChan{}
	Exec(prodconsChan, time.Millisecond*200, time.Millisecond*200)
}
