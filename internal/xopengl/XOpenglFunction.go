package xopengl

import (
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
func xglGetBindDataBufferId(bufferType BufferType) (id int32, err error) {
	bindValue := getBindingEnum(bufferType)
	var bindId int32 = 0
	gl.GetIntegerv(uint32(bindValue), &bindId)
	return bindId, CheckGLerror()
}
