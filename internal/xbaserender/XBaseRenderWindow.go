package xbaserender

import xutils "EasyGo/internal/utils"

type XBaseRenderWindow interface {
	SetWindowId(winId uintptr)
	Render() error
	MakeCurrent() bool
	DoneCurrent()
	SwapBuffers()
	Init()
	GetEventDispatcher() XBaseRenderWindowEventDispatch
	AddRender(XBaseRender) bool
	GetRenders() []XBaseRender
	Close()

	SlotResize(xutils.Vec2i) //外部窗口变化时，记录窗口大小
	GetWindowWidth() int32
	GetWindowHeigth() int32
}
