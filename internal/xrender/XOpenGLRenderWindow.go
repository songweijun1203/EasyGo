package xrender

import (
	xutils "EasyGo/internal/utils"
	"EasyGo/internal/xbaserender"
	"EasyGo/internal/xopengl"
	"EasyGo/internal/xsignal"
	"fmt"
	"sync"
)

type XOpenGLRenderWindow struct {
	ctx          *xopengl.XOpenGLContext
	eventDisptch xbaserender.XBaseRenderWindowEventDispatch
	renders      []xbaserender.XBaseRender
	mu           sync.Mutex          // 保护 renders
	connector    *xsignal.XConnector //信号管理
	windowSize   xutils.Vec2i
}

func NewXOpenGLRenderWindow() *XOpenGLRenderWindow {
	dispatch := &XRenderWindowEventDispatch{}
	w := &XOpenGLRenderWindow{
		ctx:          nil,
		eventDisptch: dispatch,
		renders:      make([]xbaserender.XBaseRender, 0),
		connector:    xsignal.NewXConnector(),
	}
	dispatch.SetRenderWindow(w)
	w.Init()
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
	renWin.bindSignalSlot()
}

func (renWin *XOpenGLRenderWindow) GetWindowWidth() int32 {
	return renWin.windowSize.X
}
func (renWin *XOpenGLRenderWindow) GetWindowHeigth() int32 {
	return renWin.windowSize.Y
}

func (renWin *XOpenGLRenderWindow) SlotResize(s xutils.Vec2i) {
	//记录窗口大小
	renWin.windowSize = s
}

func (renWin *XOpenGLRenderWindow) bindSignalSlot() {
	sig := renWin.GetEventDispatcher().GetXRenderWindowEventDispatchSig()
	xsignal.ConnectSlot(renWin.connector, &sig.SigResize, renWin.SlotResize, "SlotResize")
}

func (renWin *XOpenGLRenderWindow) Render() error {
	if !renWin.isCtxInit() {
		return fmt.Errorf("ctx not init!")
	}

	//渲染工作
	if !renWin.MakeCurrent() {
		return fmt.Errorf("MakeCurrent fail!")
	}

	for _, render := range renWin.renders {
		render.Render()
	}

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

func (renWin *XOpenGLRenderWindow) AddRender(render xbaserender.XBaseRender) bool {
	if render == nil {
		return false
	}

	renWin.mu.Lock()
	defer renWin.mu.Unlock()

	// 遍历现有列表，直接比较接口值（适用于指针类型）
	for _, r := range renWin.renders {
		if r == render {
			return false // 已存在，不重复添加
		}
	}

	renWin.renders = append(renWin.renders, render)
	render.SetRenderWindow(renWin)                   //设置渲染窗口
	render.(*XRender).connectToRenderWindowSignals() //连接信号
	return true
}
func (renWin *XOpenGLRenderWindow) GetRenders() []xbaserender.XBaseRender {
	return renWin.renders
}

// 渲染窗口关闭的时候，需要关闭opengl相关资源
func (renWin *XOpenGLRenderWindow) Close() {
	defer func() {
		//渲染器资源清理
		renders := renWin.GetRenders()
		for _, render := range renders {
			render.Close()
		}

		//上下文资源清理
		if renWin.ctx != nil {
			renWin.ctx.Close()
			renWin.ctx = nil
		}

	}()

}
