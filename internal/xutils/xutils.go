package xutils

import "fmt"

/*
函数小写，表示私有，外部包无法访问
函数大写，表示公有，外部包可以访问
*/

func Test() {
	fmt.Println("utils包中的Test函数被调用")
}
