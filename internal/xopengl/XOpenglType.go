package xopengl

import "github.com/go-gl/gl/v4.1-core/gl"

type GLDataType uint32

const (
	GLDataType_uint   GLDataType = gl.UNSIGNED_INT
	GLDataType_uint8  GLDataType = gl.UNSIGNED_BYTE
	GLDataType_uint16 GLDataType = gl.UNSIGNED_SHORT
	GLDataType_int    GLDataType = gl.INT
	GLDataType_float  GLDataType = gl.FLOAT
	GLDataType_double GLDataType = gl.DOUBLE
)
