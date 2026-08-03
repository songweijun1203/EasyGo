package xglwidget

import (
	xutils "EasyGo/internal/utils"
	"EasyGo/internal/xbaserender"
	"EasyGo/internal/xrender"
	"fmt"
	"sync"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.4/glfw"
	"github.com/quasilyte/gsignal"
)

var (
	glfwInitOnce sync.Once
	glInitOnce   sync.Once
)

type XGLWidget struct {
	window       *glfw.Window                  //内部底层窗口
	renderWindow xbaserender.XBaseRenderWindow //接口类型
	pos          xutils.Vec2i                  //鼠标的位置
	mod          xbaserender.KeyboardModifier  //辅助摁键
}

func NewXGLWidget() (*XGLWidget, error) {
	//glfw的初始化应该只调用一次
	glfwInitOnce.Do(func() {
		if err := glfw.Init(); err != nil {
			panic(err) // 或返回错误，但这里用 panic 简化，也可以使用错误变量记录
		}
	})

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	//创建窗口
	w, err := glfw.CreateWindow(400, 300, "XGLWidgte", nil, nil)
	if err != nil {
		glfw.Terminate()
		return nil, err
	}

	// 获取 Windows HWND
	hwndRaw := w.GetWin32Window()
	hwnd := uintptr(unsafe.Pointer(hwndRaw))
	if hwnd == 0 {
		w.Destroy()
		glfw.Terminate()
		return nil, fmt.Errorf("failed to get HWND")
	}

	win := &XGLWidget{
		window:       w,
		renderWindow: xrender.NewXOpenGLRenderWindow(),
		pos:          xutils.Vec2i{X: 0, Y: 0},
		mod:          xbaserender.NoModifier,
	}

	win.renderWindow.SetWindowId(hwnd)
	//加载opengl函数指针
	if win.renderWindow.MakeCurrent() {
		glInitOnce.Do(func() {
			if err := gl.Init(); err != nil {
				panic(err)
			}
		})
		win.renderWindow.DoneCurrent()
	}

	win.setupCallbacks()
	return win, nil
}

func (w *XGLWidget) GetRenderWindow() xbaserender.XBaseRenderWindow {
	return w.renderWindow
}

// 鼠标动作事件回调
func (w *XGLWidget) OnMouseButton(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	//转发信号给EventDispatch
	dispatch := w.renderWindow.GetEventDispatcher()
	var sig *gsignal.Event[xbaserender.MouseButtonEvent]
	x, y := w.window.GetCursorPos()
	w.pos = xutils.Vec2i{X: int(x), Y: int(y)}
	//pos := xutils.Vec2i{X: int(x), Y: int(y)} // 根据需要取整
	switch button {
	case glfw.MouseButton1:
		{
			switch action {
			case glfw.Press:
				sig = &dispatch.GetXRenderWindowEventDispatchSig().SigLeftButtonPress
			case glfw.Release:
				sig = &dispatch.GetXRenderWindowEventDispatchSig().SigLeftButtonRelease
			case glfw.Repeat:
				sig = &dispatch.GetXRenderWindowEventDispatchSig().SigLeftButtonDoublePress
			default:
				return // 未知动作，忽略
			}
		}
	case glfw.MouseButton2:
		{
			switch action {
			case glfw.Press:
				sig = &dispatch.GetXRenderWindowEventDispatchSig().SigRightButtonPress
			case glfw.Release:
				sig = &dispatch.GetXRenderWindowEventDispatchSig().SigRightButtonRelease
			case glfw.Repeat:
				sig = &dispatch.GetXRenderWindowEventDispatchSig().SigRightButtonDoublePress
			default:
				return // 未知动作，忽略
			}
		}
	case glfw.MouseButton3:
		{
			switch action {
			case glfw.Press:
				sig = &dispatch.GetXRenderWindowEventDispatchSig().SigMiddleButtonPress
			case glfw.Release:
				sig = &dispatch.GetXRenderWindowEventDispatchSig().SigMiddleButtonRelease
			case glfw.Repeat:
				sig = &dispatch.GetXRenderWindowEventDispatchSig().SigMiddleButtonDoublePress
			default:
				return // 未知动作，忽略
			}
		}
	case glfw.MouseButton4:
		//btnType = xbaserender.MouseButton4 // 侧键
	case glfw.MouseButton5:
		//btnType = xbaserender.MouseButton5
	default:
		// btnType = xbaserender.MouseUnknown
	}

	if sig == nil {
		return
	}

	event := xbaserender.MouseButtonEvent{Modifier: w.mod, Pos: w.pos}
	sig.Emit(event)
}

// 鼠标移动事件回调
func (w *XGLWidget) OnMouseMove(x, y float64) {
	w.pos = xutils.Vec2i{X: int(x), Y: int(y)}
	dispatch := w.renderWindow.GetEventDispatcher()
	sig := &dispatch.GetXRenderWindowEventDispatchSig().SigMouseMove
	event := xbaserender.MouseButtonEvent{Modifier: w.mod, Pos: xutils.Vec2i{X: int(x), Y: int(y)}}
	sig.Emit(event)
}

// 滚轮滚动事件
// 参数
//
//	xoff 正值—>向右滚动，负值—>向左滚动（带横向滚动的滚轮）
//	yoff 正值—>向上滚动，负值—>向下滚动
func (w *XGLWidget) OnMouseWheel(xoff, yoff float64) {
	dispatch := w.renderWindow.GetEventDispatcher()
	var sig *gsignal.Event[xbaserender.MouseButtonEvent]
	if yoff > 0 {
		sig = &dispatch.GetXRenderWindowEventDispatchSig().SigMouseWheelForward
	} else if yoff < 0 {
		sig = &dispatch.GetXRenderWindowEventDispatchSig().SigMouseWheelBackward
	}

	//需要鼠标的位置
	event := xbaserender.MouseButtonEvent{Modifier: w.mod, Pos: w.pos}
	sig.Emit(event)
}

// 窗口大小变化回调
func (w *XGLWidget) OnWindowResize(width, height int) {
	dispatch := w.renderWindow.GetEventDispatcher()
	sig := &dispatch.GetXRenderWindowEventDispatchSig().SigResize
	sig.Emit(xutils.Vec2i{X: width, Y: height})
}

// 键盘事件回调
// 参数
//
//	key 被摁下的键
//	scancode 物理键盘扫描码，通常忽略
//	action 动作类型 按下 释放 重复
//	mods 按下的修饰键（Shift、Ctrl、Alt、Super）的位掩码，可用 & 检查
func (w *XGLWidget) OnKey(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	myKey := convertKey(key)
	myMode := convertMods(mods)
	dispatch := w.renderWindow.GetEventDispatcher()
	var sig *gsignal.Event[xbaserender.KeyEvent]
	key_event := xbaserender.KeyEvent{Key: myKey, Modifier: myMode}
	w.mod = myMode
	switch action {
	case glfw.Press:
		sig = &dispatch.GetXRenderWindowEventDispatchSig().SigKeyPress
	case glfw.Release:
		sig = &dispatch.GetXRenderWindowEventDispatchSig().SigKeyRelease
	case glfw.Repeat:
		sig = nil
	}

	if sig == nil {
		return
	}

	sig.Emit(key_event)
}

// 回调转发
func (win *XGLWidget) setupCallbacks() {
	win.window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		win.OnMouseButton(button, action, mods)
	})

	win.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		win.OnKey(key, scancode, action, mods)
	})

	win.window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		win.OnMouseMove(xpos, ypos)
	})

	win.window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		win.OnWindowResize(width, height)
	})

	win.window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
		win.OnMouseWheel(xoff, yoff)
	})
}

// 获取键盘摁键
func (w *XGLWidget) getCurrentMods() xbaserender.KeyboardModifier {
	//w.window.GetKey()
	return xbaserender.NoModifier
}

func (win *XGLWidget) Run() {
	defer win.Close()

	for !win.window.ShouldClose() {
		win.Render()
		glfw.PollEvents()
	}
}

// Close 清理资源
func (win *XGLWidget) Close() {
	if win.window != nil {
		win.window.Destroy()
	}
	glfw.Terminate()
}

// ShouldClose 返回窗口是否应关闭
func (win *XGLWidget) ShouldClose() bool {
	return win.window.ShouldClose()
}

// SetShouldClose 设置关闭标志
func (win *XGLWidget) SetShouldClose(close bool) {
	win.window.SetShouldClose(close)
}

// Render 渲染场景
func (win *XGLWidget) Render() {
	win.renderWindow.Render()
}
