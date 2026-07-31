package xbaserender

type XBaseRenderWindow interface {
	SetWindowId(winId uintptr)
	Render() error
	MakeCurrent() bool
	DoneCurrent()
	SwapBuffers()
	Init()
	GetEventDispatcher() XBaseRenderWindowEventDispatch
}
