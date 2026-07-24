package utils

type XDataArray2D[T XNumber] struct {
	*XDataArray1D[T]
	row int
	col int
}

func newXDataArray2DImp[T XNumber](width int, height int, component int) *XDataArray2D[T] {
	subdata := newXDataArray1DImp[T]()
	a := &XDataArray2D[T]{
		col: width,
		row: height,
	}
	a.XDataArray1D = subdata
	subdata.SetComponent(component)
	subdata.SetNumOfTuple(width * height)
	return a
}

// GetData 获取底层的一维数组指针
func (a *XDataArray2D[T]) GetData() *XDataArray1D[T] {
	return a.XDataArray1D
}

func (a *XDataArray2D[T]) SetDimensions(width int, height int) {
	a.SetNumOfTuple(width * height)
	a.row = height
	a.col = width
}

func (a *XDataArray2D[T]) GetDimensions() (row int, col int) {
	return a.row, a.col
}

type (
	XUInt8Array2D  = XDataArray2D[uint8]
	XUInt16Array2D = XDataArray2D[uint16]
	XUInt32Array2D = XDataArray2D[uint32]
	XUInt64Array2D = XDataArray2D[uint64]

	XInt8Array2D  = XDataArray2D[int8]
	XInt16Array2D = XDataArray2D[int16]
	XInt32Array2D = XDataArray2D[int32]
	XInt64Array2D = XDataArray2D[int64]

	XFloat32Array2D = XDataArray2D[float32]
	XFloat64Array2D = XDataArray2D[float64]
)

// 将具体的构造函数赋值给别名，调用时就像调用一个类名一样！
var (
	NewXUInt8Array2D  = newXDataArray2DImp[uint8]
	NewXUInt16Array2D = newXDataArray2DImp[uint16]
	NewXUInt32Array2D = newXDataArray2DImp[uint32]
	NewXUInt64Array2D = newXDataArray2DImp[uint64]

	NewXInt8Array2D  = newXDataArray2DImp[int8]
	NewXInt16Array2D = newXDataArray2DImp[int16]
	NewXInt32Array2D = newXDataArray2DImp[int32]
	NewXInt64Array2D = newXDataArray2DImp[int64]

	NewXFloat32Array2D = newXDataArray2DImp[float32]
	NewXFloat64Array2D = newXDataArray2DImp[float64]
)
