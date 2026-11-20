package xinterp

import "golang.org/x/tools/go/ssa"

type Interp struct {
	globals map[*ssa.Global]*Cell // 全程序唯一，所有帧共享
}

func (i *Interp) RunFunction(fn *ssa.Function, args ...Value) Value {
	f := &Frame{
		interp: i,
		fn:     fn,
		block:  fn.Blocks[0],
		env:    make(map[ssa.Value]Value),
	}

	// 绑定参数到环境变量中
	if false {
		for i, param := range fn.Params {
			f.env[param] = args[i]
		}
	}

	//如果函数内部，需要对参数取地址的话，此处直接存储值，没法取地址，会出问题
	for i, param := range fn.Params {
		box := new(Cell)
		if i < len(args) {
			box.Value = args[i] // 把传进来的实参放进去
		} else {
			box.Value = zeroOf(param.Type()) // 防止参数个数不匹配时的兜底
		}
		f.env[param] = box
	}

	f.run() // 执行当前帧的指令

	return f.result // 返回函数执行的结果
}
