package xinterp // 注意：这里必须是 xinterp，而不是 xinterp_test

import (
	"fmt"
	"go/types"
	"testing"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func TestExecuteAdd(t *testing.T) {
	// 1. 初始化加载配置。
	// 这里使用 LoadAllSyntax，确保 Go 不仅解析语法树（AST），
	// 还会进行完整的类型检查（Type Checking），这是后续构建 SSA 的必要前提。
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  "./",
	}

	// 2. 加载目标代码。
	// 参数 "./" 表示加载当前目录下所有的 Go 文件。
	// 这里的“当前目录”是指 internal/xinterp。
	// 因为该目录下所有 .go 文件的包名均为 xinterp，所以它们会被合并编译为一个 Package。
	// 返回值 pkgs 是一个切片，虽然当前 "./" 只会匹配到一个包（长度为 1），
	// 但 packages.Load 支持传入多个模式（如 "./..."），因此设计为切片。
	pkgs, err := packages.Load(cfg, "./")
	if err != nil {
		t.Fatalf("Failed to load packages: %v", err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		t.Fatal("包加载有错误")
	}

	// 3. 构建 SSA 中间表示（IR）。
	// sautil.AllPackages 会递归抓取当前包及其所有依赖包（包括标准库），构建整个程序的 SSA 视图
	// prog.Build() 会对所有函数进行 SSA 构建，生成完整的中间表示。
	prog, _ := ssautil.AllPackages(pkgs, ssa.InstantiateGenerics)
	prog.Build()

	//打印package相关信息
	for _, pkg := range pkgs {
		fmt.Printf("=== 包信息 ===\n")
		fmt.Printf("PkgPath:   %v\n", pkg.PkgPath)
		fmt.Printf("Types:     %v\n", pkg.Types.String())
		fmt.Printf("GoFiles:   %d files: %v\n", len(pkg.GoFiles), pkg.GoFiles)

		// 只打印 TypesInfo 的统计信息
		ti := pkg.TypesInfo
		if ti != nil {
			fmt.Printf("TypesInfo 统计: Types=%d, Defs=%d, Uses=%d, Implicits=%d\n",
				len(ti.Types), len(ti.Defs), len(ti.Uses), len(ti.Implicits))
		}

		// 只打印 Errors 的数量和内容
		if len(pkg.Errors) > 0 {
			fmt.Printf("Errors (%d):\n", len(pkg.Errors))
			for _, err := range pkg.Errors {
				fmt.Printf("  - %s\n", err)
			}
		} else {
			fmt.Println("Errors: 无")
		}

		fmt.Printf("Types.Scope().Names():   %d names: %v\n", len(pkg.Types.Scope().Names()), pkg.Types.Scope().Names())

		fmt.Println()
	}

	var addFn *ssa.Function
	funcName := "add"

	for _, pkg := range pkgs {
		for _, mem := range pkg.Types.Scope().Names() {
			if obj, ok := pkg.Types.Scope().Lookup(mem).(*types.Func); ok {
				if fn := prog.FuncValue(obj); fn != nil && fn.Name() == funcName {
					addFn = fn
					break
				}
			}
		}
	}

	if addFn == nil {
		t.Fatal("没找到 add 函数")
	}

	// 因为现在是同包，直接调用 RunAdd 即可，不需要 xinterp. 前缀
	got := RunAdd(addFn, 10, 20)

	if got != 30 {
		t.Fatalf("RunAdd(10, 20) = %v, want 30", got)
	}
}

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
	pkgs := setupSSA(t)
	interp := &Interp{
		globals: make(map[*ssa.Global]*Value),
	}

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
}
