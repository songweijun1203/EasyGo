package xutils

import "math"

//二维向量
type Vec2[T XNumber] struct {
	X, Y T
}

//模长
func (v Vec2[T]) Norm() float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y))
}

//向量缩放
func (v Vec2[T]) Scale(scalar T) Vec2[T] {
	return Vec2[T]{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

//向量归一化
func (v Vec2[T]) Normalize() Vec2[T] {
	n := v.Norm()
	if n != 0 {
		return Vec2[T]{
			X: v.X / T(n),
			Y: v.Y / T(n),
		}
	}
	return v
}

//向量加法
func (v Vec2[T]) Add(other Vec2[T]) Vec2[T] {
	return Vec2[T]{X: v.X + other.X, Y: v.Y + other.Y}
}

//向量减法
func (v Vec2[T]) Sub(other Vec2[T]) Vec2[T] {
	return Vec2[T]{X: v.X - other.X, Y: v.Y - other.Y}
}

//向量内积
func (v Vec2[T]) Dot(other Vec2[T]) T {
	return v.X*other.X + v.Y*other.Y
}

//向量外积
func (v Vec2[T]) Outer(row Vec2[T]) Mat2[T] {
	return Mat2[T]{
		data: [4]T{
			v.X * row.X, v.Y * row.X, //col0
			v.X * row.Y, v.Y * row.Y, //col1
		},
	}
}

type Vec2f = Vec2[float32]
type Vec2d = Vec2[float64]
type Vec2i = Vec2[int32]
type Vec2u = Vec2[uint32]

type Point2f = Vec2f
type Point2d = Vec2d
type Point2i = Vec2i
type Point2u = Vec2u
