package xopengl

import (
	"EasyGo/internal/xutils"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type GLDataType uint32

const (
	GLDataType_uint   GLDataType = gl.UNSIGNED_INT
	GLDataType_uint8  GLDataType = gl.UNSIGNED_BYTE
	GLDataType_uint16 GLDataType = gl.UNSIGNED_SHORT
	GLDataType_int    GLDataType = gl.INT
	GLDataType_float  GLDataType = gl.FLOAT
	GLDataType_double GLDataType = gl.DOUBLE
)

type GLEnableTestType uint32

const (
	DEPTH_TEST          GLEnableTestType = gl.DEPTH_TEST //0x0B71
	MULTISAMPLE         GLEnableTestType = gl.MULTISAMPLE
	SCISSOR_TEST        GLEnableTestType = gl.SCISSOR_TEST
	BLEND               GLEnableTestType = gl.BLEND
	CULL_FACE           GLEnableTestType = gl.CULL_FACE
	PROGRAM_POINT_SIZE  GLEnableTestType = gl.PROGRAM_POINT_SIZE
	STENCIL_TEST        GLEnableTestType = gl.STENCIL_TEST
	DEPTH_CLAMP         GLEnableTestType = gl.DEPTH_CLAMP
	POLYGON_OFFSET_FILL GLEnableTestType = gl.POLYGON_OFFSET_FILL
)

type BlendFuncFactorType uint32

const (
	ZERO                     BlendFuncFactorType = gl.ZERO
	ONE                      BlendFuncFactorType = gl.ONE
	SRC_COLOR                BlendFuncFactorType = gl.SRC_COLOR
	ONE_MINUS_SRC_COLOR      BlendFuncFactorType = gl.ONE_MINUS_SRC_COLOR
	DST_COLOR                BlendFuncFactorType = gl.DST_COLOR
	SRC_ALPHA                BlendFuncFactorType = gl.SRC_ALPHA
	ONE_MINUS_SRC_ALPHA      BlendFuncFactorType = gl.ONE_MINUS_CONSTANT_ALPHA
	DST_ALPHA                BlendFuncFactorType = gl.DST_ALPHA
	ONE_MINUS_DST_ALPHA      BlendFuncFactorType = gl.ONE_MINUS_DST_ALPHA
	CONSTANT_COLOR           BlendFuncFactorType = gl.CONSTANT_COLOR
	ONE_MINUS_CONSTANT_COLOR BlendFuncFactorType = gl.ONE_MINUS_CONSTANT_COLOR
	CONSTANT_ALPHA           BlendFuncFactorType = gl.CONSTANT_ALPHA
	SRC_ALPHA_SATURATE       BlendFuncFactorType = gl.SRC_ALPHA_SATURATE
)

type DepthStencilCompFuntype uint32 // 深度或者模板比较函数

const (
	NEVER    DepthStencilCompFuntype = gl.NEVER    //永不通过
	LESS     DepthStencilCompFuntype = gl.LESS     //小于
	EQUAL    DepthStencilCompFuntype = gl.EQUAL    //等于
	LEQUAL   DepthStencilCompFuntype = gl.LEQUAL   //小于或等于
	GREATER  DepthStencilCompFuntype = gl.GREATER  //大于
	NOTEQUAL DepthStencilCompFuntype = gl.NOTEQUAL //不等于
	GEQUAL   DepthStencilCompFuntype = gl.GEQUAL   //大于或等于
	ALWAYS   DepthStencilCompFuntype = gl.ALWAYS   //总是通过
)

type StencilBehaviorType uint32 //板测试通过或者失败后如何跟新模板缓冲区

const (
	KEEP      StencilBehaviorType = gl.KEEP      //保持当前的默认值不变
	REPLACE   StencilBehaviorType = gl.REPLACE   //用 glStencilFunc 中的 ref 值替换当前模板值
	INCR      StencilBehaviorType = gl.INCR      //加 1（若已达最大值则保持，如 255 → 255）
	INCR_WRAP StencilBehaviorType = gl.INCR_WRAP //加 1，溢出时回绕（255 → 0）
	DECR      StencilBehaviorType = gl.DECR      //减 1（若为 0 则保持）
	DECR_WRAP StencilBehaviorType = gl.DECR_WRAP //减 1，下溢时回绕（0 → 255）
	INVERT    StencilBehaviorType = gl.INVERT    //按位取反（bitwise invert）
	ZERO1     StencilBehaviorType = gl.ZERO
)

// 深度测试状态
type StateDepthTest struct {
	isEnable           bool
	isDepthClampEnable bool
}

// 多重采样状态
type StateMultiSample struct {
	isEnable bool
}

// 剪切测试状态
type StateScissorTest struct {
	isEnable bool
	rect     xutils.Recti
}

// 融混状态
type StateBlendTest struct {
	isEnable  bool
	srcFactor BlendFuncFactorType
	dstFactor BlendFuncFactorType
}

// 裁剪面状态
type StateCullFace struct {
	isEnable bool
}

// 程序点大小状态
type StateProgramPointSize struct {
	isEnable bool
}

// 模板测试状态
type StateStencilTest struct {
	isEnable bool
	compFun  DepthStencilCompFuntype
	ref      int32
	mask     uint32
	sfail    StencilBehaviorType
	zfail    StencilBehaviorType
	zpass    StencilBehaviorType
}

type StatePolyGonOffsetFill struct {
	isEnable bool
	factor   float32
	unit     float32
}
