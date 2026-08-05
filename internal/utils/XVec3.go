package utils

//三维向量
type Vec3[T XNumber] struct {
	X, Y, Z T
}

func (v Vec3[T]) Add(other Vec3[T]) Vec3[T] {
	return Vec3[T]{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

type Vec3f = Vec3[float32]
type Vec3d = Vec3[float64]
type Vec3i = Vec3[int]
type Vec3u = Vec3[uint]

type Point3f = Vec3f
type Point3d = Vec3d
type Point3i = Vec3i
type Point3u = Vec3u
