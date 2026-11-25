package xinterp

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/ssa"
)

/*
定义 RuntimeMap 和 MapKey
修改 zeroOf，让 map 零值保留 key/value 类型
实现 ssa.MakeMap
实现 ssa.MapUpdate
实现 ssa.Lookup
扩展 Extract，确认 tuple 读取正常
扩展 len 支持 RuntimeMap
实现 delete
扩展 nativeValue 支持 map
最后再考虑 Range 和 Next
*/
type Value = interface{} // 简化：基本类型直接装

// 这是一个用于解释执行的帧结构体，包含了当前函数、基本块、环境变量和局部变量等信息。
type Frame struct {
	interp     *Interp             //全局唯一
	fn         *ssa.Function       // 当前正在执行的函数
	block      *ssa.BasicBlock     // 当前正在执行的基本块
	prevBlock  *ssa.BasicBlock     // 上一个基本块，用于处理控制流跳转
	env_values map[ssa.Value]Value // SSA 值到普通运行时值的映射
	env_slots  map[ssa.Value]*Cell // 表示可寻址位置的 SSA 值到存储槽的映射
	result     Value               // 当前帧的返回值，存储了函数执行后的结果
}

func NewFrame(i *Interp, fn *ssa.Function) *Frame {
	f := &Frame{
		interp:     i,
		fn:         fn,
		block:      fn.Blocks[0],
		env_values: make(map[ssa.Value]Value),
		env_slots:  make(map[ssa.Value]*Cell),
	}
	return f
}

// Closure 表示一个闭包对象
type Closure struct {
	Fn       *ssa.Function // 闭包关联的底层函数
	FreeVars []*Cell       // 捕获的自由变量（核心：存的是 *Cell 指针，实现状态共享）
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
		f.env_values[phi] = edgeValues[i]
	}
}

// execute 执行一条 SSA 指令，并返回当前基本块是否已经结束。
// Phi 不在这里执行；进入基本块时由 executePhis 统一处理。
func (f *Frame) execute(instr ssa.Instruction) bool {
	switch instr := instr.(type) {
	// ===== 控制流：会改变当前 Frame.block =====
	case *ssa.Jump:
		// 无条件跳到当前块唯一的后继块。
		f.block = instr.Block().Succs[0]
		return true
	case *ssa.If:
		// 根据条件选择两个后继块中的一个。
		if f.getVarValue(instr.Cond).(bool) {
			f.block = instr.Block().Succs[0]
		} else {
			f.block = instr.Block().Succs[1]
		}
		return true
	case *ssa.Return:
		// 保存返回值，并将 block 置空，让 run 结束当前函数。
		if len(instr.Results) == 1 {
			f.result = f.getVarValue(instr.Results[0])
		} else if len(instr.Results) > 1 {
			tuple := make([]Value, len(instr.Results))
			for i, result := range instr.Results {
				tuple[i] = f.getVarValue(result)
			}
			f.result = tuple
		}
		f.block = nil
		return true

	// ===== 值计算：生成或变换普通 SSA 值 =====
	case *ssa.BinOp:
		f.env_values[instr] = f.evalBinOp(instr.Op, f.getVarValue(instr.X), f.getVarValue(instr.Y))
		return false
	case *ssa.UnOp:
		f.env_values[instr] = f.evealUnOp(instr)
		return false
	case *ssa.Convert:
		val := f.getVarValue(instr.X)
		targetType := instr.Type().Underlying().(*types.Basic)

		switch targetType.Info() {
		case types.IsInteger:
			switch v := val.(type) {
			case int:
				f.env_values[instr] = v
			case float64:
				f.env_values[instr] = int(v)
			default:
				panic(fmt.Sprintf("cannot convert %T to int", val))
			}
		case types.IsFloat:
			switch v := val.(type) {
			case int:
				f.env_values[instr] = float64(v)
			case float64:
				f.env_values[instr] = v
			default:
				panic(fmt.Sprintf("cannot convert %T to float64", val))
			}
		default:
			panic(fmt.Sprintf("unsupported Convert target type: %v", targetType))
		}
		return false
	case *ssa.MakeInterface:
		// 当前解释器直接使用接口中的底层运行时值。
		f.env_values[instr] = f.getVarValue(instr.X)
		return false
	case *ssa.Extract:
		// 从多返回值调用产生的 tuple 中取出指定位置的值。
		tupleVal := f.getVarValue(instr.Tuple)
		if tuple, ok := tupleVal.([]Value); ok {
			f.env_values[instr] = tuple[instr.Index]
		} else {
			panic(fmt.Sprintf("Extract 操作数不是 Tuple ([]Value), 而是 %T", tupleVal))
		}
		return false

	// ===== 存储与地址：操作 Cell 或产生可寻址位置 =====
	case *ssa.Alloc:
		// 为局部变量分配 Cell，并写入该变量类型的零值。
		box := new(Cell)
		elemType := instr.Type().Underlying().(*types.Pointer).Elem()
		box.Value = zeroOf(elemType)
		f.env_slots[instr] = box
		return false
	case *ssa.Store:
		// 将普通运行时值写入目标 Cell。
		addr := f.getVarAddr(instr.Addr)
		addr.Value = f.getVarValue(instr.Val)
		return false
	case *ssa.FieldAddr:
		// 返回结构体字段对应的 Cell，供后续 Store 使用。
		addr := f.getVarAddr(instr.X)
		st, ok := addr.Value.(*RuntimeStruct)
		if !ok {
			panic(fmt.Sprintf("FieldAddr 操作数不是 RuntimeStruct, 而是 %T", addr.Value))
		}
		ptrType := instr.X.Type().Underlying().(*types.Pointer)
		structType := ptrType.Elem().Underlying().(*types.Struct)
		fieldName := structType.Field(instr.Field).Name()
		f.env_slots[instr] = st.Fields[fieldName]
		return false
	case *ssa.IndexAddr:
		// 返回切片元素对应的 Cell，保证元素赋值可以回写原切片。
		index := f.getVarValue(instr.Index).(int)
		collection, ok := f.getVarValue(instr.X).(*RuntimeSlice)
		if !ok || collection == nil {
			panic(fmt.Sprintf("不支持对 %T 取元素地址", f.getVarValue(instr.X)))
		}
		f.env_slots[instr] = collection.Elements[index]
		return false

	// ===== 聚合值与索引：处理结构体、切片和元素读取 =====
	case *ssa.Field:
		structVal := f.getVarValue(instr.X)
		st, ok := structVal.(*RuntimeStruct)
		if !ok {
			panic(fmt.Sprintf("Field 操作数不是 RuntimeStruct, 而是 %T", structVal))
		}
		structType := instr.X.Type().Underlying().(*types.Struct)
		fieldName := structType.Field(instr.Field).Name()
		f.env_values[instr] = readCell(st.Fields[fieldName])
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
		f.env_values[instr] = runtimeSlice
		return false
	case *ssa.Index:
		x := f.getVarValue(instr.X)
		i := f.getVarValue(instr.Index).(int)
		switch v := x.(type) {
		case string:
			f.env_values[instr] = v[i]
		case *RuntimeSlice:
			f.env_values[instr] = readCell(v.Elements[i])
		default:
			panic(fmt.Sprintf("不支持对 %T 进行索引", x))
		}
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
				f.env_values[instr] = values
				return false
			}
			if high < 0 {
				high = len(values.Elements)
			}
			f.env_values[instr] = &RuntimeSlice{
				ElemType: values.ElemType,
				Elements: values.Elements[low:high],
				Capacity: values.Cap() - low,
			}
		case string:
			if high < 0 {
				high = len(values)
			}
			f.env_values[instr] = values[low:high]
		default:
			panic(fmt.Sprintf("不支持对 %T 进行切片", collection))
		}
		return false

	// ===== Map：创建、读取和更新 RuntimeMap =====
	case *ssa.MakeMap:
		mapType, ok := instr.Type().Underlying().(*types.Map)
		if !ok {
			panic(fmt.Sprintf("MakeMap 类型不是 map: %v", instr.Type()))
		}
		reserve := 0
		if instr.Reserve != nil {
			reserve = f.getVarValue(instr.Reserve).(int)
			if reserve < 0 {
				reserve = 0
			}
		}

		runtimeMap := &RuntimeMap{
			KeyType:  mapType.Key(),
			ElemType: mapType.Elem(),
			Entries:  make(map[MapKey]*Cell, reserve),
			Nil:      false,
		}
		f.env_values[instr] = runtimeMap
		return false
	case *ssa.Lookup:
		mapVal := f.getVarValue(instr.X)
		keyVal := f.getVarValue(instr.Index)
		rm, ok := mapVal.(*RuntimeMap)
		if !ok {
			panic(fmt.Sprintf("Lookup 操作数不是 *RuntimeMap, 而是 %T", mapVal))
		}
		if rm == nil || rm.Nil {
			if instr.CommaOk {
				elemType := instr.Type().(*types.Tuple).At(0).Type()
				f.env_values[instr] = []Value{zeroOf(elemType), false}
			} else {
				f.env_values[instr] = zeroOf(instr.Type())
			}
			return false
		}

		mk := makeMapKey(rm, keyVal)
		cell, found := rm.Entries[mk]
		if instr.CommaOk {
			if found {
				f.env_values[instr] = []Value{readCell(cell), true}
			} else {
				elemType := instr.Type().(*types.Tuple).At(0).Type()
				f.env_values[instr] = []Value{zeroOf(elemType), false}
			}
		} else if found {
			f.env_values[instr] = readCell(cell)
		} else {
			f.env_values[instr] = zeroOf(instr.Type())
		}
		return false
	case *ssa.MapUpdate:
		mapVal := f.getVarValue(instr.Map)
		keyVal := f.getVarValue(instr.Key)
		valVal := f.getVarValue(instr.Value)
		rm, ok := mapVal.(*RuntimeMap)
		if !ok {
			panic(fmt.Sprintf("MapUpdate 操作数不是 *RuntimeMap, 而是 %T", mapVal))
		}
		// Go 语义：对 nil map 写入会 panic。
		if rm == nil || rm.Nil {
			panic("assignment to entry in nil map")
		}
		rm.Entries[makeMapKey(rm, keyVal)] = &Cell{Value: valVal}
		return false

	// ===== 函数与闭包：执行调用或创建闭包值 =====
	case *ssa.Call:
		return f.executeCall(instr)
	case *ssa.MakeClosure:
		// 捕获自由变量的 Cell，而不是复制自由变量当前的值。
		cells := make([]*Cell, len(instr.Bindings))
		for i, binding := range instr.Bindings {
			cells[i] = f.getVarAddr(binding)
		}

		fn, ok := instr.Fn.(*ssa.Function)
		if !ok {
			panic(fmt.Sprintf("MakeClosure 的目标不是函数: %T", instr.Fn))
		}
		f.env_values[instr] = &Closure{
			Fn:       fn,
			FreeVars: cells,
		}
		return false

	default:
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
func (f *Frame) evealUnOp(instr *ssa.UnOp) Value {
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
	//callee的可能取值*ssa.Builtin 、*ssa.Function、 *ssa.MakeClosure、*ssa.UnOp接口取出的闭包
	//先判断是不是内建函数
	if builtin, ok := callee.(*ssa.Builtin); ok {
		return f.executeBuiltin(instr, builtin)
	}

	// 静态函数直接使用 SSA 中的函数对象；动态函数值（包括闭包参数）
	// 必须先从当前帧读取，可能位于 env_values 或 env_slots 中。
	var fn *ssa.Function
	if staticFn, ok := callee.(*ssa.Function); ok {
		fn = staticFn
	} else {
		runtimeCallee := f.getVarValue(callee)
		if closure, ok := runtimeCallee.(*Closure); ok {
			return f.executeClosureCall(instr, closure)
		}

		var ok bool
		fn, ok = runtimeCallee.(*ssa.Function)
		if !ok {
			panic(fmt.Sprintf("暂时不支持调用非普通函数: %T", runtimeCallee))
		}
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

	//fmt.Printf("DEBUG: fnKey=[%s], fn.String()=[%s]\n", fnKey, fn.String())

	if _, isRegistered := nativeRegistry[fnKey]; isRegistered {
		return f.executeExternalCall(instr, fn, args)
	}

	if len(fn.Blocks) == 0 {
		panic(fmt.Sprintf("未注册的外部函数: %s", fn.String()))
	}

	result := f.interp.RunFunction(fn, args...) // 递归调用，创建一个新的Frame来执行被调用的函数

	// 4. 存结果
	if fn.Signature.Results().Len() > 0 {
		f.env_values[instr] = result
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
			f.env_values[instr] = len(v)
		case *RuntimeSlice:
			f.env_values[instr] = v.Len()
		case *RuntimeMap:
			f.env_values[instr] = v.Len()
		default:
			panic(fmt.Sprintf("len 不支持的类型: %T", args[0]))
		}

	case "cap":
		if len(args) != 1 {
			panic("cap 只能接受 1 个参数")
		}
		if slice, ok := args[0].(*RuntimeSlice); ok {
			f.env_values[instr] = slice.Cap()
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
				f.env_values[instr] = slice.appendSlice(spread)
				return false
			}
		}
		if len(args) > 1 {
			f.env_values[instr] = slice.appendValues(args[1:])
		} else {
			f.env_values[instr] = slice
		}
		if _, ok := f.env_values[instr].(*RuntimeSlice); !ok {
			panic(fmt.Sprintf("append 不支持的切片类型: %T", slice))
		}
	case "delete":
		if len(args) != 2 {
			panic("delete 需要两个参数")
		}

		runtimeMap, ok := args[0].(*RuntimeMap)
		if !ok {
			panic(fmt.Sprintf(
				"delete 第一个参数不是 *RuntimeMap，而是 %T",
				args[0],
			))
		}
		if runtimeMap == nil || runtimeMap.Nil {
			return false
		}

		key := makeMapKey(runtimeMap, args[1])
		delete(runtimeMap.Entries, key)
		return false
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
		f.env_values[instr] = out[0].Interface()
	} else if fn.Signature.Results().Len() > 1 {
		// 多返回值：ixgo 是包装成 Tuple，我们这里为了简化，先直接报错提醒
		//panic("暂不支持多返回值的外部函数（如 fmt.Println 返回的 n, err）")
		tuple := make([]Value, len(out))
		for i, v := range out {
			tuple[i] = v.Interface()
		}
		f.env_values[instr] = tuple
	}

	return false
}

// executeClosureCall 专门处理闭包的调用
func (f *Frame) executeClosureCall(instr *ssa.Call, closure *Closure) bool {
	fn := closure.Fn

	// 1. 收集调用时传入的显式参数
	var args []Value
	for _, arg := range instr.Common().Args {
		args = append(args, f.getVarValue(arg))
	}

	// 2. 创建新的执行帧
	newF := &Frame{
		interp:     f.interp,
		fn:         fn,
		block:      fn.Blocks[0],
		env_values: make(map[ssa.Value]Value),
		env_slots:  make(map[ssa.Value]*Cell),
	}

	// 3. 绑定自由变量 (FreeVars)
	// 在 SSA 中，闭包函数的自由变量会按顺序记录在 fn.FreeVars 列表中
	// 我们直接把闭包捕获的 *Cell 指针赋给新帧的 env，实现完美的状态共享！
	for i, freeVarSSA := range fn.FreeVars {
		newF.env_slots[freeVarSSA] = closure.FreeVars[i]
	}

	// 4. 绑定普通参数 (Params)
	for i, param := range fn.Params {
		box := new(Cell)
		if i < len(args) {
			box.Value = args[i]
		} else {
			box.Value = zeroOf(param.Type())
		}
		newF.env_slots[param] = box
	}

	// 5. 执行闭包函数
	newF.run()

	// 6. 将结果存回当前帧的 env
	if fn.Signature.Results().Len() > 0 {
		f.env_values[instr] = newF.result
	}

	return false // 闭包调用后继续执行当前帧
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

	// 普通 SSA 结果优先从值表读取。
	if value, ok := f.env_values[v]; ok {
		return value
	}

	// 可寻址 SSA 值从槽位表读取槽中的当前值。
	if slot, ok := f.env_slots[v]; ok {
		return slot.Value
	}

	panic(fmt.Sprintf("未找到 SSA 值: %T", v))
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

	if cell, ok := f.env_slots[v]; ok {
		return cell
	}
	panic(fmt.Sprintf("无法取地址: %T", v))
}
