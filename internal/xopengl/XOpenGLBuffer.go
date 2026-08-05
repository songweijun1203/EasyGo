package xopengl

import (
	xq_utils "EasyGo/internal/utils"
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// BufferType 对应 XOpenGLBuffer::Type
type BufferType uint32

const (
	BufferTypeVertexBuffer        BufferType = gl.ARRAY_BUFFER
	BufferTypeIndexBuffer         BufferType = gl.ELEMENT_ARRAY_BUFFER
	BufferTypeUniformBuffer       BufferType = gl.UNIFORM_BUFFER
	BufferTypeShaderStorageBuffer BufferType = gl.SHADER_STORAGE_BUFFER
	BufferTypePixelPackBuffer     BufferType = gl.PIXEL_PACK_BUFFER
	BufferTypeTransformFeedback   BufferType = gl.TRANSFORM_FEEDBACK_BUFFER
)

// UsagePattern 对应 XOpenGLBuffer::UsagePattern
type UsagePattern uint32

const (
	UsageStaticDraw  UsagePattern = gl.STATIC_DRAW
	UsageDynamicDraw UsagePattern = gl.DYNAMIC_DRAW
	UsageStreamDraw  UsagePattern = gl.STREAM_DRAW
)

// Access 对应 XOpenGLBuffer::Access
type Access uint32

const (
	AccessReadOnly  Access = gl.READ_ONLY
	AccessWriteOnly Access = gl.WRITE_ONLY
	AccessReadWrite Access = gl.READ_WRITE
)

// 封装VBO
type XOpenGLBuffer struct {
	id           uint32
	bufferType   BufferType
	usagePattern UsagePattern
}

// 构造
func NewXOpenGLBuffer(bt BufferType) *XOpenGLBuffer {
	return &XOpenGLBuffer{
		id:           0,
		bufferType:   bt,
		usagePattern: UsageStaticDraw,
	}
}

// 设置pattern
func (a *XOpenGLBuffer) SetUsagePattern(pattern UsagePattern) {
	a.usagePattern = pattern
}

// 获取pattern
func (a *XOpenGLBuffer) GetUsagePattern() (pattern UsagePattern) {
	return a.usagePattern
}

// 创建buffer
func (a *XOpenGLBuffer) Create() (uint32, error) {
	if a.IsCreated() {
		return a.id, nil
	}

	if a.bufferType == BufferTypeTransformFeedback {
		gl.GenTransformFeedbacks(1, &a.id)
	} else {
		gl.GenBuffers(1, &a.id)
	}

	var id uint32
	gl.GenBuffers(1, &id)
	a.id = id
	return id, CheckGLerror()
}

// 绑定buffer
func (a *XOpenGLBuffer) Bind() error {
	bindId, err := GetBindDataBufferId(a.bufferType)
	if err != nil {
		return err
	}

	if bindId == int32(a.id) {
		return nil
	}

	if a.bufferType == BufferTypeTransformFeedback {
		gl.BindTransformFeedback(uint32(a.bufferType), a.id)
		return CheckGLerror()
	} else {
		gl.BindBuffer(uint32(a.bufferType), a.id)
		return CheckGLerror()
	}
}

// release buffer
func (a *XOpenGLBuffer) Release() error {
	if !a.IsCreated() {
		return fmt.Errorf("XOpenGLBuffer not created!")
	}
	bindId, err := GetBindDataBufferId(a.bufferType)
	if err != nil {
		return err
	}

	//如果当前激活的不是当前buffer，则不做处理
	if bindId != int32(a.id) {
		return nil
	}

	if a.bufferType == BufferTypeTransformFeedback {
		gl.BindTransformFeedback(uint32(a.bufferType), 0)
		return CheckGLerror()
	} else {
		gl.BindBuffer(uint32(a.bufferType), 0)
		return CheckGLerror()
	}
}

func (a *XOpenGLBuffer) IsCreated() bool {
	return a.id != 0
}

func (a *XOpenGLBuffer) Destory() error {
	if a.IsCreated() {
		gl.DeleteBuffers(1, &a.id)
		a.id = 0
		return CheckGLerror()
	}
	return nil
}

// 分配存储空间,分配size个字节大小
func (a *XOpenGLBuffer) AllocateEmpty(size int) error {
	gl.BufferData(uint32(a.bufferType), size, nil, uint32(a.usagePattern))
	return CheckGLerror()
}

// AllocateFromData
func AllocateFromData[T xq_utils.XNumber](b *XOpenGLBuffer, data xq_utils.XDataArray1D[T]) error {
	if data.GetSize() == 0 {
		return b.AllocateEmpty(0)
	}

	size := data.GetBytesNum()

	// 将切片转为 unsafe.Pointer
	ptr := data.GetPointer()
	gl.BufferData(uint32(b.bufferType), size, ptr, uint32(b.usagePattern))
	return CheckGLerror()
}

func (a *XOpenGLBuffer) Write(offset int, count int, ptr unsafe.Pointer) error {
	if ptr == nil {
		return fmt.Errorf("XOpenGLBuffer.Write: ptr is nil")
	}
	if count <= 0 {
		return fmt.Errorf("XOpenGLBuffer.Write: count must be greater than 0, got %d", count)
	}
	if offset < 0 {
		return fmt.Errorf("XOpenGLBuffer.Write: offset cannot be negative, got %d", offset)
	}
	gl.BufferSubData(uint32(a.bufferType), offset, count, ptr)
	return CheckGLerror()
}

// 读取buffer数据
func (a *XOpenGLBuffer) Read(offset int, count int) (*xq_utils.XDataArray1D[uint8], error) {
	data := xq_utils.NewXUInt8Array1D()
	data.SetNumOfTuple(count)
	gl.GetBufferSubData(uint32(a.bufferType), offset, count, data.GetPointer())
	return data, CheckGLerror()
}
