package io

import (
	"bufio"
	"fmt"
	"os"
)

func Input() {

	reader := bufio.NewReader(os.Stdin)

	// for {
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("read-->", err)
	}

	fmt.Println("read-->", line)
}
