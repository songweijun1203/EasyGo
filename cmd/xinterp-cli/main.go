package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// 替换为你的真实 module 名
	"EasyGo/internal/xinterp"

	"go/types"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func main() {
	// 1. 创建临时沙盒目录
	tmpDir, err := os.MkdirTemp("", "xinterp-repl-*")
	if err != nil {
		fmt.Println("创建沙盒失败:", err)
		return
	}
	defer os.RemoveAll(tmpDir) // 退出时清理沙盒

	// 2. 写入 go.mod，让沙盒成为一个合法的 Go Module
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module repl-sandbox\n\ngo 1.20\n"), 0644)

	// 3. 读取基础库代码 (lib.go)
	libBytes, err := os.ReadFile("./cmd/xinterp-cli/lib.go")
	if err != nil {
		panic("找不到 lib.go，请在项目根目录下运行")
	}
	// 去掉原本的 package main 声明，避免重复
	libCode := strings.Replace(string(libBytes), "package main", "", 1)

	// 用来保存用户历史输入的语句（用于状态保留）
	historyCode := ""

	fmt.Println("=====================================")
	fmt.Println("  Go SSA 解释器 REPL (输入 exit 退出)")
	fmt.Println("  支持: add(1, 2), testFFI(), a := 10")
	fmt.Println("=====================================")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">>> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "exit" || input == "quit" {
			break
		}
		if input == "" {
			continue
		}

		// 执行结果状态
		execSuccess := false

		// ==========================================
		// 尝试一：作为“表达式”执行（自动打印返回值）
		// ==========================================
		codeExpr := fmt.Sprintf(`package main
%s
%s

func __repl_exec__() {
	println("=>", %s)
}
`, libCode, historyCode, input)

		tmpFile := filepath.Join(tmpDir, "main.go")
		os.WriteFile(tmpFile, []byte(codeExpr), 0644)

		cfg := &packages.Config{Mode: packages.LoadAllSyntax, Dir: tmpDir}
		pkgs, err := packages.Load(cfg, ".")

		// 如果加载没有错误，说明是合法的表达式
		if err == nil && packages.PrintErrors(pkgs) == 0 {
			runExec(pkgs)
			execSuccess = true
		} else {
			// ==========================================
			// 尝试二：作为“语句”执行（无返回值，或者是变量声明）
			// ==========================================
			codeStmt := fmt.Sprintf(`package main
%s
%s

func __repl_exec__() {
	%s
}
`, libCode, historyCode, input)

			os.WriteFile(tmpFile, []byte(codeStmt), 0644)
			pkgs, err = packages.Load(cfg, ".")

			if err == nil && packages.PrintErrors(pkgs) == 0 {
				runExec(pkgs)
				execSuccess = true

				// 【状态保留】如果执行成功，且不是表达式，把它加到历史记录中
				// 防止每次输入都被闭包在函数内，导致变量丢失
				if strings.Contains(input, ":=") || strings.Contains(input, "=") {
					// 简单粗暴地将 := 转换为包级别的 var 声明，防止编译报错
					// 注意：这里只做了极简处理，真实场景需要 AST 解析
					if strings.Contains(input, ":=") {
						parts := strings.SplitN(input, ":=", 2)
						varName := strings.TrimSpace(parts[0])
						varVal := strings.TrimSpace(parts[1])
						historyCode += fmt.Sprintf("\nvar %s = %s", varName, varVal)
					} else {
						// 如果是普通语句（如函数调用），直接追加到历史
						// 但这里为了防止重复执行，暂时不追加非赋值语句
					}
				}
			}
		}

		if !execSuccess {
			fmt.Println("语法错误或执行失败，请重新输入")
		}
	}
}

// 运行生成的 __repl_exec__ 函数
func runExec(pkgs []*packages.Package) {
	prog, _ := ssautil.AllPackages(pkgs, ssa.InstantiateGenerics)
	prog.Build()

	var execFn *ssa.Function
	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup("__repl_exec__")
		if obj != nil {
			if fnObj, ok := obj.(*types.Func); ok {
				execFn = prog.FuncValue(fnObj)
				break
			}
		}
	}

	if execFn != nil {
		// 捕获 Panic，防止解释器崩溃
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("\n[执行异常] %v\n", r)
				}
			}()
			xinterp.RunFunction(execFn)
		}()
	}
}
