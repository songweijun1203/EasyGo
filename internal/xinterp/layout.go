package xinterp

import "golang.org/x/tools/go/ssa"

// blockLayout 是一个基本块的只读指令划分。
// Phi 必须在进入块时并行求值，其余指令顺序执行，这个划分是块的纯函数。
type blockLayout struct {
	phis   []*ssa.Phi
	others []ssa.Instruction
}

// SSA 函数的只读执行布局
type funcLayout struct {
	size     int               // 寄存器总数
	regIndex map[ssa.Value]int // SSA 值 -> 寄存器索引；构建完成后只读
	entry    *ssa.BasicBlock   // 入口块；外部函数（无 Blocks）为 nil

	blocks  []blockLayout // 与 fn.Blocks 同序，用 block.Index 直接索引
	maxPhis int           // 所有块中 phi 数量的最大值，用于预分配暂存区
}

func buildFuncLayout(fn *ssa.Function) *funcLayout {
	regIndex := make(map[ssa.Value]int)
	idx := 0

	// 参数和自由变量不在任何 Block.Instrs 里，必须单独分配。
	for _, param := range fn.Params {
		regIndex[param] = idx
		idx++
	}
	for _, freeVar := range fn.FreeVars {
		regIndex[freeVar] = idx
		idx++
	}

	//在遍历block的同时，把phi也一并筛选出来
	layout := &funcLayout{regIndex: regIndex}
	layout.blocks = make([]blockLayout, len(fn.Blocks))

	for bi, block := range fn.Blocks {
		var bl blockLayout
		for _, instr := range block.Instrs {
			if v, ok := instr.(ssa.Value); ok {
				regIndex[v] = idx
				idx++
			}
			if phi, ok := instr.(*ssa.Phi); ok {
				bl.phis = append(bl.phis, phi)
			} else {
				bl.others = append(bl.others, instr)
			}
		}
		if n := len(bl.phis); n > layout.maxPhis {
			layout.maxPhis = n
		}
		layout.blocks[bi] = bl
	}
	layout.size = idx
	if len(fn.Blocks) > 0 {
		layout.entry = fn.Blocks[0]
	}
	return layout
}

func (i *Interp) layoutOf(fn *ssa.Function) *funcLayout {
	i.layoutMu.RLock()
	layout := i.layouts[fn]
	i.layoutMu.RUnlock()
	if layout != nil {
		return layout
	}

	i.layoutMu.Lock()
	defer i.layoutMu.Unlock()
	// 等锁期间可能已有别的 goroutine 建好了。
	if layout = i.layouts[fn]; layout != nil {
		return layout
	}
	layout = buildFuncLayout(fn)
	if i.layouts == nil {
		i.layouts = make(map[*ssa.Function]*funcLayout)
	}
	i.layouts[fn] = layout
	return layout
}
