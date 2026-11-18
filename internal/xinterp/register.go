package xinterp

import (
	"fmt"
	"reflect"
	"strings"
)

// 核心注册表：存储SSA中的函数名->真实的Go函数反射值
var nativeRegistry = make(map[string]reflect.Value)

// 提供给外部（比如init函数）注册真实函数的方法
func RegisterNative(name string, fn interface{}) {
	nativeRegistry[name] = reflect.ValueOf(fn)
}

// 预注册一些标准库函数
func init() {
	RegisterNative("fmt.Println", fmt.Println)
	RegisterNative("strings.ToUpper", strings.ToUpper)
}

// 测试外部库：定长参数，单返回值
func testExternalUpper() string {
	return strings.ToUpper("hello")
}
