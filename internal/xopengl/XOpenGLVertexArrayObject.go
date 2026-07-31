package xopengl

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type XOpenGLVertexArrayObject struct {
	id      uint32
	buffers map[uint32]*XOpenGLBuffer
}

func NewXOpenGLVertexArrayObjec() *XOpenGLVertexArrayObject {
	return &XOpenGLVertexArrayObject{
		id:      0,
		buffers: make(map[uint32]*XOpenGLBuffer),
	}
}

func (a *XOpenGLVertexArrayObject) Create() error {
	gl.GenVertexArrays(1, &a.id)
	return CheckGLerror()
}

func (a *XOpenGLVertexArrayObject) IsCreated() bool {
	return a.id != 0
}

func (a *XOpenGLVertexArrayObject) Destory() error {
	if a.IsCreated() {
		gl.DeleteVertexArrays(1, &a.id)
		return CheckGLerror()
	}
	return nil
}

func (a *XOpenGLVertexArrayObject) Bind() error {
	if !a.IsCreated() {
		return fmt.Errorf("XOpenGLVertexArrayObjec not created!")
	}

	gl.BindVertexArray(a.id)
	return CheckGLerror()
}

func (a *XOpenGLVertexArrayObject) Release() error {
	if !a.IsCreated() {
		return fmt.Errorf("XOpenGLVertexArrayObjec not created!")
	}
	gl.BindVertexArray(0)
	return CheckGLerror()
}

// 绑定顶点属性到顶点数组对象（VAO）
// 参数
//
//		index:顶点属性的索引
//		buffer:与该顶点属性关联的VBO缓冲区
//		size:每个顶点属性关联包含的分量数量
//		dataType:属性的底层数据类型（如 GL_FLOAT）
//		stride:步长，单位为字节（相邻属性之间的内存跨度）
//	  	offset: 偏移量，单位为字节（属性在缓冲区中的起始位置）
//	  	isInstance: 若为 true，则将该属性配置为实例化渲染（Instanced Rendering）模式
func (a *XOpenGLVertexArrayObject) AddBuffer(index uint32, buffer *XOpenGLBuffer, size uint32, dataType GLDataType, stride uint32, offset uint32, isInstance bool) error {
	a.Bind()
	buffer.Bind()
	//ptr := unsafe.Add(unsafe.Pointer(nil), uintptr(offset))
	gl.VertexAttribPointer(index, int32(size), uint32(dataType), false, int32(stride), gl.PtrOffset(int(offset))) //unsafe.Pointer(uintptr(offset))
	buffer.Release()
	a.Release()
	return CheckGLerror()
}

func (a *XOpenGLVertexArrayObject) EnableAttribute(index uint32) error {
	a.Bind()
	gl.EnableVertexAttribArray(index)
	a.Release()
	return CheckGLerror()
}

func (a *XOpenGLVertexArrayObject) DiableAttribute(index uint32) error {
	a.Bind()
	gl.DisableVertexAttribArray(index)
	a.Release()
	return CheckGLerror()
}
