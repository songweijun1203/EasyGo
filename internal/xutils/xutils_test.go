package xutils

import (
	"fmt"
	"testing"
)

// 测试单位矩阵
func TestMat4Identity(t *testing.T) {
	m := NewMat4f()
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			val := m.At(uint32(col), uint32(row))
			if col == row {
				if val != 1.0 {
					t.Errorf("对角元素 (%d,%d) 应为 1，实际为 %v", row, col, val)
				}
			} else {
				if val != 0.0 {
					t.Errorf("非对角元素 (%d,%d) 应为 0，实际为 %v", row, col, val)
				}
			}
		}
	}
	m.FormatPrint()
}

// 测试矩阵乘法（右乘）
func TestMat4MulRight(t *testing.T) {
	affine := NewAffine3f()
	affine.matrix.FormatPrint()

	affine.Translate(10, 20, 30)
	affine.matrix.FormatPrint()

	affine.RotateZ(Radians(30))
	affine.RotateX(Radians(40))
	affine.RotateY(Radians(50))
	affine.matrix.FormatPrint()

	affine.Scale(2, 4, 6)
	affine.matrix.FormatPrint()

}

// 测试极分解（需要 gonum）
func TestComputeRotationScaling(t *testing.T) {
	affine := NewAffine3f()

	affine.Translate(10, 20, 30)

	affine.RotateZ(Radians(30))
	affine.RotateX(Radians(40))
	affine.RotateY(Radians(50))

	affine.Scale(2, 4, 6)
	affine.matrix.FormatPrint()

	r, s := affine.ComputeRotationScaling()
	fmt.Println("ComputeRotationScaling")
	r.FormatPrint()
	s.FormatPrint()
}

// 测试 ComputeScalingRotation
func TestComputeScalingRotation(t *testing.T) {
	affine := NewAffine3f()

	affine.Translate(10, 20, 30)

	affine.RotateZ(Radians(30))
	affine.RotateX(Radians(40))
	affine.RotateY(Radians(50))

	affine.Scale(2, 4, 6)
	affine.matrix.FormatPrint()

	s, r := affine.ComputeScalingRotation()
	r.FormatPrint()
	s.FormatPrint()
}

// 测试 TestEulerAngles欧拉角提取测试
func TestEulerAngles(t *testing.T) {
	affine := NewAffine3f()
	affine.RotateZ(Radians(100))
	affine.RotateX(Radians(90))
	affine.RotateY(Radians(50))
	rotmat := affine.matrix.SubMat3()
	rotmat.FormatPrint()
	result := rotmat.EulerAngles(2, 0, 1)
	fmt.Printf("Z: %v,X:%v,Y:%v ", Degrees(result.X), Degrees(result.Y), Degrees(result.Z))

}
