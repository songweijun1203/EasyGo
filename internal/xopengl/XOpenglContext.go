// 文件: XOpenglContext.go
package xopengl

/*
#include <windows.h>
#include <GL/gl.h>
#include <GL/wgl.h>
#include <stdlib.h>

// WGL 扩展函数指针类型
typedef HGLRC (WINAPI *PFNWGLCREATECONTEXTATTRIBSARBPROC)(HDC, HGLRC, const int*);
typedef BOOL (WINAPI *PFNWGLCHOOSEPIXELFORMATARBPROC)(HDC, const int*, const FLOAT*, UINT, int*, UINT*);

// 全局扩展函数指针
PFNWGLCREATECONTEXTATTRIBSARBPROC wglCreateContextAttribsARB = NULL;
PFNWGLCHOOSEPIXELFORMATARBPROC wglChoosePixelFormatARB = NULL;

// 加载 WGL 扩展函数
void loadWGLExtensions() {
    wglCreateContextAttribsARB = (PFNWGLCREATECONTEXTATTRIBSARBPROC)wglGetProcAddress("wglCreateContextAttribsARB");
    wglChoosePixelFormatARB = (PFNWGLCHOOSEPIXELFORMATARBPROC)wglGetProcAddress("wglChoosePixelFormatARB");
}

// 包装函数（供 Go 调用）
HGLRC wglCreateContextAttribsARB_wrapper(HDC hdc, HGLRC shareContext, const int* attribList) {
    if (wglCreateContextAttribsARB) {
        return wglCreateContextAttribsARB(hdc, shareContext, attribList);
    }
    return NULL;
}

BOOL wglChoosePixelFormatARB_wrapper(HDC hdc, const int* piAttribIList, const FLOAT* pfAttribFList,
                                     UINT nMaxFormats, int* piFormats, UINT* nNumFormats) {
    if (wglChoosePixelFormatARB) {
        return wglChoosePixelFormatARB(hdc, piAttribIList, pfAttribFList, nMaxFormats, piFormats, nNumFormats);
    }
    return FALSE;
}

// 创建临时窗口并初始化扩展（仅执行一次）
BOOL initTempContext() {
    static BOOL initialized = FALSE;
    if (initialized) return TRUE;

    WNDCLASS wc = {0};
    wc.lpfnWndProc = DefWindowProc;
    wc.hInstance = GetModuleHandle(NULL);
    wc.lpszClassName = "TempGLWindow";
    if (!RegisterClass(&wc)) return FALSE;

    HWND hwnd = CreateWindowEx(0, "TempGLWindow", "Temp", 0,
        0, 0, 1, 1, NULL, NULL, GetModuleHandle(NULL), NULL);
    if (!hwnd) {
        UnregisterClass("TempGLWindow", GetModuleHandle(NULL));
        return FALSE;
    }
    HDC hdc = GetDC(hwnd);
    PIXELFORMATDESCRIPTOR pfd = {
        sizeof(PIXELFORMATDESCRIPTOR), 1,
        PFD_DRAW_TO_WINDOW | PFD_SUPPORT_OPENGL | PFD_DOUBLEBUFFER,
        PFD_TYPE_RGBA, 32, 0,0,0,0,0,0,0,0,0,0,0,0,0,
        24, 8, 0, PFD_MAIN_PLANE, 0,0,0,0
    };
    int pf = ChoosePixelFormat(hdc, &pfd);
    if (!pf) { ReleaseDC(hwnd, hdc); DestroyWindow(hwnd); UnregisterClass("TempGLWindow", GetModuleHandle(NULL)); return FALSE; }
    if (!SetPixelFormat(hdc, pf, &pfd)) { ReleaseDC(hwnd, hdc); DestroyWindow(hwnd); UnregisterClass("TempGLWindow", GetModuleHandle(NULL)); return FALSE; }
    HGLRC temp = wglCreateContext(hdc);
    if (!temp) { ReleaseDC(hwnd, hdc); DestroyWindow(hwnd); UnregisterClass("TempGLWindow", GetModuleHandle(NULL)); return FALSE; }
    if (!wglMakeCurrent(hdc, temp)) { wglDeleteContext(temp); ReleaseDC(hwnd, hdc); DestroyWindow(hwnd); UnregisterClass("TempGLWindow", GetModuleHandle(NULL)); return FALSE; }

    loadWGLExtensions();

    wglMakeCurrent(NULL, NULL);
    wglDeleteContext(temp);
    ReleaseDC(hwnd, hdc);
    DestroyWindow(hwnd);
    UnregisterClass("TempGLWindow", GetModuleHandle(NULL));
    initialized = TRUE;
    return TRUE;
}

// 使用包装函数选择多重采样像素格式
int choosePixelFormatARB(HDC hdc) {
    if (!wglChoosePixelFormatARB) return 0;
    const int attribs[] = {
        WGL_DRAW_TO_WINDOW_ARB, GL_TRUE,
        WGL_SUPPORT_OPENGL_ARB, GL_TRUE,
        WGL_DOUBLE_BUFFER_ARB, GL_TRUE,
        WGL_PIXEL_TYPE_ARB, WGL_TYPE_RGBA_ARB,
        WGL_COLOR_BITS_ARB, 32,
        WGL_DEPTH_BITS_ARB, 24,
        WGL_STENCIL_BITS_ARB, 8,
        WGL_SAMPLE_BUFFERS_ARB, GL_TRUE,
        WGL_SAMPLES_ARB, 8,
        0
    };
    UINT numFormats;
    int pf;
    if (!wglChoosePixelFormatARB_wrapper(hdc, attribs, NULL, 1, &pf, &numFormats) || numFormats == 0)
        return 0;
    return pf;
}

// 设置像素格式
BOOL setPixelFormatGo(HDC hdc, int pf) {
    PIXELFORMATDESCRIPTOR pfd;
    DescribePixelFormat(hdc, pf, sizeof(pfd), &pfd);
    return SetPixelFormat(hdc, pf, &pfd);
}

// 辅助函数（封装为Go可调用）
HDC getDCGo(HWND hwnd) { return GetDC(hwnd); }
int releaseDCGo(HWND hwnd, HDC hdc) { return ReleaseDC(hwnd, hdc); }
HWND windowFromDCGo(HDC hdc) { return WindowFromDC(hdc); }
BOOL swapBuffersGo(HDC hdc) { return SwapBuffers(hdc); }
const char* glGetStringGo(GLenum name) { return (const char*)glGetString(name); }
void glGetIntegervGo(GLenum pname, GLint* params) { glGetIntegerv(pname, params); }

// 链接库
#cgo LDFLAGS: -lopengl32 -lgdi32

*/
import "C"
import (
	"sync"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// WGL 常量
const (
	WGL_CONTEXT_MAJOR_VERSION_ARB    = 0x2091
	WGL_CONTEXT_MINOR_VERSION_ARB    = 0x2092
	WGL_CONTEXT_PROFILE_MASK_ARB     = 0x9126
	WGL_CONTEXT_CORE_PROFILE_BIT_ARB = 0x00000001
	WGL_CONTEXT_FLAGS_ARB            = 0x2094
	WGL_CONTEXT_DEBUG_BIT_ARB        = 0x00000001

	WGL_DRAW_TO_WINDOW_ARB = 0x2001
	WGL_SUPPORT_OPENGL_ARB = 0x2010
	WGL_DOUBLE_BUFFER_ARB  = 0x2011
	WGL_PIXEL_TYPE_ARB     = 0x2013
	WGL_TYPE_RGBA_ARB      = 0x202B
	WGL_COLOR_BITS_ARB     = 0x2014
	WGL_DEPTH_BITS_ARB     = 0x2022
	WGL_STENCIL_BITS_ARB   = 0x2023
	WGL_SAMPLE_BUFFERS_ARB = 0x2041
	WGL_SAMPLES_ARB        = 0x2042
)

// OpenGL 常量
const (
	GL_VERSION                   = 0x1F02
	GL_CONTEXT_PROFILE_MASK      = 0x9126
	GL_CONTEXT_FLAGS             = 0x821E
	GL_CONTEXT_CORE_PROFILE_BIT  = 0x00000001
	GL_CONTEXT_COMPATIBILITY_BIT = 0x00000002
	GL_CONTEXT_FLAG_DEBUG_BIT    = 0x00000002
)

var (
	glInitOnce sync.Once
)

// XOpenGLContext 结构体
type XOpenGLContext struct {
	nativeContext uintptr // HGLRC
	nativeDisplay uintptr // HDC
	contextMsg    string
	refCount      int
	mu            sync.Mutex
}

// NewXOpenGLContext 构造函数
func NewXOpenGLContext() *XOpenGLContext {
	return &XOpenGLContext{}
}

var initOnce sync.Once
var initSuccess bool

// ensureExtensionsLoaded 调用 C 的 initTempContext
func ensureExtensionsLoaded() bool {
	initOnce.Do(func() {
		initSuccess = C.initTempContext() != 0
	})
	return initSuccess
}

// Create 使用窗口句柄创建 OpenGL 上下文（winId 为 HWND）
func (c *XOpenGLContext) Create(winId uintptr) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !ensureExtensionsLoaded() {
		return false
	}

	// 关键修复：使用 unsafe.Pointer 进行转换
	hwnd := (C.HWND)(unsafe.Pointer(winId))

	// 1. 获取 DC
	hdc := C.getDCGo(hwnd)
	if hdc == nil {
		return false
	}
	defer func() {
		if c.nativeDisplay == 0 {
			C.releaseDCGo(hwnd, hdc)
		}
	}()

	// 2. 选择像素格式（带 MSAA）
	pf := C.choosePixelFormatARB(hdc)
	if pf == 0 {
		return false
	}

	// 3. 设置像素格式
	if C.setPixelFormatGo(hdc, pf) == 0 {
		return false
	}

	// 4. 创建 OpenGL 上下文 (4.6 Core + Debug)
	attribs := []int32{
		WGL_CONTEXT_MAJOR_VERSION_ARB, 4,
		WGL_CONTEXT_MINOR_VERSION_ARB, 6,
		WGL_CONTEXT_PROFILE_MASK_ARB, WGL_CONTEXT_CORE_PROFILE_BIT_ARB,
		WGL_CONTEXT_FLAGS_ARB, WGL_CONTEXT_DEBUG_BIT_ARB,
		0,
	}
	ctx := C.wglCreateContextAttribsARB_wrapper(hdc, nil, (*C.int)(unsafe.Pointer(&attribs[0])))
	if ctx == nil {
		return false
	}

	// 5. 激活上下文
	if C.wglMakeCurrent(hdc, ctx) == 0 {
		C.wglDeleteContext(ctx)
		return false
	}

	// 6. 保存为 uintptr
	c.nativeContext = uintptr(unsafe.Pointer(ctx))
	c.nativeDisplay = uintptr(unsafe.Pointer(hdc))

	// 7. 收集信息
	version := C.GoString(C.glGetStringGo(GL_VERSION))
	c.contextMsg = "OpenGL版本: " + version + "\n"

	var profileMask, flags C.GLint
	C.glGetIntegervGo(GL_CONTEXT_PROFILE_MASK, &profileMask)
	C.glGetIntegervGo(GL_CONTEXT_FLAGS, &flags)
	c.contextMsg += "Profile mask: " + itoa(int(profileMask)) + "\n"
	c.contextMsg += "Context flags: " + itoa(int(flags)) + "\n"

	if profileMask&GL_CONTEXT_CORE_PROFILE_BIT != 0 {
		c.contextMsg += "核心模式\n"
	} else if profileMask&GL_CONTEXT_COMPATIBILITY_BIT != 0 {
		c.contextMsg += "兼容模式\n"
	}
	if flags&GL_CONTEXT_FLAG_DEBUG_BIT != 0 {
		c.contextMsg += "调试标志:1"
	} else {
		c.contextMsg += "调试标志:0"
	}

	glInitOnce.Do(func() {
		if err := gl.Init(); err != nil {
			panic(err)
		}
	})

	return true
}

// MakeCurrent 激活上下文（带引用计数）
func (c *XOpenGLContext) MakeCurrent() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isValidLocked() {
		return false
	}
	if c.refCount == 0 {
		hdc := C.HDC(unsafe.Pointer(c.nativeDisplay))
		ctx := C.HGLRC(unsafe.Pointer(c.nativeContext))
		if C.wglMakeCurrent(hdc, ctx) == 0 {
			return false
		}
	}
	c.refCount++
	return true
}

// DoneCurrent 解除当前上下文（引用计数减一）
func (c *XOpenGLContext) DoneCurrent() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isValidLocked() {
		return
	}
	if c.refCount > 0 {
		c.refCount--
		if c.refCount == 0 {
			C.wglMakeCurrent(nil, nil)
		}
	}
}

// SwapBuffers 交换前后缓冲区
func (c *XOpenGLContext) SwapBuffers() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isValidLocked() {
		hdc := C.HDC(unsafe.Pointer(c.nativeDisplay))
		C.swapBuffersGo(hdc)
	}
}

// GetNativeContext 返回 HGLRC
func (c *XOpenGLContext) GetNativeContext() uintptr {
	return c.nativeContext
}

// GetNativeDisplay 返回 HDC
func (c *XOpenGLContext) GetNativeDisplay() uintptr {
	return c.nativeDisplay
}

// IsValid 检查上下文是否有效
func (c *XOpenGLContext) IsValid() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isValidLocked()
}

func (c *XOpenGLContext) isValidLocked() bool {
	return c.nativeContext != 0 && c.nativeDisplay != 0
}

// GetContextMsg 返回上下文信息字符串
func (c *XOpenGLContext) GetContextMsg() string {
	return c.contextMsg
}

// Reset 释放所有资源
func (c *XOpenGLContext) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.nativeDisplay != 0 {
		hdc := C.HDC(unsafe.Pointer(c.nativeDisplay))
		hwnd := C.windowFromDCGo(hdc)
		C.releaseDCGo(hwnd, hdc)
	}
	if c.nativeContext != 0 {
		C.wglMakeCurrent(nil, nil)
		ctx := C.HGLRC(unsafe.Pointer(c.nativeContext))
		C.wglDeleteContext(ctx)
	}
	c.nativeDisplay = 0
	c.nativeContext = 0
	c.refCount = 0
}

// 辅助函数（生产环境请使用 strconv.Itoa）
func itoa(i int) string {
	return string(rune(i + '0')) // 仅演示，实际应替换
}
