//go:build darwin

package xopengl

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework OpenGL
#import <Cocoa/Cocoa.h>

static void *createOffscreenContextGo(void) {
    NSOpenGLPixelFormatAttribute attributes[] = {
        NSOpenGLPFAOpenGLProfile, NSOpenGLProfileVersion4_1Core,
        NSOpenGLPFAAccelerated,
        NSOpenGLPFADoubleBuffer,
        NSOpenGLPFAColorSize, 24,
        NSOpenGLPFADepthSize, 24,
        NSOpenGLPFAStencilSize, 8,
        0
    };
    NSOpenGLPixelFormat *format = [[NSOpenGLPixelFormat alloc] initWithAttributes:attributes];
    if (format == nil) {
        return NULL;
    }
    NSOpenGLContext *context = [[NSOpenGLContext alloc] initWithFormat:format shareContext:nil];
    [format release];
    return context;
}

static void destroyOffscreenContextGo(void *context) {
    NSOpenGLContext *openGLContext = (NSOpenGLContext *)context;
    if (openGLContext != nil) {
        if ([NSOpenGLContext currentContext] == openGLContext) {
            [NSOpenGLContext clearCurrentContext];
        }
        [openGLContext clearDrawable];
        [openGLContext release];
    }
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type XOffsetWindow struct {
	context unsafe.Pointer
}

func NewXOffsetWindow() (*XOffsetWindow, error) {
	context := C.createOffscreenContextGo()
	if context == nil {
		return nil, fmt.Errorf("failed to create offscreen OpenGL context")
	}
	return &XOffsetWindow{context: context}, nil
}

func (w *XOffsetWindow) WinId() uintptr {
	return uintptr(w.context)
}

func (w *XOffsetWindow) Close() error {
	if w.context != nil {
		C.destroyOffscreenContextGo(w.context)
		w.context = nil
	}
	return nil
}
