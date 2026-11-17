package xinterp

func add(a, b int) int {
	return a + b
}

// 测试 if
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
