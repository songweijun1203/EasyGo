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

func mutateParamCopy(value int) int {
	ptr := &value
	*ptr = 99
	return value
}

func testParamAddress() int {
	value := 10
	changedCopy := mutateParamCopy(value)
	return value*100 + changedCopy
}

func setThroughPointer(ptr *int) {
	*ptr = 7
}

func testPointerParam() int {
	value := 1
	setThroughPointer(&value)
	return value
}

func testPointerPhi(useSecond bool) int {
	first := 1
	second := 2
	ptr := &first
	if useSecond {
		ptr = &second
	}
	*ptr = 9
	return first*10 + second
}

func testRangeStringUTF8() int {
	sum := 0
	for i := range "中a" {
		sum += i
	}
	return sum // 期望 3（0 + 3），不是 1
}

func testRangeMap() int {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	m["c"] = 3
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum // 6
}

func testRangeStringValue() int {
	sum := 0
	for _, r := range "中a" {
		sum += int(r)
	}
	return sum // 20013 + 97 = 20110
}

// 测试 20：string 索引 → Convert，验证 byte 不能泄漏出去
func testStringIndex() int {
	s := "abc"
	return int(s[0]) // 'a' = 97
}

// 测试 21：bool 相等比较
func testBoolEqual(a, b bool) int {
	if a == b {
		return 1
	}
	return 0
}

// 测试 22：位运算与移位
func testBitOps(x int) int {
	return (x&1)<<2 | x>>1 // x=5 → 1<<2 | 2 = 4|2 = 6
}

// 测试 23：指针比较的是地址 identity，不是 Cell 里的值
func testPointerEqual() int {
	a, b := 1, 2
	p, q := &a, &a
	r := &b
	if p == q && p != r {
		return 1
	}
	return 0
}
