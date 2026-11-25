package xinterp

import (
	"fmt"
	"go/types"
	"reflect"
)

// Cell 表示一个可寻址的解释器存储槽；*Cell 是 SSA 指针值的运行时表示。
// 普通值直接保存在寄存器中，不需要经过 Cell。
type Cell struct {
	Value Value
}

// RuntimeSlice 是解释器内部的数组/切片表示。
// Elements 保存槽位，而不是把槽位指针直接暴露给原生 Go 函数。
type RuntimeSlice struct {
	ElemType types.Type
	Elements []*Cell
	Capacity int
	Nil      bool
}

// RuntimeStruct 是解释器内部的结构体表示。
type RuntimeStruct struct {
	Fields map[string]*Cell
}

// RuntimeMap 是解释器内部的map表示
type RuntimeMap struct {
	KeyType  types.Type
	ElemType types.Type
	Entries  map[MapKey]*Cell
	Nil      bool //用于表示map的零值，不能用普通的零值，应为还需要知道key/value类型
}

type RuntimeRange struct {
	Collection Value    //被迭代的集合
	Index      int      //当前迭代的索引
	Keys       []MapKey //仅用于Map,预提取的key列表，用于确定遍历顺序
}

func (m *RuntimeMap) Len() int {
	if m == nil || m.Nil {
		return 0
	}
	return len(m.Entries)
}

type MapKey struct {
	Type  string //类型，如int
	Value interface{}
}

func makeMapKey(runtimeMap *RuntimeMap, value Value) MapKey {
	value = readCell(value)

	switch value.(type) {
	case int, string, bool:
		return MapKey{
			Type:  runtimeMap.KeyType.String(),
			Value: value,
		}
	default:
		panic(fmt.Sprintf(
			"暂不支持 %T 作为 map key",
			value,
		))
	}
}

func readCell(value Value) Value {
	cell, ok := value.(*Cell)
	if !ok || cell == nil {
		return value
	}
	return cell.Value
}

func cloneCell(value Value) *Cell {
	return &Cell{Value: value}
}

func (s *RuntimeSlice) Len() int {
	if s == nil {
		return 0
	}
	return len(s.Elements)
}

func (s *RuntimeSlice) Cap() int {
	if s == nil {
		return 0
	}
	if s.Capacity < len(s.Elements) {
		return len(s.Elements)
	}
	return s.Capacity
}

func (s *RuntimeSlice) appendValues(values []Value) *RuntimeSlice {
	if s == nil {
		s = &RuntimeSlice{}
	}
	s.Nil = false
	for _, value := range values {
		s.Elements = append(s.Elements, cloneCell(value))
	}
	if s.Capacity < len(s.Elements) {
		s.Capacity = len(s.Elements)
	}
	return s
}

func (s *RuntimeSlice) appendSlice(other *RuntimeSlice) *RuntimeSlice {
	if other == nil {
		return s
	}
	for _, cell := range other.Elements {
		s = s.appendValues([]Value{readCell(cell)})
	}
	return s
}

func nativeMap(runtimeMap *RuntimeMap, target reflect.Type) reflect.Value {
	if runtimeMap == nil {
		return reflect.Zero(target)
	}

	// 目标是 interface{} 时，根据 RuntimeMap 的类型推导真实 map 类型
	if target.Kind() == reflect.Interface {
		target = reflect.MapOf(
			reflectTypeFor(runtimeMap.KeyType),
			reflectTypeFor(runtimeMap.ElemType),
		)
	}

	if target.Kind() != reflect.Map {
		panic(fmt.Sprintf(
			"不能将 RuntimeMap 转换为 %v",
			target,
		))
	}

	// typed nil map
	if runtimeMap.Nil {
		return reflect.Zero(target)
	}

	result := reflect.MakeMapWithSize(
		target,
		len(runtimeMap.Entries),
	)

	for key, cell := range runtimeMap.Entries {
		nativeKey := nativeValue(
			key.Value,
			target.Key(),
		)

		nativeValue := nativeValue(
			readCell(cell),
			target.Elem(),
		)

		result.SetMapIndex(nativeKey, nativeValue)
	}

	return result
}

func nativeValue(value Value, target reflect.Type) reflect.Value {
	value = readCell(value)
	if value == nil {
		return reflect.Zero(target)
	}

	switch value := value.(type) {
	case *RuntimeSlice:
		return nativeSlice(value, target)
	case *RuntimeMap:
		return nativeMap(value, target)
	case *RuntimeStruct:
		if target.Kind() != reflect.Interface {
			panic(fmt.Sprintf("不能将 RuntimeStruct 转换为 %v", target))
		}
		fields := make(map[string]interface{}, len(value.Fields))
		for name, cell := range value.Fields {
			fields[name] = nativeInterfaceValue(readCell(cell))
		}
		return reflect.ValueOf(fields)
	}

	result := reflect.ValueOf(value)
	if result.Type().AssignableTo(target) {
		return result
	}
	if result.Type().ConvertibleTo(target) {
		return result.Convert(target)
	}
	if target.Kind() == reflect.Interface && result.Type().Implements(target) {
		return result
	}
	panic(fmt.Sprintf("不能将解释器值 %T 转换为 %v", value, target))
}

func nativeSlice(slice *RuntimeSlice, target reflect.Type) reflect.Value {
	if target.Kind() == reflect.Interface {
		if slice == nil || slice.ElemType == nil {
			return reflect.Zero(target)
		}
		target = reflect.SliceOf(reflectTypeFor(slice.ElemType))
	}
	if slice == nil || slice.Nil {
		return reflect.Zero(target)
	}
	if target.Kind() != reflect.Slice {
		panic(fmt.Sprintf("不能将 RuntimeSlice 转换为 %v", target))
	}

	result := reflect.MakeSlice(target, len(slice.Elements), len(slice.Elements))
	for i, cell := range slice.Elements {
		result.Index(i).Set(nativeValue(readCell(cell), target.Elem()))
	}
	return result
}

func nativeInterfaceValue(value Value) interface{} {
	value = readCell(value)
	switch value := value.(type) {
	case *RuntimeSlice:
		return nativeSlice(value, reflect.SliceOf(reflectTypeFor(value.ElemType))).Interface()
	case *RuntimeStruct:
		fields := make(map[string]interface{}, len(value.Fields))
		for name, cell := range value.Fields {
			fields[name] = nativeInterfaceValue(readCell(cell))
		}
		return fields
	case *RuntimeMap:
		return nativeMap(
			value,
			reflect.TypeOf((*interface{})(nil)).Elem(),
		).Interface()
	default:
		return value
	}
}

func reflectTypeFor(t types.Type) reflect.Type {
	switch t := t.Underlying().(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool, types.UntypedBool:
			return reflect.TypeOf(false)
		case types.String, types.UntypedString:
			return reflect.TypeOf("")
		case types.Float32:
			return reflect.TypeOf(float32(0))
		case types.Float64, types.UntypedFloat:
			return reflect.TypeOf(float64(0))
		default:
			if t.Info()&types.IsInteger != 0 {
				return reflect.TypeOf(int(0))
			}
		}
	case *types.Slice:
		return reflect.SliceOf(reflectTypeFor(t.Elem()))
	case *types.Interface:
		return reflect.TypeOf((*interface{})(nil)).Elem()
	}
	return reflect.TypeOf((*interface{})(nil)).Elem()
}
