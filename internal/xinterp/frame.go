package xinterp

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/ssa"
)

type Value = interface{} // 简化：基本类型直接装

// 这是一个用于解释执行的帧结构体，包含了当前函数、基本块、环境变量和局部变量等信息。
type Frame struct {
	interp    *Interp             //全局唯一
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

// 执行 Phi 指令，根据前一个基本块的值选择正确的分支。
func (f *Frame) executePhi(phi *ssa.Phi) {
	for i, pred := range phi.Block().Preds {
		if pred == f.prevBlock {
			f.env[phi] = f.getVarValue(phi.Edges[i])
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
		x := f.getVarValue(instr.X)
		y := f.getVarValue(instr.Y)
		// 统一调用计算函数，不管它是 int 还是 float
		f.env[instr] = f.evalBinOp(instr.Op, x, y)
		return false //普通指令，继续执行
	case *ssa.Return:
		// 处理返回指令
		if len(instr.Results) > 0 {
			f.result = f.getVarValue(instr.Results[0]) // 修复 2：直接读 env
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
		cond := f.getVarValue(instr.Cond).(bool)
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
		//addr := f.env[instr.Addr].(*Value) // 用 env 拿地址
		addr := f.getVarAddr(instr.Addr)
		val := f.getVarValue(instr.Val)
		*addr = val
		return false
	case *ssa.UnOp:
		// 处理一元操作指令
		// 常见情况：取地址、取值、逻辑非等
		switch instr.Op {
		case token.MUL: // 取值操作
			//addr := f.env[instr.X].(*Value)
			addr := f.getVarAddr(instr.X)
			f.env[instr] = *addr

		default:
			panic(fmt.Sprintf("unsupported UnOp: %v", instr.Op))
		}
		return false
	case *ssa.Call:
		// 处理函数调用指令
		return f.executeCall(instr)
	case *ssa.Index:
		x := f.getVarValue(instr.X)           // 被索引的对象（字符串或数组）
		i := f.getVarValue(instr.Index).(int) // 索引下标
		switch v := x.(type) {
		case string:
			f.env[instr] = v[i]
		case []Value: // 如果你自己实现了切片
			f.env[instr] = v[i]
		default:
			panic(fmt.Sprintf("不支持对 %T 进行索引", x))
		}
		return false
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
	//先判断是不是内建函数
	if builtin, ok := callee.(*ssa.Builtin); ok {
		return f.executeBuiltin(instr, builtin)
	}

	fn, ok := callee.(*ssa.Function)
	if !ok {
		panic(fmt.Sprintf("暂时不支持调用非普通函数: %T", callee))
	}

	//收集参数
	var args []Value
	for _, arg := range instr.Common().Args {
		args = append(args, f.getVarValue(arg))
	}

	//FFI边界分流
	var fnKey string
	if fn.Pkg != nil && fn.Pkg.Pkg != nil {
		fnKey = fn.Pkg.Pkg.Path() + "." + fn.Name()
	} else {
		fnKey = fn.String() // 兜底
	}

	fmt.Printf("DEBUG: fnKey=[%s], fn.String()=[%s]\n", fnKey, fn.String())

	if _, isRegistered := nativeRegistry[fnKey]; isRegistered {
		return f.executeExternalCall(instr, fn, args)
	}

	if len(fn.Blocks) == 0 {
		panic(fmt.Sprintf("未注册的外部函数: %s", fn.String()))
	}

	result := f.interp.RunFunction(fn, args...) // 递归调用，创建一个新的Frame来执行被调用的函数

	// 4. 存结果
	if fn.Signature.Results().Len() > 0 {
		f.env[instr] = result
	}

	return false // 函数调用后继续执行当前帧
}

// 调用内建函数
func (f *Frame) executeBuiltin(instr *ssa.Call, builtin *ssa.Builtin) bool {
	// 1. 收集参数（复用原来的逻辑）
	var args []Value
	for _, arg := range instr.Common().Args {
		args = append(args, f.getVarValue(arg))
	}

	switch builtin.Name() {
	case "println":
		// 模拟 println 的行为，用 Go 的 fmt.Println
		fmt.Println(args...)
		// 无返回值，直接返回

	case "print":
		fmt.Print(args...)

	case "len":
		// 模拟 len 的行为
		if len(args) != 1 {
			panic("len 只能接受 1 个参数")
		}
		// 因为 args[0] 是 interface{}，我们用 type switch 来算长度
		switch v := args[0].(type) {
		case string:
			f.env[instr] = len(v)
		case []interface{}: // 预留：如果以后支持了切片
			f.env[instr] = len(v)
		default:
			panic(fmt.Sprintf("len 不支持的类型: %T", args[0]))
		}

	case "cap":
		// 类似 len 的实现
		f.env[instr] = 0 // 简化处理

	default:
		panic(fmt.Sprintf("不支持的 Builtin 函数: %s", builtin.Name()))
	}
	return false
}

// 调用外部函数
func (f *Frame) executeExternalCall(instr *ssa.Call, fn *ssa.Function, args []Value) bool {
	//1. 从注册表里查真实函数
	nativeFunc, ok := nativeRegistry[fn.String()]
	if !ok {
		panic(fmt.Sprintf("未注册的外部函数：%s", fn.String()))
	}

	//2. 将解释器的[]value转换为反射需要的[]reflect.Value
	var in []reflect.Value
	for _, arg := range args {
		// 注意：如果 arg 是 nil，reflect.ValueOf(nil) 会 panic，所以需要特殊处理
		if arg == nil {
			in = append(in, reflect.Zero(reflect.TypeOf((*interface{})(nil)).Elem()))
		} else {
			in = append(in, reflect.ValueOf(arg))
		}
	}

	//反射调用真实go函数
	out := nativeFunc.Call(in)

	// 4. 将返回值塞回 SSA 的 env 中
	if fn.Signature.Results().Len() == 1 {
		f.env[instr] = out[0].Interface()
	} else if fn.Signature.Results().Len() > 1 {
		// 多返回值：ixgo 是包装成 Tuple，我们这里为了简化，先直接报错提醒
		panic("暂不支持多返回值的外部函数（如 fmt.Println 返回的 n, err）")
	}

	return false
}

// 获取变量值
func (f *Frame) getVarValue(v ssa.Value) Value {
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
		case "untyped string", "string":
			return constant.StringVal(c.Value)
		}
		panic(fmt.Sprintf("unsupported const type: %v", c.Type()))
	}

	//是否是全局变量
	if g, ok := v.(*ssa.Global); ok {
		box := f.getGlobalAddr(g)
		return *box
	}

	//局部变量
	val := f.env[v]
	if ptr, ok := val.(*Value); ok {
		val = *ptr
	}
	return val
}

// 获取全局变量地址
func (f *Frame) getGlobalAddr(g *ssa.Global) *Value {
	if box, ok := f.interp.globals[g]; ok {
		return box
	}

	//懒初始化：第一次访问，分配该类型的零值
	elemType := g.Type().Underlying().(*types.Pointer).Elem()
	z := zeroOf(elemType)
	box := &z
	f.interp.globals[g] = box // 存入解释器的全局表
	return box
}

// 获取变量地址
func (f *Frame) getVarAddr(v ssa.Value) *Value {
	if g, ok := v.(*ssa.Global); ok {
		return f.getGlobalAddr(g)
	}

	val := f.env[v]
	if ptr, ok := val.(*Value); ok {
		return ptr
	}
	panic(fmt.Sprintf("无法取地址: %T", v))
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
