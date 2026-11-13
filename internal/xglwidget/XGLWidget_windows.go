//go:build windows

package xglwidget

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.4/glfw"
)

func configureOpenGLWindow() {
	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	glfw.WindowHint(glfw.Resizable, glfw.True)
}

func nativeOpenGLHandle(window *glfw.Window) uintptr {
	return uintptr(unsafe.Pointer(window.GetWin32Window()))
}
