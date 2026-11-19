package xinterp

import "fmt"

var g_count int

func add(a, b int) int {
	return a + b
}

// 测试
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 测试 for
func sum(n int) int {
	total := 0
	for i := 1; i <= n; i++ {
		total += i
	}
	return total
}

func mul(a int, factor int) int {
	return a * factor
}
func addAndMul(a, b int, factor int) int {
	sum := a + b
	return mul(sum, factor)
}

func testBuiltin() int {
	println("Hello, SSA!")
	s := "hello"
	return len(s) // 应该返回 5
}

func printCount() {
	//fmt.Printf("g_count: %d\n", g_count)
	println("Hello, g_count!")
	g_count++
}

func swap(a, b int) (int, int) {
	return b, a
}

func testSwap() {
	x, y := swap(10, 20)
	// SSA 这里会产生一个 Call，然后跟着两个 Extract #0 和 Extract #1
	println(x, y) // 期望输出：20 10
}

func getTwoValue() (int, int) {
	println(10, 20)
	return 10, 20
}

func myPrintIn() {
	fmt.Println(10)
}
