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

type Person struct {
	Name string
	Age  int
}

func testStruct() {
	// 1. 结构体初始化 (SSA 会调用 Alloc 造一个 map)
	var p Person

	// 2. 字段赋值 (SSA 会生成 FieldAddr 获取地址，然后 Store 写入)
	p.Name = "Alice"
	p.Age = 25

	// 3. 字段读取 (SSA 会生成 Field 读取值)
	fmt.Printf("Name: %s, Age: %d\n", p.Name, p.Age)

	// 4. 传递结构体指针 (SSA 会传递指向 map 的指针)
	updateAge(&p)
	fmt.Println("Updated Age:", p.Age)
}
func updateAge(ptr *Person) {
	ptr.Age = 30
}

func testAppend() {
	// 1. 普通追加
	s := []int{1, 2, 3}
	s = append(s, 4, 5)
	fmt.Println("After append:", s) // 期望：[1 2 3 4 5]

	// 2. 展开追加
	t := []int{6, 7, 8}
	s = append(s, t...)
	fmt.Println("After spread append:", s) // 期望：[1 2 3 4 5 6 7 8]

	// 3. 对 nil 切片 append
	var nilSlice []int
	nilSlice = append(nilSlice, 100)
	fmt.Println("Nil slice append:", nilSlice) // 期望：[100]

	// 4. 验证 len
	fmt.Println("Final len:", len(s)) // 期望：8
}

func printSlice() {
	s := []int{1, 2, 3}
	fmt.Println(s)
}

func testMapUpdate() int {
	m := make(map[string]int)
	m["answer"] = 1
	m["answer"] = 42
	return m["answer"]
}

// 测试闭包：计数器
func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func testClosure() {
	counter := makeCounter()

	println("First call:", counter())  // 期望输出: 1
	println("Second call:", counter()) // 期望输出: 2
	println("Third call:", counter())  // 期望输出: 3
}
