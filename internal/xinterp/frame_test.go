package xinterp // 注意：这里必须是 xinterp，而不是 xinterp_test

import (
	"go/types"
	"testing"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// setupSSA 是一个辅助函数，用于在测试中设置 SSA 环境。
// 准备环境，只加载一次，后面的测试共用
func setupSSA(t *testing.T) []*packages.Package {
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  "./",
	}
	pkgs, err := packages.Load(cfg, "./")
	if err != nil || packages.PrintErrors(pkgs) > 0 {
		t.Fatal("加载包失败")
	}
	return pkgs
}

// 按照字符串查找函数
func findFunctionByName(t *testing.T, pkgs []*packages.Package, funcName string) *ssa.Function {
	prog, _ := ssautil.AllPackages(pkgs, ssa.InstantiateGenerics)
	prog.Build()

	for _, pkg := range pkgs {
		for _, mem := range pkg.Types.Scope().Names() {
			if obj, ok := pkg.Types.Scope().Lookup(mem).(*types.Func); ok {
				if fn := prog.FuncValue(obj); fn != nil && fn.Name() == funcName {
					return fn
				}
			}
		}
	}
	t.Fatalf("没找到函数: %s", funcName)
	return nil
}

func TestControlFlow(t *testing.T) {
	//pkgs := setupSSA(t)
	interp := NewInterp()
	err := interp.Load(
		&packages.Config{
			Mode: packages.LoadAllSyntax,
			Dir:  "./",
		},
		"./",
	)
	if err != nil {
		return
	}

	pkgs := interp.sourcePkgs
	// 测试 1：add (线性)
	{
		addFn := findFunctionByName(t, pkgs, "add")
		result := interp.RunFunction(addFn, 10, 20)

		if result != 30 {
			t.Fatalf("add(10, 20) = %v, want 30", result)
		}
		println("add(10, 20) = ", result.(int))
	}

	// 测试 2：max (if)
	{
		maxFn := findFunctionByName(t, pkgs, "max")
		result := interp.RunFunction(maxFn, 10, 20)

		if result != 20 {
			t.Fatalf("max(10, 20) = %v, want 20", result)
		}
		println("max(10, 20) = ", result.(int))
	}

	// 测试 3：sum (for)
	{
		sumFn := findFunctionByName(t, pkgs, "sum")
		result := interp.RunFunction(sumFn, 5)

		if result != 15 {
			t.Fatalf("sum(5) = %v, want 15", result)
		}
		println("sum(5) = ", result.(int))
	}

	//测试4 函数调用
	{
		addAndmulFn := findFunctionByName(t, pkgs, "addAndMul")
		result := interp.RunFunction(addAndmulFn, 2, 5, 2)

		if result != 14 {
			t.Fatalf("addAndmulFn(2,5,2) = %v, want 14", result)
		}
		println("addAndmulFn(2,5,2) = ", result.(int))
	}

	//测试5 内建函数的调用
	{
		builtinFn := findFunctionByName(t, pkgs, "testBuiltin")
		gotBuiltin := interp.RunFunction(builtinFn)
		if gotBuiltin != 5 {
			t.Fatalf("testBuiltin 失败: got %v, want 5", gotBuiltin)
		}
	}

	//测试6 外部库调用
	{
		upperFn := findFunctionByName(t, pkgs, "testExternalUpper")
		got := interp.RunFunction(upperFn)

		if got != "HELLO" {
			t.Fatalf("testExternalUpper() = %v, want HELLO", got)
		}
		println("testExternalUpper() = ", got.(string))
	}

	// 测试7 全局变量测试
	{
		printFn := findFunctionByName(t, pkgs, "printCount")
		interp.RunFunction(printFn)
		interp.RunFunction(printFn)
		interp.RunFunction(printFn)
	}

	// 测试8 多返回值
	{
		Fn := findFunctionByName(t, pkgs, "testSwap")
		interp.RunFunction(Fn)
	}

	// 测试9 print
	{
		Fn := findFunctionByName(t, pkgs, "myPrintIn")
		interp.RunFunction(Fn)
	}

	// 测试10 struct测试
	{
		Fn := findFunctionByName(t, pkgs, "testStruct")
		interp.RunFunction(Fn)
	}

	// 测试 11：append
	{
		Fn := findFunctionByName(t, pkgs, "testAppend")
		interp.RunFunction(Fn)
	}

	{
		Fn := findFunctionByName(t, pkgs, "printSlice")
		interp.RunFunction(Fn)
	}

	// 测试 12：map 写入、更新和读取
	{
		fn := findFunctionByName(t, pkgs, "testMapUpdate")
		result := interp.RunFunction(fn)
		if result != 42 {
			t.Fatalf("testMapUpdate() = %v, want 42", result)
		}
	}

	// 测试 13：闭包测试
	{
		Fn := findFunctionByName(t, pkgs, "testClosure")
		interp.RunFunction(Fn)
	}

	// 测试 14：对形参取地址只修改被调用函数中的参数副本
	{
		fn := findFunctionByName(t, pkgs, "testParamAddress")
		if result := interp.RunFunction(fn); result != 1099 {
			t.Fatalf("testParamAddress() = %v, want 1099", result)
		}
	}

	// 测试 15：指针参数保留地址并写回调用者的存储槽
	{
		fn := findFunctionByName(t, pkgs, "testPointerParam")
		if result := interp.RunFunction(fn); result != 7 {
			t.Fatalf("testPointerParam() = %v, want 7", result)
		}
	}

	// 测试 16：Phi 在分支汇合时保留指针值，不读取 Cell 内容
	{
		fn := findFunctionByName(t, pkgs, "testPointerPhi")
		if result := interp.RunFunction(fn, false); result != 92 {
			t.Fatalf("testPointerPhi(false) = %v, want 92", result)
		}
		if result := interp.RunFunction(fn, true); result != 19 {
			t.Fatalf("testPointerPhi(true) = %v, want 19", result)
		}
	}
}
