package utils

//三维向量
type Vec3[T XNumber] struct {
	X, Y, Z T
}

type Vec3f = Vec3[float32]
type Vec3d = Vec3[float64]
type Vec3i = Vec3[int]
type Vec3u = Vec3[uint]

func (v Vec3[T]) Add(other Vec3[T]) Vec3[T] {
	return Vec3[T]{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}
