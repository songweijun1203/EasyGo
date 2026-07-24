package utils

import (
	"fmt"
	"unsafe"
)

type XDataArray1D[T XNumber] struct {
	data      []T    //数据
	component int    //组份
	name      string //名称
}

func newXDataArray1DImp[T XNumber]() *XDataArray1D[T] {
	a := &XDataArray1D[T]{
		data:      nil,
		component: 1,
		name:      "",
	}
	return a
}

func (a *XDataArray1D[T]) SetName(n string) {
	a.name = n
}

func (a *XDataArray1D[T]) SetComponent(n int) {
	a.component = n
}

func (a *XDataArray1D[T]) GetName() string {
	return a.name
}

func (a *XDataArray1D[T]) GetSize() int {
	return len(a.data)
}

func (a *XDataArray1D[T]) GetComponent() int {
	return a.component
}

func (a *XDataArray1D[T]) GetBytesNum() int {
	var zero T
	return len(a.data) * int(unsafe.Sizeof(zero))
}

func (a *XDataArray1D[T]) SetNumOfTuple(n int) {
	if a.GetSize() != n*a.GetComponent() {
		//重新分配
		a.data = make([]T, n*a.GetComponent())
	}
}

func (a *XDataArray1D[T]) Clear() {
	a.data = nil
}

func (a *XDataArray1D[T]) SetTuple(index int, values ...T) {
	if len(values) != a.component {
		panic("values length does not match component")
	}

	if index < 0 || index >= a.GetTupleNum() {
		panic(fmt.Sprintf("index %d out of range [0, %d)", index, a.GetTupleNum()))
	}

	start := index * a.component
	for i, v := range values {
		a.data[start+i] = v
	}
}

func (a *XDataArray1D[T]) GetTupleNum() int {
	return len(a.data) / a.component
}

func (a *XDataArray1D[T]) GetTuple(index int) []T {
	if index < 0 || index >= a.GetTupleNum() {
		panic(fmt.Sprintf("index %d out of range [0, %d)", index, a.GetTupleNum()))
	}
	start := index * a.component
	end := start + a.component
	return a.data[start:end]
}

func (a *XDataArray1D[T]) GetTupleCopy(index int) []T {
	start := index * a.component
	end := start + a.component

	result := make([]T, a.component)
	copy(result, a.data[start:end])

	return result
}

func (a *XDataArray1D[T]) ForEach_tuple(op func(values []T)) {
	//循环遍历
	num := a.GetTupleNum()
	for i := 0; i < num; i++ {
		op(a.GetTuple(i))
	}
}

type (
	XUInt8Array1D  = XDataArray1D[uint8]
	XUInt16Array1D = XDataArray1D[uint16]
	XUInt32Array1D = XDataArray1D[uint32]
	XUInt64Array1D = XDataArray1D[uint64]

	XInt8Array1D  = XDataArray1D[int8]
	XInt16Array1D = XDataArray1D[int16]
	XInt32Array1D = XDataArray1D[int32]
	XInt64Array1D = XDataArray1D[int64]

	XFloat32Array1D = XDataArray1D[float32]
	XFloat64Array1D = XDataArray1D[float64]
)

// 将具体的构造函数赋值给别名，调用时就像调用一个类名一样！
var (
	NewXUInt8Array1D  = newXDataArray1DImp[uint8]
	NewXUInt16Array1D = newXDataArray1DImp[uint16]
	NewXUInt32Array1D = newXDataArray1DImp[uint32]
	NewXUInt64Array1D = newXDataArray1DImp[uint64]

	NewXInt8Array1D  = newXDataArray1DImp[int8]
	NewXInt16Array1D = newXDataArray1DImp[int16]
	NewXInt32Array1D = newXDataArray1DImp[int32]
	NewXInt64Array1D = newXDataArray1DImp[int64]

	NewXFloat32Array1D = newXDataArray1DImp[float32]
	NewXFloat64Array1D = newXDataArray1DImp[float64]
)
