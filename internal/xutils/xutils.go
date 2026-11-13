package xutils

import (
	"math"
)

func Radians(deg float32) float32 {
	return math.Pi * deg / 180
}

func Degrees(rad float32) float32 {
	return rad * 180 / math.Pi
}

//旋转矩阵转欧拉角
