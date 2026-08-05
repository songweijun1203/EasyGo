package utils

//四维向量
type Vec4[T XNumber] struct {
	X, Y, Z, W T
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
