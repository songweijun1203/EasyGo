package utils

//四维向量
type Vec4[T XNumber] struct {
	X, Y, Z, W T
}

type Vec4f = Vec4[float32]
type Vec4d = Vec4[float64]
type Vec4i = Vec4[int]
type Vec4u = Vec4[uint]
