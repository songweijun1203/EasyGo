package xutils

import "math"

//三维向量
type Vec3[T XNumber] struct {
	X, Y, Z T
}

//模长
func (v Vec3[T]) Norm() float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z))
}

//向量缩放
func (v Vec3[T]) Scale(scalar T) Vec3[T] {
	return Vec3[T]{
		X: v.X * scalar,
		Y: v.Y * scalar,
		Z: v.Z * scalar,
	}
}

//向量归一化
func (v Vec3[T]) Normalize() Vec3[T] {
	n := v.Norm()
	if n != 0 {
		return Vec3[T]{
			X: v.X / T(n),
			Y: v.Y / T(n),
			Z: v.Z / T(n),
		}
	}
	return v
}

//向量加法
func (v Vec3[T]) Add(other Vec3[T]) Vec3[T] {
	return Vec3[T]{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

//向量减法
func (v Vec3[T]) Sub(other Vec3[T]) Vec3[T] {
	return Vec3[T]{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z}
}

//向量点积(内积)
func (v Vec3[T]) Dot(other Vec3[T]) T {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

//向量外积，列向量X行向量
func (v Vec3[T]) Outer(row Vec3[T]) Mat3[T] {
	return Mat3[T]{
		data: [9]T{
			v.X * row.X, v.Y * row.X, v.Z * row.X,
			v.X * row.Y, v.Y * row.Y, v.Z * row.Y,
			v.X * row.Z, v.Y * row.Z, v.Z * row.Z,
		},
	}
}

//其次坐标归一化
func (v Vec3[T]) Homogenize() Vec3[T] {
	return Vec3[T]{X: v.X / v.Z, Y: v.Y / v.Z, Z: 1}
}

//向量叉积 v×t
func (v Vec3[T]) Cross(t Vec3[T]) Vec3[T] {
	return Vec3[T]{
		X: v.Y*t.Z - v.Z*t.Y,
		Y: v.Z*t.X - v.X*t.Z,
		Z: v.X*t.Y - v.Y*t.X,
	}
}

type Vec3f = Vec3[float32]
type Vec3d = Vec3[float64]
type Vec3i = Vec3[int]
type Vec3u = Vec3[uint]

type Point3f = Vec3f
type Point3d = Vec3d
type Point3i = Vec3i
type Point3u = Vec3u
