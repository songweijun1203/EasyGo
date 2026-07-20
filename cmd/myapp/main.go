package main //package main表示一个可独立执行的程序，每个 Go 应用程序都包含一个名为 main 的包

//fmt 格式化输入输出的标准库，提供了格式化的 I/O 函数
import (
	xq_utils "EasyGo/internal/utils"
	"fmt"
)

func var_init() {
	var age = 10
	var isboy = true
	var weight float32 = 60.5
	var width float64 = 100.5
	var heigt int = 180
	fmt.Println("age:", age)
	fmt.Println("isboy:", isboy)
	fmt.Println("weight:", weight)
	fmt.Println("width:", width)
	fmt.Println("heigt:", heigt)
}

func test() {
	numbers := []int{1, 2, 3, 4, 5}

	// 标准 Go 的循环写法
	squares := make([]int, len(numbers))
	for i, n := range numbers {
		squares[i] = n * n
	}

	xq_utils.Test() // 调用 utils 包中的 test 函数

	fmt.Println("原始数组:", numbers)
	fmt.Println("平方数组:", squares)
}

func main() {
	var_init()
}
