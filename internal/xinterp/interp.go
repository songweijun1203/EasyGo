package xinterp

import "golang.org/x/tools/go/ssa"

type Interp struct {
	globals map[*ssa.Global]*Cell // 全程序唯一，所有帧共享
}

func (i *Interp) RunFunction(fn *ssa.Function, args ...Value) Value {
	f := &Frame{
		interp:     i,
		fn:         fn,
		block:      fn.Blocks[0],
		env_values: make(map[ssa.Value]Value),
		env_slots:  make(map[ssa.Value]*Cell),
	}

	// 当前实现把参数统一放入槽位表，保证参数被取地址时仍然可写。
	for i, param := range fn.Params {
		box := new(Cell)
		if i < len(args) {
			box.Value = args[i] // 把传进来的实参放进去
		} else {
			box.Value = zeroOf(param.Type()) // 防止参数个数不匹配时的兜底
		}
		f.env_slots[param] = box
	}

	f.run() // 执行当前帧的指令

	return f.result // 返回函数执行的结果
}
