package utils

//三维向量
type Vev3[T XNumber] struct {
	X, Y, Z T
}

type Vev3f = Vev3[float32]
type Vev3d = Vev3[float64]
type Vev3i = Vev3[int]
type Vev3u = Vev3[uint]

func (v Vev3[T]) Add(other Vev3[T]) Vev3[T] {
	return Vev3[T]{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}
