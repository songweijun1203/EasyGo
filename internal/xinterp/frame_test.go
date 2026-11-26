package xinterp // 注意：这里必须是 xinterp，而不是 xinterp_test

import (
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestControlFlow(t *testing.T) {
	packagePath := "EasyGo/internal/xinterp"
	interp := NewInterp()
	err := interp.Load(
		&packages.Config{
			Mode: packages.LoadAllSyntax,
			Dir:  "./",
		},
		"./",
	)
	if err != nil {
		t.Fatalf("Load 失败: %v", err)
	}

	// 测试 1：add (线性)
	{
		addFn, _ := interp.LookupFunction(packagePath, "add")
		result := interp.RunFunction(addFn, 10, 20)

		if result != 30 {
			t.Fatalf("add(10, 20) = %v, want 30", result)
		}
		println("add(10, 20) = ", result.(int))
	}

	// 测试 2：max (if)
	{
		maxFn, _ := interp.LookupFunction(packagePath, "max")
		result := interp.RunFunction(maxFn, 10, 20)

		if result != 20 {
			t.Fatalf("max(10, 20) = %v, want 20", result)
		}
		println("max(10, 20) = ", result.(int))
	}

	// 测试 3：sum (for)
	{
		sumFn, _ := interp.LookupFunction(packagePath, "sum")
		result := interp.RunFunction(sumFn, 5)

		if result != 15 {
			t.Fatalf("sum(5) = %v, want 15", result)
		}
		println("sum(5) = ", result.(int))
	}

	//测试4 函数调用
	{
		addAndmulFn, _ := interp.LookupFunction(packagePath, "addAndMul")

		result := interp.RunFunction(addAndmulFn, 2, 5, 2)

		if result != 14 {
			t.Fatalf("addAndmulFn(2,5,2) = %v, want 14", result)
		}
		println("addAndmulFn(2,5,2) = ", result.(int))
	}

	//测试5 内建函数的调用
	{
		builtinFn, _ := interp.LookupFunction(packagePath, "testBuiltin")
		gotBuiltin := interp.RunFunction(builtinFn)
		if gotBuiltin != 5 {
			t.Fatalf("testBuiltin 失败: got %v, want 5", gotBuiltin)
		}
	}

	//测试6 外部库调用
	{
		upperFn, _ := interp.LookupFunction(packagePath, "testExternalUpper")
		got := interp.RunFunction(upperFn)

		if got != "HELLO" {
			t.Fatalf("testExternalUpper() = %v, want HELLO", got)
		}
		println("testExternalUpper() = ", got.(string))
	}

	// 测试7 全局变量测试
	{
		printFn, _ := interp.LookupFunction(packagePath, "printCount")
		interp.RunFunction(printFn)
		interp.RunFunction(printFn)
		interp.RunFunction(printFn)
	}

	// 测试8 多返回值
	{
		Fn, _ := interp.LookupFunction(packagePath, "testSwap")
		interp.RunFunction(Fn)
	}

	// 测试9 print
	{
		Fn, _ := interp.LookupFunction(packagePath, "myPrintIn")
		interp.RunFunction(Fn)
	}

	// 测试10 struct测试
	{
		Fn, _ := interp.LookupFunction(packagePath, "testStruct")
		interp.RunFunction(Fn)
	}

	// 测试 11：append
	{
		Fn, _ := interp.LookupFunction(packagePath, "testAppend")
		interp.RunFunction(Fn)
	}

	{
		Fn, _ := interp.LookupFunction(packagePath, "printSlice")
		interp.RunFunction(Fn)
	}

	// 测试 12：map 写入、更新和读取
	{
		fn, _ := interp.LookupFunction(packagePath, "testMapUpdate")
		result := interp.RunFunction(fn)
		if result != 42 {
			t.Fatalf("testMapUpdate() = %v, want 42", result)
		}
	}

	// 测试 13：闭包测试
	{
		Fn, _ := interp.LookupFunction(packagePath, "testClosure")
		interp.RunFunction(Fn)
	}

	// 测试 14：对形参取地址只修改被调用函数中的参数副本
	{
		fn, _ := interp.LookupFunction(packagePath, "testParamAddress")
		if result := interp.RunFunction(fn); result != 1099 {
			t.Fatalf("testParamAddress() = %v, want 1099", result)
		}
	}

	// 测试 15：指针参数保留地址并写回调用者的存储槽
	{
		fn, _ := interp.LookupFunction(packagePath, "testPointerParam")
		if result := interp.RunFunction(fn); result != 7 {
			t.Fatalf("testPointerParam() = %v, want 7", result)
		}
	}

	// 测试 16：Phi 在分支汇合时保留指针值，不读取 Cell 内容
	{
		fn, _ := interp.LookupFunction(packagePath, "testPointerPhi")
		if result := interp.RunFunction(fn, false); result != 92 {
			t.Fatalf("testPointerPhi(false) = %v, want 92", result)
		}
		if result := interp.RunFunction(fn, true); result != 19 {
			t.Fatalf("testPointerPhi(true) = %v, want 19", result)
		}
	}

	// 测试 17：range map 累加所有 value（顺序无关，断言安全）
	{
		fn, _ := interp.LookupFunction(packagePath, "testRangeMap")
		if result := interp.RunFunction(fn); result != 6 {
			t.Fatalf("testRangeMap() = %v, want 6", result)
		}
	}

	// 测试 18：range string 的 key 是字节下标，必须按 UTF-8 宽度前进
	{
		fn, _ := interp.LookupFunction(packagePath, "testRangeStringUTF8")
		if result := interp.RunFunction(fn); result != 3 {
			t.Fatalf("testRangeStringUTF8() = %v, want 3", result)
		}
	}

	// 测试 19：range string 的 value 是 rune，且与本解释器 int 归一化自洽
	{
		fn, _ := interp.LookupFunction(packagePath, "testRangeStringValue")
		if result := interp.RunFunction(fn); result != 20110 {
			t.Fatalf("testRangeStringValue() = %v, want 20110", result)
		}
	}
}
