package xutils

import (
	"fmt"
	"math"
)

// 3*3列矩阵 列主序
type Mat3[T XNumber] struct {
	data [9]T
}

// 默认创建单位矩阵
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

// 获取矩阵
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

// 获取列数据
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

// 矩阵向量乘法
func (m *Mat3[T]) MulVec(v Vec3[T]) Vec3[T] {
	col0 := m.ColData(0)                                             //第0列数据
	col1 := m.ColData(1)                                             //第1列数据
	col2 := m.ColData(2)                                             //第2列数据
	return col0.Scale(v.X).Add(col1.Scale(v.Y)).Add(col2.Scale(v.Y)) //列向量的线性组合
}

// 矩阵乘法,右乘 m*n
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

// 矩阵乘法,左乘 n*m
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

func (m *Mat3[T]) FormatPrint() {
	fmt.Println("matrix: ")
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			fmt.Printf("%v ", m.At(uint32(c), uint32(r)))
		}
		fmt.Println()
	}

}

// EulerAngles 提取欧拉角，参数 a0, a1, a2 表示旋转顺序（0:X, 1:Y, 2:Z）
// 返回 Vec3[T] 分别对应 a0, a1, a2 轴的角度（弧度）。
// 注意：T 应为浮点数类型（float32 或 float64）。
func (m *Mat3[T]) EulerAngles(a0, a1, a2 int) Vec3[T] {
	if a0 < 0 || a0 > 2 || a1 < 0 || a1 > 2 || a2 < 0 || a2 > 2 {
		panic("EulerAngles: axis indices must be 0, 1, or 2")
	}

	odd := 0
	if (a0+1)%3 != a1 {
		odd = 1
	}
	i := a0
	j := (a0 + 1 + odd) % 3
	k := (a0 + 2 - odd) % 3

	// 注意：get(列, 行)
	get := func(col, row int) T { return m.At(uint32(col), uint32(row)) }
	toF64 := func(v T) float64 { return float64(v) }

	res := Vec3[T]{}

	if a0 == a2 {
		// 此分支较少用，但为了完整，也修正索引
		val0 := toF64(get(i, j)) // coeff(j,i) → get(i,j)
		val1 := toF64(get(i, k)) // coeff(k,i) → get(i,k)
		angle0 := math.Atan2(val0, val1)

		if (odd != 0 && angle0 < 0) || (odd == 0 && angle0 > 0) {
			if angle0 > 0 {
				angle0 -= math.Pi
			} else {
				angle0 += math.Pi
			}
			s2 := math.Hypot(val0, val1)
			val2 := toF64(get(i, i)) // coeff(i,i) → get(i,i)
			if odd != 0 {
				res.Y = T(-math.Atan2(s2, val2))
			} else {
				res.Y = T(math.Atan2(s2, val2))
			}
		} else {
			s2 := math.Hypot(val0, val1)
			res.Y = T(math.Atan2(s2, toF64(get(i, i))))
		}

		s1 := math.Sin(float64(res.X))
		c1 := math.Cos(float64(res.X))
		// coeff(j,k) → get(k,j), coeff(k,k) → get(k,k)
		// coeff(j,j) → get(j,j), coeff(k,j) → get(j,k)
		res.Z = T(math.Atan2(
			c1*toF64(get(k, j))-s1*toF64(get(k, k)),
			c1*toF64(get(j, j))-s1*toF64(get(j, k)),
		))
		res.X = T(angle0)
	} else {
		// ===== 普通分支（a0 != a2） =====
		// res[0] = atan2(coeff(j,k), coeff(k,k))
		res.X = T(math.Atan2(toF64(get(k, j)), toF64(get(k, k))))

		// c2 = sqrt(coeff(i,i)^2 + coeff(i,j)^2)
		c2 := math.Hypot(toF64(get(i, i)), toF64(get(j, i)))

		if (odd != 0 && float64(res.X) < 0) || (odd == 0 && float64(res.X) > 0) {
			if float64(res.X) > 0 {
				res.X = T(float64(res.X) - math.Pi)
			} else {
				res.X = T(float64(res.X) + math.Pi)
			}
			// res[1] = atan2(-coeff(i,k), -c2)
			res.Y = T(math.Atan2(-toF64(get(k, i)), -c2))
		} else {
			// res[1] = atan2(-coeff(i,k), c2)
			res.Y = T(math.Atan2(-toF64(get(k, i)), c2))
		}

		s1 := math.Sin(float64(res.X))
		c1 := math.Cos(float64(res.X))
		// res[2] = atan2( s1*coeff(k,i) - c1*coeff(j,i),  c1*coeff(j,j) - s1*coeff(k,j) )
		res.Z = T(math.Atan2(
			s1*toF64(get(i, k))-c1*toF64(get(i, j)),
			c1*toF64(get(j, j))-s1*toF64(get(j, k)),
		))
	}

	if odd == 0 {
		res.X = -res.X
		res.Y = -res.Y
		res.Z = -res.Z
	}

	return res
}

type Mat3f = Mat3[float32]
type Mat3i = Mat3[int32]
type Mat3u = Mat3[uint32]
type Mat3d = Mat3[float64]
