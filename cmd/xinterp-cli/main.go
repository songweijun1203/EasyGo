package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"EasyGo/internal/xinterp" // 替换为你的真实 module 名

	"go/types"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func main() {
	// 1. 准备临时沙盒目录（存放在系统临时目录中，避免污染项目）
	tmpDir, err := os.MkdirTemp("", "xinterp-repl-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir) // 程序退出时清理沙盒

	// 2. 将 lib.go 的内容读出来，准备拼接
	libBytes, err := os.ReadFile("./cmd/xinterp-cli/lib.go")
	if err != nil {
		panic("找不到 lib.go，请在 cmd/xinterp-cli 目录下运行")
	}
	// 去掉原本的 package main，避免重复
	libCode := strings.Replace(string(libBytes), "package main", "", 1)

	// 3. 写入 go.mod，让沙盒成为一个合法的 Go Module
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module repl-sandbox\n\ngo 1.20\n"), 0644)

	fmt.Println("=====================================")
	fmt.Println("  Go SSA 解释器 REPL (输入 exit 退出)")
	fmt.Println("  支持: add(1, 2), fib(10), testFFI()")
	fmt.Println("=====================================")

	scanner := bufio.NewScanner(os.Stdin)

	// 4. 进入 REPL 循环
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

		// 5. 动态拼接代码
		// 把用户的输入包装在 __repl_exec__ 中，并用 println 打印结果
		sourceCode := fmt.Sprintf(`package main
%s

func __repl_exec__() {
	println("=>", %s)
}
`, libCode, input)

		// 写入临时沙盒文件
		tmpFile := filepath.Join(tmpDir, "main.go")
		os.WriteFile(tmpFile, []byte(sourceCode), 0644)

		// 6. 加载并构建 SSA
		cfg := &packages.Config{
			Mode: packages.LoadAllSyntax,
			Dir:  tmpDir, // 指向我们的沙盒目录
		}
		pkgs, err := packages.Load(cfg, ".")
		if err != nil || packages.PrintErrors(pkgs) > 0 {
			fmt.Println("语法错误，请重新输入")
			continue
		}

		prog, _ := ssautil.AllPackages(pkgs, ssa.InstantiateGenerics)
		prog.Build()

		// 7. 寻找执行入口 __repl_exec__
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

		// 8. 执行并捕获 Panic，防止解释器崩溃退出 REPL
		if execFn != nil {
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("执行崩溃: %v\n", r)
					}
				}()
				// 调用你的解释器执行这个无参函数
				xinterp.RunFunction(execFn)
			}()
		}
	}
}
