package xbaserender

type XBaseRender interface {
	SetViewPort(x float32, y float32, w float32, h float32)
	GetViewPort() (x float32, y float32, w float32, h float32)
	GetConvertViewPort() (x, y, z, w int32)

	SetRenderWindow(XBaseRenderWindow)
	GetRenderWindow() XBaseRenderWindow
	Render()
	MakeCurrent() bool
	DoneCurrent()
	Close()
}
