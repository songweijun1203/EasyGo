package xrender

import (
	"EasyGo/internal/xbaserender"
	"EasyGo/internal/xopengl"
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type XOpenGLRenderWindow struct {
	ctx          *xopengl.XOpenGLContext
	eventDisptch xbaserender.XBaseRenderWindowEventDispatch
}

func NewXOpenGLRenderWindow() *XOpenGLRenderWindow {
	dispatch := &XRenderWindowEventDispatch{}
	w := &XOpenGLRenderWindow{
		ctx:          nil,
		eventDisptch: dispatch,
	}
	dispatch.SetRenderWindow(w)
	return w
}

// SetWindowId 设置窗口句柄
func (renWin *XOpenGLRenderWindow) SetWindowId(winId uintptr) {
	ctx := xopengl.NewXOpenGLContext()
	if ctx.Create(winId) {
		renWin.ctx = ctx
	}
}
func (renWin *XOpenGLRenderWindow) Init() {
	//do nth
}

func (renWin *XOpenGLRenderWindow) Render() error {
	if !renWin.isCtxInit() {
		return fmt.Errorf("ctx not init!")
	}

	//渲染工作
	if !renWin.MakeCurrent() {
		return fmt.Errorf("MakeCurrent fail!")
	}

	gl.ClearColor(0.1, 0.1, 0.2, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	renWin.SwapBuffers()
	renWin.DoneCurrent()
	return nil
}

func (renWin *XOpenGLRenderWindow) MakeCurrent() bool {
	if !renWin.isCtxInit() {
		return false
	}
	return renWin.ctx.MakeCurrent()
}

func (renWin *XOpenGLRenderWindow) DoneCurrent() {
	if !renWin.isCtxInit() {
		return
	}
	renWin.ctx.DoneCurrent()
}

func (renWin *XOpenGLRenderWindow) SwapBuffers() {
	if !renWin.isCtxInit() {
		return
	}
	renWin.ctx.SwapBuffers()
}

// 渲染上下文是否创建完成
func (renWin *XOpenGLRenderWindow) isCtxInit() bool {
	if renWin.ctx == nil {
		return false
	}
	return true
}

func (renWin *XOpenGLRenderWindow) GetEventDispatcher() xbaserender.XBaseRenderWindowEventDispatch {
	return renWin.eventDisptch
}
