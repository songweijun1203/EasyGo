package xutils

//2*2矩阵 列主序
type Mat2[T XNumber] struct {
	data [4]T
}

func NewMat2[T XNumber]() Mat2[T] {
	return Mat2[T]{
		data: [4]T{1, 0, 0, 1},
	}
}

func NewMat2f() Mat2f {
	return NewMat2[float32]()
}

func NewMat2d() Mat2d {
	return NewMat2[float64]()
}

func (m *Mat2[T]) SetIdentity() {
	m.data = [4]T{1, 0, 0, 1}
}

//获取矩阵
func (m *Mat2[T]) Pointer() *T {
	return &m.data[0]
}

func (m *Mat2[T]) At(col uint32, row uint32) T {
	if col > 1 {
		return 0
	}
	if row > 1 {
		return 0
	}
	index := col*2 + row
	return m.data[index]
}

func (m *Mat2[T]) SetData(col uint32, row uint32, v T) T {
	if col > 1 {
		return 0
	}
	if row > 1 {
		return 0
	}
	index := col*2 + row

	old := m.data[index]
	m.data[index] = v
	return old
}

//矩阵相加
func (m *Mat2[T]) Add(t Mat2[T]) Mat2[T] {
	return Mat2[T]{
		data: [4]T{
			m.data[0] + t.data[0], m.data[1] + t.data[1],
			m.data[2] + t.data[2], m.data[3] + t.data[3],
		},
	}
}

//获取列数据
func (m *Mat2[T]) ColData(colIdx uint32) Vec2[T] {
	if colIdx == 0 {
		return Vec2[T]{X: m.data[0], Y: m.data[1]}
	}

	if colIdx == 1 {
		return Vec2[T]{X: m.data[2], Y: m.data[3]}
	}

	return Vec2[T]{0, 0}
}

//矩阵向量乘法
func (m *Mat2[T]) MulVec(v Vec2[T]) Vec2[T] {
	col0 := m.ColData(0)                        //第0列数据
	col1 := m.ColData(1)                        //第1列数据
	return col0.Scale(v.X).Add(col1.Scale(v.Y)) //列向量的线性组合
}

//矩阵乘法,右乘 m*n
func (m *Mat2[T]) MulRight(n Mat2[T]) Mat2[T] {
	c0 := m.MulVec(n.ColData(0))
	c1 := m.MulVec(n.ColData(1))

	return Mat2[T]{
		data: [4]T{
			c0.X, c0.Y,
			c1.X, c1.Y,
		},
	}
}

//矩阵乘法,左乘 n*m
func (m *Mat2[T]) MulLeft(n Mat2[T]) Mat2[T] {
	return n.MulRight(*m)
}

type Mat2f = Mat2[float32]
type Mat2i = Mat2[int32]
type Mat2u = Mat2[uint32]
type Mat2d = Mat2[float64]
