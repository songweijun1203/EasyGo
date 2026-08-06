package xutils

import "math"

//四维向量
type Vec4[T XNumber] struct {
	X, Y, Z, W T
}

//模长
func (v Vec4[T]) Norm() float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W))
}

//向量缩放
func (v Vec4[T]) Scale(scalar T) Vec4[T] {
	return Vec4[T]{
		X: v.X * scalar,
		Y: v.Y * scalar,
		Z: v.Z * scalar,
		W: v.W * scalar,
	}
}

//向量归一化
func (v Vec4[T]) Normalize() Vec4[T] {
	n := v.Norm()
	if n != 0 {
		return Vec4[T]{
			X: v.X / T(n),
			Y: v.Y / T(n),
			Z: v.Z / T(n),
			W: v.W / T(n),
		}
	}
	return v
}

// 向量加法
func (v Vec4[T]) Add(other Vec4[T]) Vec4[T] {
	return Vec4[T]{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
		W: v.W + other.W,
	}
}

// 向量减法
func (v Vec4[T]) Sub(other Vec4[T]) Vec4[T] {
	return Vec4[T]{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z, W: v.W - other.W}
}

// 向量点积(内积)
func (v Vec4[T]) Dot(other Vec4[T]) T {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z + v.W*other.W
}

// 向量外积，列向量X行向量
func (v Vec4[T]) Outer(row Vec4[T]) Mat4[T] {
	return Mat4[T]{
		data: [16]T{
			v.X * row.X, v.Y * row.X, v.Z * row.X, v.W * row.X,
			v.X * row.Y, v.Y * row.Y, v.Z * row.Y, v.W * row.Y,
			v.X * row.Z, v.Y * row.Z, v.Z * row.Z, v.W * row.Z,
			v.X * row.W, v.Y * row.W, v.Z * row.W, v.W * row.W,
		},
	}
}

//其次坐标归一化
func (v Vec4[T]) Homogenize() Vec4[T] {
	return Vec4[T]{X: v.X / v.W, Y: v.Y / v.W, Z: v.Z / v.W, W: 1}
}

type Vec4f = Vec4[float32]
type Vec4d = Vec4[float64]
type Vec4i = Vec4[int32]
type Vec4u = Vec4[uint32]

type Rectf = Vec4f
type Rectd = Vec4d
type Recti = Vec4i
type Rectu = Vec4u

type Point4f = Vec4f
type Point4d = Vec4d
type Point4i = Vec4i
type Point4u = Vec4u
