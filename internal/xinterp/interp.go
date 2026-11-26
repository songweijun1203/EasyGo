package xinterp

import (
	"errors"
	"fmt"
	"sync"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type Interp struct {
	globals map[*ssa.Global]*Cell // 全程序唯一，所有帧共享

	prog       *ssa.Program
	sourcePkgs []*packages.Package
	ssaPkgs    []*ssa.Package

	// 用导入路径快速查找 SSA 包
	packagesByPath map[string]*ssa.Package

	layoutMu sync.RWMutex                  // 保护 layouts
	layouts  map[*ssa.Function]*funcLayout // 函数布局缓存
}

func NewInterp() *Interp {
	return &Interp{
		globals: make(map[*ssa.Global]*Cell),
	}
}

func (i *Interp) RunFunction(fn *ssa.Function, args ...Value) Value {
	f := NewFrame(i, fn)

	// 参数是 SSA 值，直接绑定到寄存器。需要取地址时，SSA 会显式生成 Alloc。
	for index, param := range fn.Params {
		value := zeroOf(param.Type())
		if index < len(args) {
			value = args[index]
		}
		f.setValue(param, value)
	}

	f.run() // 执行当前帧的指令

	return f.result // 返回函数执行的结果
}

// 加载磁盘上的包
func (i *Interp) Load(
	cfg *packages.Config,
	patterns ...string,
) error {
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return err
	}

	// 检查 pkgs 中的编译和类型错误。
	prog, ssaPkgs := ssautil.AllPackages(
		pkgs,
		ssa.InstantiateGenerics,
	)
	prog.Build()

	i.prog = prog
	i.sourcePkgs = pkgs
	i.ssaPkgs = ssaPkgs
	i.globals = make(map[*ssa.Global]*Cell)

	i.layoutMu.Lock()
	i.layouts = make(map[*ssa.Function]*funcLayout)
	i.layoutMu.Unlock()

	i.packagesByPath = make(map[string]*ssa.Package)

	for _, pkg := range ssaPkgs {
		if pkg != nil && pkg.Pkg != nil {
			i.packagesByPath[pkg.Pkg.Path()] = pkg
		}
	}

	return nil
}

func (i *Interp) LookupFunction(
	packagePath string,
	functionName string,
) (*ssa.Function, error) {
	if i.prog == nil {
		return nil, errors.New("解释器尚未加载程序")
	}

	pkg := i.packagesByPath[packagePath]
	if pkg == nil {
		return nil, fmt.Errorf("未找到包: %s", packagePath)
	}

	fn := pkg.Func(functionName)
	if fn == nil {
		return nil, fmt.Errorf(
			"包 %s 中未找到函数 %s",
			packagePath,
			functionName,
		)
	}

	return fn, nil
}
