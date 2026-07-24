package utils

//二维向量
type Vev2[T XNumber] struct {
	X, Y T
}

type Vev2f = Vev2[float32]
type Vev2d = Vev2[float64]
type Vev2i = Vev2[int]
type Vev2u = Vev2[uint]

//模长
func (v Vev2[T]) Norm() float64 {
	return float64(v.X*v.X + v.Y*v.Y)
}

//向量缩放
func (v *Vev2[T]) Scale(scalar T) {
	v.X *= scalar
	v.Y *= scalar
}

//向量归一化
func (v *Vev2[T]) Normalize() {
	n := v.Norm()
	if n != 0 {
		v.X /= T(n)
		v.Y /= T(n)
	}
}

//向量加法
func (v Vev2[T]) Add(other Vev2[T]) Vev2[T] {
	return Vev2[T]{X: v.X + other.X, Y: v.Y + other.Y}
}

//向量减法
func (v Vev2[T]) Sub(other Vev2[T]) Vev2[T] {
	return Vev2[T]{X: v.X - other.X, Y: v.Y - other.Y}
}

//向量点积
func (v Vev2[T]) Dot(other Vev2[T]) T {
	return v.X*other.X + v.Y*other.Y
}
