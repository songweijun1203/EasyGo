package xinterp

import (
	"go/types"
	"reflect"
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
		//结构体的零值用一个map表达，key是字段名，value是字段值的地址
		st := make(map[string]*Value)
		for i := 0; i < u.NumFields(); i++ {
			field := u.Field(i)
			box := new(Value)
			//递归调用自身
			*box = zeroOf(field.Type())
			st[field.Name()] = box
		}
		return st

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

func isNil(v Value) bool {
	if v == nil {
		return true
	}
	// 使用 reflect 检查引用类型是否为 nil
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		return rv.IsNil()
	}
	return false
}
