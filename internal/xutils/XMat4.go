package xutils

type Mat4[T XNumber] struct {
	data [16]T
}

func NewMat4[T XNumber]() Mat4[T] {
	return Mat4[T]{
		data: [16]T{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		},
	}
}

func NewMat4f() Mat4f {
	return NewMat4[float32]()
}

func NewMat4d() Mat4d {
	return NewMat4[float64]()
}

func (m *Mat4[T]) SetIdentity() {
	m.data = [16]T{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

// 获取矩阵
func (m *Mat4[T]) Pointer() *T {
	return &m.data[0]
}

func (m *Mat4[T]) At(col uint32, row uint32) T {
	if col > 3 {
		return 0
	}
	if row > 3 {
		return 0
	}
	index := col*4 + row
	return m.data[index]
}

func (m *Mat4[T]) SetData(col uint32, row uint32, v T) T {
	if col > 3 {
		return 0
	}
	if row > 3 {
		return 0
	}
	index := col*4 + row

	old := m.data[index]
	m.data[index] = v
	return old
}

func (m *Mat4[T]) Add(t Mat4[T]) Mat4[T] {
	return Mat4[T]{
		data: [16]T{
			m.data[0] + t.data[0], m.data[1] + t.data[1], m.data[2] + t.data[2], m.data[3] + t.data[3],
			m.data[4] + t.data[4], m.data[5] + t.data[5], m.data[6] + t.data[6], m.data[7] + t.data[7],
			m.data[8] + t.data[8], m.data[9] + t.data[9], m.data[10] + t.data[10], m.data[11] + t.data[11],
			m.data[12] + t.data[12], m.data[13] + t.data[13], m.data[14] + t.data[14], m.data[15] + t.data[15],
		},
	}
}

// 获取列数据
func (m *Mat4[T]) ColData(colIdx uint32) Vec4[T] {
	if colIdx == 0 {
		return Vec4[T]{X: m.data[0], Y: m.data[1], Z: m.data[2], W: m.data[3]}
	}

	if colIdx == 1 {
		return Vec4[T]{X: m.data[4], Y: m.data[5], Z: m.data[6], W: m.data[7]}
	}

	if colIdx == 2 {
		return Vec4[T]{X: m.data[8], Y: m.data[9], Z: m.data[10], W: m.data[11]}
	}

	if colIdx == 3 {
		return Vec4[T]{X: m.data[12], Y: m.data[13], Z: m.data[14], W: m.data[15]}
	}

	return Vec4[T]{0, 0, 0, 0}
}

// 矩阵向量乘法
func (m *Mat4[T]) MulVec(v Vec4[T]) Vec4[T] {
	col0 := m.ColData(0) //第0列数据
	col1 := m.ColData(1) //第1列数据
	col2 := m.ColData(2) //第2列数据
	col3 := m.ColData(3) //第3列数据
	return col0.Scale(v.X).Add(col1.Scale(v.Y)).Add(col2.Scale(v.Y)).Add(col3.Scale(v.W))
}

// 矩阵乘法,右乘 m*n
func (m *Mat4[T]) MulRight(n Mat4[T]) Mat4[T] {
	c0 := m.MulVec(n.ColData(0))
	c1 := m.MulVec(n.ColData(1))
	c2 := m.MulVec(n.ColData(2))
	c3 := m.MulVec(n.ColData(3))

	return Mat4[T]{
		data: [16]T{
			c0.X, c0.Y, c0.Z, c0.W,
			c1.X, c1.Y, c1.Z, c1.W,
			c2.X, c2.Y, c2.Z, c2.W,
			c3.X, c3.Y, c3.Z, c3.W,
		},
	}
}

// 矩阵乘法,左乘 n*m
func (m *Mat4[T]) MulLeft(n Mat4[T]) Mat4[T] {
	return n.MulRight(*m)
}

func (m *Mat4[T]) SubMat3() Mat3[T] {
	return Mat3[T]{
		data: [9]T{
			m.data[0], m.data[1], m.data[2],
			m.data[4], m.data[5], m.data[6],
			m.data[8], m.data[9], m.data[10],
		},
	}
}

type Mat4f = Mat4[float32]
type Mat4i = Mat4[int32]
type Mat4u = Mat4[uint32]
type Mat4d = Mat4[float64]
