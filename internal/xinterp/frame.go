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
	return (&Interp{globals: make(map[*ssa.Global]*Cell)}).RunFunction(fn, args...)
}

// 执行当前帧的指令，直到遇到返回指令或异常。
func (f *Frame) run() {
	for f.block != nil {
		//1 分类收集指令
		var phis []*ssa.Phi
		var otherInstrs []ssa.Instruction
		{
			for _, instr := range f.block.Instrs {
				phi, ok := instr.(*ssa.Phi)
				if ok {
					phis = append(phis, phi)
				} else {
					otherInstrs = append(otherInstrs, instr)
				}
			}
		}

		//2 先执行开头的全部phi指令
		f.executePhis(phis)

		//3 执行别的指令
		{
			for _, instr := range otherInstrs {
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

func (f *Frame) executePhis(phis []*ssa.Phi) {
	var edgeValues []Value
	for _, phi := range phis {
		//找到从前驱块来的那条边对应的值
		var val Value
		found := false
		for i, pred := range phi.Block().Preds {
			if pred == f.prevBlock {
				val = f.getVarValue(phi.Edges[i]) // 此时读到的绝对是旧值
				found = true
				break
			}
		}

		if !found && f.prevBlock != nil {
			// 初始进入 entry 块时 prevBlock 为 nil，跳过校验
			panic("Phi 找不到对应的前驱块")
		}
		edgeValues = append(edgeValues, val)
	}

	for i, phi := range phis {
		f.env[phi] = edgeValues[i]
	}
}

// 执行单条指令，根据指令类型进行不同的处理。
func (f *Frame) execute(instr ssa.Instruction) bool {
	switch instr := instr.(type) {
	case *ssa.BinOp:
		// 处理二元操作指令
		f.env[instr] = f.evalBinOp(instr.Op, f.getVarValue(instr.X), f.getVarValue(instr.Y))
		return false //普通指令，继续执行
	case *ssa.UnOp:
		// 处理一元操作指令
		// 常见情况：取地址、取值、逻辑非等
		f.env[instr] = f.evealUnop(instr)
		return false
	case *ssa.Return:
		// 处理返回指令
		if len(instr.Results) == 1 {
			f.result = f.getVarValue(instr.Results[0]) // 修复 2：直接读 env
		} else if len(instr.Results) > 1 {
			tuple := make([]Value, len(instr.Results))
			for i, r := range instr.Results {
				tuple[i] = f.getVarValue(r)
			}
			f.result = tuple
		}
		f.block = nil // 置空，主循环退出
		return true
	case *ssa.Extract:
		tupleVal := f.getVarValue(instr.Tuple)
		if tuple, ok := tupleVal.([]Value); ok {
			f.env[instr] = tuple[instr.Index]
		} else {
			panic(fmt.Sprintf("Extract 操作数不是 Tuple ([]Value), 而是 %T", tupleVal))
		}
		return false
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
		box := new(Cell)
		elemType := instr.Type().Underlying().(*types.Pointer).Elem()
		box.Value = zeroOf(elemType)
		f.env[instr] = box
		return false
	case *ssa.Store:
		// 处理存储指令
		addr := f.getVarAddr(instr.Addr)
		val := f.getVarValue(instr.Val)
		addr.Value = val
		return false
	case *ssa.MakeInterface:
		// 接口值在解释器中直接使用其底层值表示。
		f.env[instr] = f.getVarValue(instr.X)
		return false
	case *ssa.MakeSlice:
		length := f.getVarValue(instr.Len).(int)
		capacity := f.getVarValue(instr.Cap).(int)
		sliceType := instr.Type().Underlying().(*types.Slice)
		runtimeSlice := &RuntimeSlice{
			ElemType: sliceType.Elem(),
			Capacity: capacity,
		}
		for i := 0; i < length; i++ {
			runtimeSlice.Elements = append(runtimeSlice.Elements, &Cell{Value: zeroOf(sliceType.Elem())})
		}
		f.env[instr] = runtimeSlice
		return false
	case *ssa.IndexAddr:
		index := f.getVarValue(instr.Index).(int)
		collection, ok := f.getVarValue(instr.X).(*RuntimeSlice)
		if !ok || collection == nil {
			panic(fmt.Sprintf("不支持对 %T 取元素地址", f.getVarValue(instr.X)))
		}
		f.env[instr] = collection.Elements[index]
		return false

	case *ssa.Slice:
		collection := f.getVarValue(instr.X)
		low, high := 0, -1
		if instr.Low != nil {
			low = f.getVarValue(instr.Low).(int)
		}
		if instr.High != nil {
			high = f.getVarValue(instr.High).(int)
		}
		switch values := collection.(type) {
		case *RuntimeSlice:
			if values == nil || values.Nil {
				f.env[instr] = values
				return false
			}
			if high < 0 {
				high = len(values.Elements)
			}
			f.env[instr] = &RuntimeSlice{
				ElemType: values.ElemType,
				Elements: values.Elements[low:high],
				Capacity: values.Cap() - low,
			}
		case string:
			if high < 0 {
				high = len(values)
			}
			f.env[instr] = values[low:high]
		default:
			panic(fmt.Sprintf("不支持对 %T 进行切片", collection))
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
		case *RuntimeSlice:
			f.env[instr] = readCell(v.Elements[i])
		default:
			panic(fmt.Sprintf("不支持对 %T 进行索引", x))
		}
		return false
	case *ssa.Convert:
		val := f.getVarValue(instr.X)
		targetType := instr.Type().Underlying().(*types.Basic)

		switch targetType.Info() {
		case types.IsInteger:
			switch v := val.(type) {
			case int:
				f.env[instr] = v
			case float64:
				f.env[instr] = int(v) // float 转 int，截断小数
			default:
				panic(fmt.Sprintf("cannot convert %T to int", val))
			}
		case types.IsFloat:
			switch v := val.(type) {
			case int:
				f.env[instr] = float64(v) // int 转 float
			case float64:
				f.env[instr] = v
			default:
				panic(fmt.Sprintf("cannot convert %T to float64", val))
			}
		default:
			panic(fmt.Sprintf("unsupported Convert target type: %v", targetType))
		}
		return false
	case *ssa.Field: //结构体字段的值
		// 获取结构体字段值。
		structVal := f.getVarValue(instr.X)
		st, ok := structVal.(*RuntimeStruct)
		if !ok {
			panic(fmt.Sprintf("Field 操作数不是 RuntimeStruct, 而是 %T", structVal))
		}
		structType := instr.X.Type().Underlying().(*types.Struct)
		fieldName := structType.Field(instr.Field).Name()
		f.env[instr] = readCell(st.Fields[fieldName])
		return false
	case *ssa.FieldAddr: //结构体字段的地址
		addr := f.getVarAddr(instr.X)
		st, ok := addr.Value.(*RuntimeStruct)
		if !ok {
			panic(fmt.Sprintf("FieldAddr 操作数不是 RuntimeStruct, 而是 %T", addr.Value))
		}
		ptrType := instr.X.Type().Underlying().(*types.Pointer)
		structType := ptrType.Elem().Underlying().(*types.Struct)
		fieldName := structType.Field(instr.Field).Name()
		f.env[instr] = st.Fields[fieldName]
		return false
	default:
		// 处理其他类型的指令
		panic(fmt.Sprintf("unsupported instruction: %T", instr))
	}
}

// 专门处理二元运算
func (f *Frame) evalBinOp(op token.Token, x, y Value) Value {
	if op == token.EQL || op == token.NEQ {
		if isNil(x) || isNil(y) {
			// 如果两边都是 nil，或者一边是 nil 一边不是
			res := (isNil(x) && isNil(y))
			if op == token.NEQ {
				res = !res
			}
			return res
		}
	}
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
	case string:
		yv := y.(string)
		switch op {
		case token.ADD:
			return xv + yv
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

// 专门处理一元运算
func (f *Frame) evealUnop(instr *ssa.UnOp) Value {
	switch instr.Op {
	case token.MUL: // 取值操作
		addr := f.getVarAddr(instr.X)
		return addr.Value
	case token.NOT: //逻辑非
		return !f.getVarValue(instr.X).(bool)
	case token.SUB: //取负
		x := f.getVarValue(instr.X)
		switch v := x.(type) {
		case int:
			return -v
		case float64:
			return -v
		default:
			panic(fmt.Sprintf("unsupported type for SUB: %T", x))
		}
	case token.XOR: //按位取反
		return ^f.getVarValue(instr.X).(int)
	}
	panic(fmt.Sprintf("unsupported UnOp: %v", instr.Op))
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
		case *RuntimeSlice:
			f.env[instr] = v.Len()
		default:
			panic(fmt.Sprintf("len 不支持的类型: %T", args[0]))
		}

	case "cap":
		if len(args) != 1 {
			panic("cap 只能接受 1 个参数")
		}
		if slice, ok := args[0].(*RuntimeSlice); ok {
			f.env[instr] = slice.Cap()
		} else {
			panic(fmt.Sprintf("cap 不支持的类型: %T", args[0]))
		}
	case "append":
		//append一般两种用法
		//append(s,1,2,3)
		//t :=int{4,5}
		//s=append(s,t...)
		if len(args) < 1 {
			panic("append 至少需要一个参数（切片本身）")
		}
		var slice *RuntimeSlice
		if args[0] != nil {
			var ok bool
			slice, ok = args[0].(*RuntimeSlice)
			if !ok {
				panic(fmt.Sprintf("append 不支持的切片类型: %T", args[0]))
			}
		}
		if slice == nil {
			sliceType := instr.Type().Underlying().(*types.Slice)
			slice = &RuntimeSlice{ElemType: sliceType.Elem()}
		}

		if len(args) == 2 {
			if spread, ok := args[1].(*RuntimeSlice); ok {
				f.env[instr] = slice.appendSlice(spread)
				return false
			}
		}
		if len(args) > 1 {
			f.env[instr] = slice.appendValues(args[1:])
		} else {
			f.env[instr] = slice
		}
		if _, ok := f.env[instr].(*RuntimeSlice); !ok {
			panic(fmt.Sprintf("append 不支持的切片类型: %T", slice))
		}

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

	// 将解释器值物化成原生函数签名要求的 Go 值。
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = nativeValue(arg, nativeFunc.Type().In(i))
	}

	var out []reflect.Value
	if nativeFunc.Type().IsVariadic() {
		// SSA 已把可变参数放在 Args 的最后一个切片中。
		out = nativeFunc.CallSlice(in)
	} else {
		out = nativeFunc.Call(in)
	}

	// 4. 将返回值塞回 SSA 的 env 中
	if fn.Signature.Results().Len() == 1 {
		f.env[instr] = out[0].Interface()
	} else if fn.Signature.Results().Len() > 1 {
		// 多返回值：ixgo 是包装成 Tuple，我们这里为了简化，先直接报错提醒
		//panic("暂不支持多返回值的外部函数（如 fmt.Println 返回的 n, err）")
		tuple := make([]Value, len(out))
		for i, v := range out {
			tuple[i] = v.Interface()
		}
		f.env[instr] = tuple
	}

	return false
}

// 获取变量值
func (f *Frame) getVarValue(v ssa.Value) Value {
	if c, ok := v.(*ssa.Const); ok {
		if c.Value == nil {
			return zeroOf(c.Type())
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
		return box.Value
	}

	//局部变量
	return readCell(f.env[v])
}

// 获取全局变量地址
func (f *Frame) getGlobalAddr(g *ssa.Global) *Cell {
	if box, ok := f.interp.globals[g]; ok {
		return box
	}

	//懒初始化：第一次访问，分配该类型的零值
	elemType := g.Type().Underlying().(*types.Pointer).Elem()
	box := &Cell{Value: zeroOf(elemType)}
	f.interp.globals[g] = box // 存入解释器的全局表
	return box
}

// 获取变量地址
func (f *Frame) getVarAddr(v ssa.Value) *Cell {
	if g, ok := v.(*ssa.Global); ok {
		return f.getGlobalAddr(g)
	}

	if cell, ok := f.env[v].(*Cell); ok {
		return cell
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
