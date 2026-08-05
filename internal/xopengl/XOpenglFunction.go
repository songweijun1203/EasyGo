package xopengl

import (
	xutils "EasyGo/internal/utils"
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func getBindingEnum(t BufferType) uint32 {
	switch t {
	case BufferTypeVertexBuffer:
		return gl.ARRAY_BUFFER_BINDING
	case BufferTypeIndexBuffer:
		return gl.ELEMENT_ARRAY_BUFFER_BINDING
	case BufferTypePixelPackBuffer:
		return gl.PIXEL_PACK_BUFFER_BINDING
	case BufferTypeUniformBuffer:
		return gl.UNIFORM_BUFFER_BINDING
	case BufferTypeShaderStorageBuffer:
		return gl.SHADER_STORAGE_BUFFER_BINDING
	case BufferTypeTransformFeedback:
		// 注意：这是 GL_TRANSFORM_FEEDBACK_BUFFER_BINDING（用于 buffer 目标）
		return gl.TRANSFORM_FEEDBACK_BUFFER_BINDING
	default:
		return 0
	}
}

// 判断函数调用是否发生了错误
func CheckGLerror() error {
	err := gl.GetError()
	if err == gl.NO_ERROR {
		return nil
	}
	// 将错误码转换为可读字符串
	errMsg := "UNKNOWN_ERROR"
	switch err {
	case gl.INVALID_ENUM:
		errMsg = "INVALID_ENUM"
	case gl.INVALID_VALUE:
		errMsg = "INVALID_VALUE"
	case gl.INVALID_OPERATION:
		errMsg = "INVALID_OPERATION"
	case gl.OUT_OF_MEMORY:
		errMsg = "OUT_OF_MEMORY"
	case gl.INVALID_FRAMEBUFFER_OPERATION:
		errMsg = "INVALID_FRAMEBUFFER_OPERATION"
	}
	//fmt.Printf("[OpenGL Error] Code: %d (%s)\n", err, errMsg)
	return fmt.Errorf("[OpenGL Error] Code: %d (%s)\n", err, errMsg)
}

// 获取当前指定的buffer类型的激活id
func GetBindDataBufferId(bufferType BufferType) (id int32, err error) {
	bindValue := getBindingEnum(bufferType)
	var bindId int32 = 0
	gl.GetIntegerv(uint32(bindValue), &bindId)
	return bindId, CheckGLerror()
}

// 设置视口，返回旧的视口
func XglViewPort(x, y, width, height int32) (xutils.Recti, error) {
	viewport := [4]int32{}
	gl.GetIntegerv(gl.VIEWPORT, &viewport[0])
	gl.Viewport(int32(x), int32(y), int32(width), int32(height))
	return xutils.Recti{X: viewport[0], Y: viewport[1], Z: viewport[2], W: viewport[3]}, CheckGLerror()
}

// 设置剪切区域，返回旧的剪切区域
func XglScissor(x, y, width, height int32) (xutils.Recti, error) {
	scissor := [4]int32{}
	gl.GetIntegerv(gl.SCISSOR_BOX, &scissor[0])
	gl.Scissor(x, y, width, height)
	return xutils.Recti{X: scissor[0], Y: scissor[1], Z: scissor[2], W: scissor[3]}, CheckGLerror()
}

// 获取剪切区域
func XglGetScissor() (xutils.Recti, error) {
	scissor := [4]int32{}
	gl.GetIntegerv(gl.SCISSOR_BOX, &scissor[0])
	return xutils.Recti{X: scissor[0], Y: scissor[1], Z: scissor[2], W: scissor[3]}, CheckGLerror()
}

// 获取混合函数
func XgetBlendFunc() (srcFactor BlendFuncFactorType, dstFactor BlendFuncFactorType, err error) {
	var s, d int32
	gl.GetIntegerv(uint32(gl.BLEND_SRC_RGB), &s)
	gl.GetIntegerv(uint32(gl.BLEND_DST_RGB), &d)
	return BlendFuncFactorType(s), BlendFuncFactorType(d), CheckGLerror()
}

// 设置混合函数
func XglBlendFunc(srcFactor BlendFuncFactorType, dstFactor BlendFuncFactorType) error {
	gl.BlendFunc(uint32(srcFactor), uint32(dstFactor))
	return CheckGLerror()
}

// 获取模板设置相关参数
func XglGetStencilFunc() (fun DepthStencilCompFuntype, ref int32, mask uint32, err error) {
	var f, r, m int32
	gl.GetIntegerv(uint32(gl.STENCIL_FUNC), &f)
	gl.GetIntegerv(uint32(gl.STENCIL_REF), &r)
	gl.GetIntegerv(uint32(gl.STENCIL_VALUE_MASK), &m)
	return DepthStencilCompFuntype(f), r, uint32(m), CheckGLerror()
}

// 设置模板参数
func XglStencilFunc(fun DepthStencilCompFuntype, ref int32, mask uint32) error {
	gl.StencilFunc(uint32(fun), ref, mask)
	return CheckGLerror()
}

// 获取模板通过行为
func XGetStencilOp() (sfail StencilBehaviorType, zfail StencilBehaviorType, zpass StencilBehaviorType, err error) {
	var sfail_, zfail_, zpass_ int32
	gl.GetIntegerv(gl.STENCIL_FAIL, &sfail_)
	gl.GetIntegerv(gl.STENCIL_PASS_DEPTH_FAIL, &zfail_)
	gl.GetIntegerv(gl.STENCIL_PASS_DEPTH_PASS, &zpass_)
	return StencilBehaviorType(sfail_), StencilBehaviorType(zfail_), StencilBehaviorType(zpass_), CheckGLerror()
}

// 设置模板通过行为
func XglSetStencilOp(sfail StencilBehaviorType, zfail StencilBehaviorType, zpass StencilBehaviorType) error {
	gl.StencilOp(uint32(sfail), uint32(zfail), uint32(zpass))
	return CheckGLerror()
}

// 获取polygonOffet参数
func XGetPolygonOffsetParm() (factor float32, unit float32, err error) {
	var factor_, unit_ float32
	gl.GetFloatv(uint32(gl.POLYGON_OFFSET_FACTOR), &factor_)
	gl.GetFloatv(uint32(gl.POLYGON_OFFSET_UNITS), &unit_)
	return factor_, unit_, CheckGLerror()
}

// 设置polygonOffet参数
func XglPolygonOffset(factor float32, unit float32) error {
	gl.PolygonOffset(factor, unit)
	return CheckGLerror()
}
