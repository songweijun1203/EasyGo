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

func (w *XGLWidget) OnMouseButton(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("%p OnMouseButton\n", w)
	//转发信号给EventDispatch
	dispatch := w.renderWindow.GetEventDispatcher()
	var sig *gsignal.Event[xbaserender.MouseButtonEvent]
	x, y := w.window.GetCursorPos()
	pos := xutils.Vec2i{X: int(x), Y: int(y)} // 根据需要取整
	switch button {
	case glfw.MouseButton1:
		{
			switch action {
			case glfw.Press:
				sig = dispatch.GetSigLeftButtonPress()
			case glfw.Release:
				sig = dispatch.GetSigLeftButtonRelease()
			case glfw.Repeat:
				sig = dispatch.GetSigLeftButtonDoublePress()
			default:
				return // 未知动作，忽略
			}
		}
	case glfw.MouseButton2:
		{
			switch action {
			case glfw.Press:
				sig = dispatch.GetSigRightButtonPress()
			case glfw.Release:
				sig = dispatch.GetSigRightButtonRelease()
			case glfw.Repeat:
				sig = dispatch.GetSigRightButtonDoublePress()
			default:
				return // 未知动作，忽略
			}
		}
	case glfw.MouseButton3:
		{
			switch action {
			case glfw.Press:
				sig = dispatch.GetSigMiddleButtonPress()
			case glfw.Release:
				sig = dispatch.GetSigMiddleButtonRelease()
			case glfw.Repeat:
				sig = dispatch.GetSigMiddleButtonDoublePress()
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

	event := xbaserender.MouseButtonEvent{Modifier: xbaserender.NoModifier, Pos: pos}
	sig.Emit(event)
}

// type MouseButtonCallback func(w *Window, button MouseButton, action Action, mods ModifierKey)
func (win *XGLWidget) setupCallbacks() {
	win.window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		win.OnMouseButton(button, action, mods)
	})
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
