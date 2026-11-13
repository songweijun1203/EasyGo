//go:build darwin

package xopengl

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework OpenGL
#import <Cocoa/Cocoa.h>
#include <stdint.h>

static void makeCurrentContextGo(uintptr_t context) {
    [(NSOpenGLContext *)(void *)context makeCurrentContext];
}

static void doneCurrentContextGo(uintptr_t context) {
    if ([NSOpenGLContext currentContext] == (NSOpenGLContext *)(void *)context) {
        [NSOpenGLContext clearCurrentContext];
    }
}

static void swapBuffersGo(uintptr_t context) {
    [(NSOpenGLContext *)(void *)context flushBuffer];
}
*/
import "C"

import (
	"strconv"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	glInitOnce sync.Once
	glInitErr  error
)

// XOpenGLContext wraps the NSOpenGLContext owned by the GLFW window.
type XOpenGLContext struct {
	nativeContext uintptr
	contextMsg    string
	refCount      int
	mu            sync.Mutex
}

func NewXOpenGLContext() *XOpenGLContext {
	return &XOpenGLContext{}
}

func (c *XOpenGLContext) Create(winId uintptr) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if winId == 0 {
		return false
	}

	c.nativeContext = winId
	C.makeCurrentContextGo(C.uintptr_t(c.nativeContext))
	defer C.doneCurrentContextGo(C.uintptr_t(c.nativeContext))

	glInitOnce.Do(func() {
		glInitErr = gl.Init()
	})
	if glInitErr != nil {
		c.nativeContext = 0
		return false
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	c.contextMsg = "OpenGL版本: " + version + "\n"

	var profileMask, flags int32
	gl.GetIntegerv(gl.CONTEXT_PROFILE_MASK, &profileMask)
	gl.GetIntegerv(gl.CONTEXT_FLAGS, &flags)
	c.contextMsg += "Profile mask: " + strconv.Itoa(int(profileMask)) + "\n"
	c.contextMsg += "Context flags: " + strconv.Itoa(int(flags)) + "\n"

	if profileMask&gl.CONTEXT_CORE_PROFILE_BIT != 0 {
		c.contextMsg += "核心模式\n"
	}
	if flags&gl.CONTEXT_FLAG_FORWARD_COMPATIBLE_BIT != 0 {
		c.contextMsg += "前向兼容标志:1"
	} else {
		c.contextMsg += "前向兼容标志:0"
	}

	return true
}

func (c *XOpenGLContext) MakeCurrent() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isValidLocked() {
		return false
	}
	if c.refCount == 0 {
		C.makeCurrentContextGo(C.uintptr_t(c.nativeContext))
	}
	c.refCount++
	return true
}

func (c *XOpenGLContext) DoneCurrent() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isValidLocked() || c.refCount == 0 {
		return
	}
	c.refCount--
	if c.refCount == 0 {
		C.doneCurrentContextGo(C.uintptr_t(c.nativeContext))
	}
}

func (c *XOpenGLContext) SwapBuffers() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isValidLocked() {
		C.swapBuffersGo(C.uintptr_t(c.nativeContext))
	}
}

func (c *XOpenGLContext) GetNativeContext() uintptr {
	return c.nativeContext
}

func (c *XOpenGLContext) GetNativeDisplay() uintptr {
	return 0
}

func (c *XOpenGLContext) IsValid() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isValidLocked()
}

func (c *XOpenGLContext) isValidLocked() bool {
	return c.nativeContext != 0
}

func (c *XOpenGLContext) GetContextMsg() string {
	return c.contextMsg
}

func (c *XOpenGLContext) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.nativeContext != 0 {
		C.doneCurrentContextGo(C.uintptr_t(c.nativeContext))
	}
	c.nativeContext = 0
	c.refCount = 0
}
