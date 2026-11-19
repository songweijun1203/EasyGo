package xinterp

import (
	"fmt"
	"go/types"
)

func zeroOf(t types.Type) Value {
	// 获取底层类型（Underlying 会自动穿透 type MyInt int 这种自定义别名）
	switch u := t.Underlying().(type) {
	case *types.Basic:
		info := u.Info()
		switch {
		case info&types.IsBoolean != 0:
			return false
		case info&types.IsInteger != 0:
			// 【关键】配合你的 BinOp，所有整型统一返回 Go 的 int(0)
			// 这样后面的 x.(int) 断言才不会报错
			return int(0)
		case info&types.IsFloat != 0:
			// 所有浮点统一返回 float64(0.0)
			return float64(0.0)
		case info&types.IsString != 0:
			return ""
		}

	case *types.Pointer, *types.Map, *types.Chan, *types.Signature, *types.Interface:
		// 引用类型（指针、map、通道、函数、接口）的零值统一是 nil
		return nil

	case *types.Struct:
		// 如果你的玩具解释器还没实现结构体，这里先 panic 阻断
		// 如果实现了，这里需要构造一个包含所有字段零值的结构体 map/slice
		panic(fmt.Sprintf("玩具解释器暂不支持结构体零值分配: %v", t))

	case *types.Array:
		//panic(fmt.Sprintf("玩具解释器暂不支持数组零值分配: %v", t))
		arr := make([]*Value, u.Len())
		for i := 0; i < int(u.Len()); i++ {
			box := new(Value)
			*box = zeroOf(u.Elem())
			arr[i] = box
		}
		return arr
	case *types.Slice:
		// 切片的零值是 nil
		return [](*Value)(nil)
	default:
		// 兜底：如果遇到无法识别的类型，返回 nil
		return nil
	}
	return nil
}
