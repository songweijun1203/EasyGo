package xutils

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// 三维仿射变换
type XAffine3[T XNumber] struct {
	matrix Mat4[T]
}

func NewAffine3[T XNumber]() XAffine3[T] {
	return XAffine3[T]{
		matrix: NewMat4[T](),
	}
}

func NewAffine3f() XAffine3f {
	return XAffine3f{
		matrix: NewMat4f(),
	}
}

func NewAffine3d() XAffine3d {
	return XAffine3d{
		matrix: NewMat4d(),
	}
}

func (a *XAffine3[T]) SetMatrix(m Mat4[T]) {
	a.matrix = m
}

func (a *XAffine3[T]) SetIdentity() {
	a.matrix.SetIdentity()
}

func (a *XAffine3[T]) Translation() Vec3[T] {
	v4 := a.matrix.ColData(3)
	return Vec3[T]{X: v4.X, Y: v4.Y, Z: v4.Z}
}

func (a *XAffine3[T]) ComputeRotationScaling() (r Mat3[T], s Mat3[T]) {
	a_double := Mat3d{
		data: [9]float64{
			float64(a.matrix.data[0]), float64(a.matrix.data[1]), float64(a.matrix.data[2]),
			float64(a.matrix.data[4]), float64(a.matrix.data[5]), float64(a.matrix.data[6]),
			float64(a.matrix.data[8]), float64(a.matrix.data[9]), float64(a.matrix.data[10]),
		},
	}
	mat := mat3fToMatDense(a_double)
	R, S := computeRotationScaling(mat)
	r_double := matDenseToMat3f(R)
	s_double := matDenseToMat3f(S)
	return Mat3[T]{
			data: [9]T{
				T(r_double.data[0]), T(r_double.data[1]), T(r_double.data[2]),
				T(r_double.data[3]), T(r_double.data[4]), T(r_double.data[5]),
				T(r_double.data[6]), T(r_double.data[7]), T(r_double.data[8]),
			},
		},
		Mat3[T]{
			data: [9]T{
				T(s_double.data[0]), T(s_double.data[1]), T(s_double.data[2]),
				T(s_double.data[3]), T(s_double.data[4]), T(s_double.data[5]),
				T(s_double.data[6]), T(s_double.data[7]), T(s_double.data[8]),
			},
		}
}

func (a *XAffine3[T]) ComputeScalingRotation() (s Mat3[T], r Mat3[T]) {
	a_double := Mat3d{
		data: [9]float64{
			float64(a.matrix.data[0]), float64(a.matrix.data[1]), float64(a.matrix.data[2]),
			float64(a.matrix.data[4]), float64(a.matrix.data[5]), float64(a.matrix.data[6]),
			float64(a.matrix.data[8]), float64(a.matrix.data[9]), float64(a.matrix.data[10]),
		},
	}
	mat := mat3fToMatDense(a_double)
	R, S := computeScalingRotation(mat)
	r_double := matDenseToMat3f(R)
	s_double := matDenseToMat3f(S)
	return Mat3[T]{
			data: [9]T{
				T(s_double.data[0]), T(s_double.data[1]), T(s_double.data[2]),
				T(s_double.data[3]), T(s_double.data[4]), T(s_double.data[5]),
				T(s_double.data[6]), T(s_double.data[7]), T(s_double.data[8]),
			},
		},
		Mat3[T]{
			data: [9]T{
				T(r_double.data[0]), T(r_double.data[1]), T(r_double.data[2]),
				T(r_double.data[3]), T(r_double.data[4]), T(r_double.data[5]),
				T(r_double.data[6]), T(r_double.data[7]), T(r_double.data[8]),
			},
		}
}

func (a *XAffine3[T]) MulRight(b XAffine3[T]) XAffine3[T] {
	m := a.matrix.MulRight(b.matrix)
	a_new := NewAffine3[T]()
	a_new.SetMatrix(m)
	return a_new
}

func (a *XAffine3[T]) MulLeft(b XAffine3[T]) XAffine3[T] {
	m := a.matrix.MulLeft(b.matrix)
	a_new := NewAffine3[T]()
	a_new.SetMatrix(m)
	return a_new
}

// 叠加变换
func (a *XAffine3[T]) Translate(tx, ty, tz T) {
	t := Translate(tx, ty, tz)
	*a = a.MulRight(t)
}

func (a *XAffine3[T]) Scale(sx, sy, sz T) {
	t := Scale(sx, sy, sz)
	*a = a.MulRight(t)
}

func (a *XAffine3[T]) RotateX(radian float32) {
	r := RotateX[T](radian)
	*a = a.MulRight(r)
}

func (a *XAffine3[T]) RotateY(radian float32) {
	r := RotateY[T](radian)
	*a = a.MulRight(r)
}

func (a *XAffine3[T]) RotateZ(radian float32) {
	r := RotateZ[T](radian)
	*a = a.MulRight(r)
}

type XAffine3f = XAffine3[float32]
type XAffine3d = XAffine3[float64]

// ------------------------------------------------------------------
func mat3fToMatDense(m Mat3d) *mat.Dense {
	dense := mat.NewDense(3, 3, nil)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			// m 是列主序：索引 j*3 + i
			dense.Set(i, j, float64(m.data[j*3+i]))
		}
	}
	return dense
}

func matDenseToMat3f(m *mat.Dense) Mat3d {
	if m == nil {
		return Mat3d{}
	}
	r, c := m.Dims()
	if r != 3 || c != 3 {
		panic("matrix must be 3x3")
	}
	var res Mat3d
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			// 列主序：第 j 列第 i 行 -> res[j*3+i]
			res.data[j*3+i] = float64(m.At(i, j))
		}
	}
	return res
}

func computeRotationScaling(M *mat.Dense) (R, S *mat.Dense) {
	r, c := M.Dims()
	if r != 3 || c != 3 {
		panic("matrix must be 3x3")
	}

	// 1. SVD: M = U * Σ * V^T
	var svd mat.SVD
	ok := svd.Factorize(M, mat.SVDThin)
	if !ok {
		panic("SVD factorization failed")
	}

	var U, V mat.Dense
	svd.UTo(&U)
	svd.VTo(&V)
	vals := svd.Values(nil) // 奇异值切片，降序

	// 2. 计算 x = det(U * V^T)，应为 ±1
	var UV mat.Dense
	UV.Mul(&U, V.T())
	x := mat.Det(&UV)
	// 浮点误差修正
	if math.Abs(x+1) < 1e-10 {
		x = -1
	} else if math.Abs(x-1) < 1e-10 {
		x = 1
	}

	// 3. 调整奇异值：第一个乘以 x
	vals[0] *= x

	// 4. 缩放矩阵 S = V * diag(vals) * V^T
	diag := mat.NewDiagDense(len(vals), vals)
	var temp mat.Dense
	temp.Mul(&V, diag)
	S = new(mat.Dense)
	S.Mul(&temp, V.T())

	// 5. 构造 diagX: [1/x, 1, 1]
	diagX := mat.NewDiagDense(len(vals), nil)
	diagX.SetDiag(0, 1/x)
	for i := 1; i < len(vals); i++ {
		diagX.SetDiag(i, 1)
	}
	// 旋转矩阵 R = U * diagX * V^T
	var rotTemp mat.Dense
	rotTemp.Mul(&U, diagX)
	R = new(mat.Dense)
	R.Mul(&rotTemp, V.T())

	return R, S
}

func computeScalingRotation(M *mat.Dense) (R, S *mat.Dense) {
	r, c := M.Dims()
	if r != 3 || c != 3 {
		panic("matrix must be 3x3")
	}

	var svd mat.SVD
	ok := svd.Factorize(M, mat.SVDThin)
	if !ok {
		panic("SVD factorization failed")
	}

	var U, V mat.Dense
	svd.UTo(&U)
	svd.VTo(&V)
	vals := svd.Values(nil)

	// x = det(U * V^T)
	var UV mat.Dense
	UV.Mul(&U, V.T())
	x := mat.Det(&UV)
	if math.Abs(x+1) < 1e-10 {
		x = -1
	} else if math.Abs(x-1) < 1e-10 {
		x = 1
	}

	vals[0] *= x

	// 缩放矩阵 S = U * diag(vals) * U^T
	diag := mat.NewDiagDense(len(vals), vals)
	var temp mat.Dense
	temp.Mul(&U, diag)
	S = new(mat.Dense)
	S.Mul(&temp, U.T())

	// 旋转矩阵 R = U * diagX * V^T（与 ComputeRotationScaling 相同）
	diagX := mat.NewDiagDense(len(vals), nil)
	diagX.SetDiag(0, 1/x)
	for i := 1; i < len(vals); i++ {
		diagX.SetDiag(i, 1)
	}
	var rotTemp mat.Dense
	rotTemp.Mul(&U, diagX)
	R = new(mat.Dense)
	R.Mul(&rotTemp, V.T())

	return R, S
}

// 基本矩阵
func Translate[T XNumber](tx, ty, tz T) XAffine3[T] {
	a := NewAffine3[T]()
	a.matrix.SetData(3, 0, tx)
	a.matrix.SetData(3, 1, ty)
	a.matrix.SetData(3, 2, tz)
	return a
}

func Scale[T XNumber](sx, sy, sz T) XAffine3[T] {
	m := NewMat4[T]()
	m.SetData(0, 0, sx)
	m.SetData(1, 1, sy)
	m.SetData(2, 2, sz)

	a := NewAffine3[T]()
	a.matrix = m
	return a
}

func RotateZ[T XNumber](radian float32) XAffine3[T] {
	a := NewAffine3[T]()
	a.matrix.data = [16]T{
		T(math.Cos(float64(radian))), T(math.Sin(float64(radian))), 0, 0,
		-T(math.Sin(float64(radian))), T(math.Cos(float64(radian))), 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	return a
}

func RotateX[T XNumber](radian float32) XAffine3[T] {
	a := NewAffine3[T]()
	a.matrix.data = [16]T{
		1, 0, 0, 0,
		0, T(math.Cos(float64(radian))), T(math.Sin(float64(radian))), 0,
		0, -T(math.Sin(float64(radian))), T(math.Cos(float64(radian))), 0,
		0, 0, 0, 1,
	}
	return a
}

func RotateY[T XNumber](radian float32) XAffine3[T] {
	a := NewAffine3[T]()
	a.matrix.data = [16]T{
		T(math.Cos(float64(radian))), 0, -T(math.Sin(float64(radian))), 0,
		0, 1, 0, 0,
		T(math.Sin(float64(radian))), 0, T(math.Cos(float64(radian))), 0,
		0, 0, 0, 1,
	}
	return a
}
