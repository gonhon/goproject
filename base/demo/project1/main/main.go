package main

import "fmt"

func cla(a, b int, test byte) int {
	switch test {
	case '+':
		return a + b
	case '-':
		return a - b
	case '*':
		return a * b
	case '/':
		return a / b
	}
	return 0
}

// 闭包
func add(test byte) func(a, b int) int {
	return func(a, b int) int {
		switch test {
		case '+':
			return a + b
		case '-':
			return a - b
		case '*':
			return a * b
		case '/':
			return a / b
		}
		return 0
	}
}
func main() {
	fmt.Println("123")
	a := add('+')
	fmt.Println(cla(12, 1, '+'))
	fmt.Println(a(12, 12))
}
