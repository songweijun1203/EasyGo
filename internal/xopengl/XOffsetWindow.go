package xopengl

/*
#include <windows.h>

// 窗口过程：所有离屏窗口共用
LRESULT CALLBACK offscreenWndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam) {
    return DefWindowProc(hwnd, msg, wParam, lParam);
}

// 注册窗口类（返回 ATOM，0 表示失败）
static ATOM RegisterOffscreenClass(HINSTANCE hInstance) {
    WNDCLASS wc = {0};
    wc.lpfnWndProc = offscreenWndProc;
    wc.hInstance = hInstance;
    wc.lpszClassName = "OffscreenOpenGL";
    return RegisterClass(&wc);
}

// 创建隐藏窗口
static HWND CreateOffscreenWindow(HINSTANCE hInstance) {
    HWND hwnd = CreateWindowEx(
        0,
        "OffscreenOpenGL",
        "Temp Window",
        0,          // 无样式
        0, 0, 1, 1, // 位置和尺寸
        NULL, NULL, hInstance, NULL
    );
    if (hwnd) {
        ShowWindow(hwnd, SW_HIDE);
    }
    return hwnd;
}

// 销毁窗口
static void DestroyOffscreenWindow(HWND hwnd) {
    if (hwnd) DestroyWindow(hwnd);
}
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

var (
	once    sync.Once
	hInst   C.HINSTANCE
	hwnd    C.HWND
	errInit error
)

// XOffsetWindow 封装一个隐藏的离屏窗口（单例模式）
type XOffsetWindow struct {
	hwnd C.HWND
}

// NewXOffsetWindow 返回单例隐藏窗口（线程安全）
func NewXOffsetWindow() (*XOffsetWindow, error) {
	once.Do(func() {
		// 获取当前模块句柄
		hInst = C.GetModuleHandle(nil)
		if hInst == nil {
			errInit = fmt.Errorf("GetModuleHandle failed")
			return
		}
		// 注册窗口类
		atom := C.RegisterOffscreenClass(hInst)
		if atom == 0 {
			errInit = fmt.Errorf("RegisterClass failed")
			return
		}
		// 创建窗口
		hwnd = C.CreateOffscreenWindow(hInst)
		if hwnd == nil {
			errInit = fmt.Errorf("CreateWindowEx failed")
			return
		}
	})
	if errInit != nil {
		return nil, errInit
	}
	return &XOffsetWindow{hwnd: hwnd}, nil
}

// WinId 返回窗口句柄（可用于 OpenGL 上下文创建）
func (w *XOffsetWindow) WinId() uintptr {
	return uintptr(unsafe.Pointer(w.hwnd))
}

// Close 销毁窗口（注意：如需注销类，需确保没有其他窗口使用）
func (w *XOffsetWindow) Close() error {
	if w.hwnd != nil {
		C.DestroyOffscreenWindow(w.hwnd)
		w.hwnd = nil
		// 注意：由于类可能被多个窗口共享，通常不在这里注销
		// 如果需要彻底清理，可以在所有窗口销毁后调用 UnregisterClass
	}
	return nil
}
