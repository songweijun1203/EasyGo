package utils

//二维向量
type Vec2[T XNumber] struct {
	X, Y T
}

//模长
func (v Vec2[T]) Norm() float64 {
	return float64(v.X*v.X + v.Y*v.Y)
}

//向量缩放
func (v *Vec2[T]) Scale(scalar T) {
	v.X *= scalar
	v.Y *= scalar
}

//向量归一化
func (v *Vec2[T]) Normalize() {
	n := v.Norm()
	if n != 0 {
		v.X /= T(n)
		v.Y /= T(n)
	}
}

//向量加法
func (v Vec2[T]) Add(other Vec2[T]) Vec2[T] {
	return Vec2[T]{X: v.X + other.X, Y: v.Y + other.Y}
}

//向量减法
func (v Vec2[T]) Sub(other Vec2[T]) Vec2[T] {
	return Vec2[T]{X: v.X - other.X, Y: v.Y - other.Y}
}

//向量点积
func (v Vec2[T]) Dot(other Vec2[T]) T {
	return v.X*other.X + v.Y*other.Y
}

type Vec2f = Vec2[float32]
type Vec2d = Vec2[float64]
type Vec2i = Vec2[int32]
type Vec2u = Vec2[uint32]

type Point2f = Vec2f
type Point2d = Vec2d
type Point2i = Vec2i
type Point2u = Vec2u
