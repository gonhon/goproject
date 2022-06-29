package student

import (
	"fmt"
	"testing"
)

func TestStu(t *testing.T) {
	// stu := student{}
	stu := New(24, "gao")
	fmt.Println(stu)
}
