package utils

//四维向量
type Vev4[T XNumber] struct {
	X, Y, Z, W T
}

type Vev4f = Vev4[float32]
type Vev4d = Vev4[float64]
type Vev4i = Vev4[int]
type Vev4u = Vev4[uint]
