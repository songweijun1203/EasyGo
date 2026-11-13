//go:build darwin

package xglwidget

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.4/glfw"
)

func configureOpenGLWindow() {
	glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLAPI)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 8)
	glfw.WindowHint(glfw.DepthBits, 24)
	glfw.WindowHint(glfw.StencilBits, 8)
	glfw.WindowHint(glfw.Resizable, glfw.True)
}

func nativeOpenGLHandle(window *glfw.Window) uintptr {
	return uintptr(unsafe.Pointer(window.GetNSGLContext()))
}
