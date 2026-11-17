package xinterp

import (
	"fmt"
	"go/constant"
	"go/token"

	"golang.org/x/tools/go/ssa"
)

type Value = interface{} // 简化：基本类型直接装

// 这是一个用于解释执行的帧结构体，包含了当前函数、基本块、环境变量和局部变量等信息。
type Frame struct {
	fn        *ssa.Function       // 当前正在执行的函数
	block     *ssa.BasicBlock     // 当前正在执行的基本块
	prevBlock *ssa.BasicBlock     // 上一个基本块，用于处理控制流跳转
	env       map[ssa.Value]Value // 当前帧的环境变量映射，存储了 SSA 值到实际值的映射
	locals    []Value             // 当前帧的局部变量列表，存储了函数的局部变量值
	result    Value               // 当前帧的返回值，存储了函数执行后的结果
}

func RunFunction(fn *ssa.Function, args ...Value) Value {
	f := &Frame{
		fn:    fn,
		block: fn.Blocks[0],
		env:   make(map[ssa.Value]Value),
	}

	// 绑定参数到环境变量中
	for i, param := range fn.Params {
		f.env[param] = args[i]
	}

	f.run() // 执行当前帧的指令

	return f.result // 返回函数执行的结果
}

// 执行当前帧的指令，直到遇到返回指令或异常。
func (f *Frame) run() {
	for f.block != nil {
		for _, instr := range f.block.Instrs {
			if _, ok := instr.(*ssa.Phi); ok {
				f.executePhi(instr.(*ssa.Phi)) // 处理 Phi 指令
			} else {
				break // 如果遇到非 Phi 指令，跳出循环，继续执行后续指令
			}
		}

		for _, instr := range f.block.Instrs {
			if _, ok := instr.(*ssa.Phi); ok {
				continue // 跳过 Phi
			}

			// 解释执行每条指令
			oldBlock := f.block
			jumped := f.execute(instr)
			// 如果执行过程中遇到返回指令，退出循环
			if jumped {
				f.prevBlock = oldBlock // 更新上一个基本块
				break                  //跳出当前基本块的指令循环，继续执行下一条指令
			}
		}
	}
}

// 统一的取值方法（处理常量、内存指针、普通值）
func (f *Frame) get(v ssa.Value) Value {
	if c, ok := v.(*ssa.Const); ok {
		if c.Value == nil {
			return nil
		}
		typeStr := c.Type().String()
		switch typeStr {
		case "untyped int", "int", "int64", "int32": // 注意这里可以列举你关心的整数类型
			return int(c.Int64())
		case "untyped float", "float64", "float32":
			// c.Float64() 返回 float64
			val, _ := constant.Float64Val(c.Value)
			return val
		case "untyped bool", "bool":
			return constant.BoolVal(c.Value)
		}
		panic(fmt.Sprintf("unsupported const type: %v", c.Type()))
	}

	val := f.env[v]
	if ptr, ok := val.(*Value); ok {
		val = *ptr
	}
	return val
}

// 执行 Phi 指令，根据前一个基本块的值选择正确的分支。
func (f *Frame) executePhi(phi *ssa.Phi) {
	for i, pred := range phi.Block().Preds {
		if pred == f.prevBlock {
			f.env[phi] = f.get(phi.Edges[i])
			return
		}
	}
	panic("Phi 找不到对应的前驱块")
}

// 执行单条指令，根据指令类型进行不同的处理。
func (f *Frame) execute(instr ssa.Instruction) bool {
	switch instr := instr.(type) {
	case *ssa.BinOp:
		// 处理二元操作指令
		x := f.get(instr.X)
		y := f.get(instr.Y)
		// 统一调用计算函数，不管它是 int 还是 float
		f.env[instr] = f.evalBinOp(instr.Op, x, y)
		return false //普通指令，继续执行
	case *ssa.Return:
		// 处理返回指令
		if len(instr.Results) > 0 {
			f.result = f.get(instr.Results[0]) // 修复 2：直接读 env
		}
		f.block = nil // 置空，主循环退出
		return true
		// ====== 新增的控制流指令 ======
	case *ssa.Jump:
		// Jump 指令无条件跳转到第一个后继块
		f.block = instr.Block().Succs[0] // 修复 1：使用 Block().Succs[0] 获取后继块
		return true
	case *ssa.If:
		// If 指令根据条件跳转到两个后继块之一
		cond := f.get(instr.Cond).(bool)
		if cond {
			f.block = instr.Block().Succs[0] // 条件为真，跳转到第一个后继块
		} else {
			f.block = instr.Block().Succs[1] // 条件为假，跳转到第二个后继块
		}
		return true
	case *ssa.Alloc:
		// 处理内存分配指令
		var v Value
		f.env[instr] = &v
		return false
	case *ssa.Store:
		// 处理存储指令
		addr := f.env[instr.Addr].(*Value) // 用 env 拿地址
		val := f.get(instr.Val)
		*addr = val
		return false
	case *ssa.UnOp:
		// 处理一元操作指令
		// 常见情况：取地址、取值、逻辑非等
		switch instr.Op {
		case token.MUL: // 取值操作
			addr := f.env[instr.X].(*Value)
			f.env[instr] = *addr
		default:
			panic(fmt.Sprintf("unsupported UnOp: %v", instr.Op))
		}
		return false
	case *ssa.Call:
		// 处理函数调用指令
		return f.executeCall(instr)
	default:
		// 处理其他类型的指令
		panic(fmt.Sprintf("unsupported instruction: %T", instr))
	}
}

// 专门处理二元运算
func (f *Frame) evalBinOp(op token.Token, x, y Value) Value {
	switch xv := x.(type) {
	case int:
		yv := y.(int)
		switch op {
		case token.ADD:
			return xv + yv
		case token.SUB:
			return xv - yv
		case token.MUL:
			return xv * yv
		case token.QUO:
			return xv / yv
		case token.REM:
			return xv % yv
		case token.EQL:
			return xv == yv
		case token.NEQ:
			return xv != yv
		case token.LSS:
			return xv < yv
		case token.LEQ:
			return xv <= yv
		case token.GTR:
			return xv > yv
		case token.GEQ:
			return xv >= yv
		}
	case float64:
		yv := y.(float64)
		switch op {
		case token.ADD:
			return xv + yv
		case token.SUB:
			return xv - yv
		case token.MUL:
			return xv * yv
		case token.QUO:
			return xv / yv
		case token.EQL:
			return xv == yv
		case token.NEQ:
			return xv != yv
		case token.LSS:
			return xv < yv
		case token.LEQ:
			return xv <= yv
		case token.GTR:
			return xv > yv
		case token.GEQ:
			return xv >= yv
		}
	}
	panic(fmt.Sprintf("unsupported types for op %v: %T, %T", op, x, y))
}

// 处理函数调用
func (f *Frame) executeCall(instr *ssa.Call) bool {
	callee := instr.Common().Value
	fn, ok := callee.(*ssa.Function)
	if !ok {
		panic(fmt.Sprintf("暂时不支持调用非普通函数: %T", callee))
	}

	//收集参数
	var args []Value
	for _, arg := range instr.Common().Args {
		args = append(args, f.get(arg))
	}

	result := RunFunction(fn, args...) // 递归调用，创建一个新的Frame来执行被调用的函数

	// 4. 存结果
	if fn.Signature.Results().Len() > 0 {
		f.env[instr] = result
	}

	return false // 函数调用后继续执行当前帧
}

// 后续不需要
func RunAdd(fn *ssa.Function, a, b int) Value {
	f := &Frame{
		fn:    fn,
		block: fn.Blocks[0],
		env:   make(map[ssa.Value]Value),
	}
	// 绑定参数
	f.env[fn.Params[0]] = a
	f.env[fn.Params[1]] = b

	f.run()
	return f.result
}
