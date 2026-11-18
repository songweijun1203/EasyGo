package main

import "fmt"

// 你可以在这里写任何你希望 REPL 支持的业务函数
func add(a, b int) int {
	return a + b
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func testFFI() {
	// 测试外部库调用
	fmt.Println("FFI Bridge Works!")
}
