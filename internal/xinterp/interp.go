package xinterp

import "golang.org/x/tools/go/ssa"

type Interp struct {
	globals map[*ssa.Global]*Value // 全程序唯一，所有帧共享
}

func (i *Interp) RunFunction(fn *ssa.Function, args ...Value) Value {
	f := &Frame{
		interp: i,
		fn:     fn,
		block:  fn.Blocks[0],
		env:    make(map[ssa.Value]Value),
	}

	// 绑定参数到环境变量中
	for i, param := range fn.Params {
		f.env[param] = args[i]
	}

	f.run() // 执行当前帧的指令

	return f.result // 返回函数执行的结果
}
