package xutils

//3*3列矩阵 列主序
type Mat3[T XNumber] struct {
	data [9]T
}

//默认创建单位矩阵
func NewMat3[T XNumber]() Mat3[T] {
	return Mat3[T]{
		data: [9]T{1, 0, 0, 0, 1, 0, 0, 0, 1},
	}
}

func NewMat3f() Mat3f {
	return NewMat3[float32]()
}

func NewMat3d() Mat3d {
	return NewMat3[float64]()
}

func (m *Mat3[T]) SetIdentity() {
	m.data = [9]T{1, 0, 0, 0, 1, 0, 0, 0, 1}
}

//获取矩阵
func (m *Mat3[T]) Pointer() *T {
	return &m.data[0]
}

func (m *Mat3[T]) At(col uint32, row uint32) T {
	if col > 2 {
		return 0
	}
	if row > 2 {
		return 0
	}
	index := col*3 + row
	return m.data[index]
}

func (m *Mat3[T]) SetData(col uint32, row uint32, v T) T {
	if col > 2 {
		return 0
	}
	if row > 2 {
		return 0
	}
	index := col*3 + row

	old := m.data[index]
	m.data[index] = v
	return old
}

func (m *Mat3[T]) Add(t Mat3[T]) Mat3[T] {
	return Mat3[T]{
		data: [9]T{
			m.data[0] + t.data[0], m.data[1] + t.data[1], m.data[2] + t.data[2],
			m.data[3] + t.data[3], m.data[4] + t.data[4], m.data[5] + t.data[5],
			m.data[6] + t.data[6], m.data[7] + t.data[7], m.data[8] + t.data[8],
		},
	}
}

//获取列数据
func (m *Mat3[T]) ColData(colIdx uint32) Vec3[T] {
	if colIdx == 0 {
		return Vec3[T]{X: m.data[0], Y: m.data[1], Z: m.data[2]}
	}

	if colIdx == 1 {
		return Vec3[T]{X: m.data[3], Y: m.data[4], Z: m.data[5]}
	}

	if colIdx == 2 {
		return Vec3[T]{X: m.data[6], Y: m.data[7], Z: m.data[8]}
	}

	return Vec3[T]{0, 0, 0}
}

//矩阵向量乘法
func (m *Mat3[T]) MulVec(v Vec3[T]) Vec3[T] {
	col0 := m.ColData(0)                                             //第0列数据
	col1 := m.ColData(1)                                             //第1列数据
	col2 := m.ColData(2)                                             //第2列数据
	return col0.Scale(v.X).Add(col1.Scale(v.Y)).Add(col2.Scale(v.Y)) //列向量的线性组合
}

//矩阵乘法,右乘 m*n
func (m *Mat3[T]) MulRight(n Mat3[T]) Mat3[T] {
	c0 := m.MulVec(n.ColData(0))
	c1 := m.MulVec(n.ColData(1))
	c2 := m.MulVec(n.ColData(2))

	return Mat3[T]{
		data: [9]T{
			c0.X, c0.Y, c0.Z,
			c1.X, c1.Y, c1.Z,
			c2.X, c2.Y, c2.Z,
		},
	}
}

//矩阵乘法,左乘 n*m
func (m *Mat3[T]) MulLeft(n Mat3[T]) Mat3[T] {
	return n.MulRight(*m)
}

func (m *Mat3[T]) SubMat2() Mat2[T] {
	return Mat2[T]{
		data: [4]T{
			m.data[0], m.data[1],
			m.data[3], m.data[4],
		},
	}
}

type Mat3f = Mat3[float32]
type Mat3i = Mat3[int32]
type Mat3u = Mat3[uint32]
type Mat3d = Mat3[float64]
