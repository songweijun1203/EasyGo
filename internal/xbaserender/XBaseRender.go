package xbaserender

type XBaseRender interface {
	SetRenderWindow(XBaseRenderWindow)
	GetRenderWindow() XBaseRenderWindow
	Render()
	MakeCurrent() bool
	DoneCurrent()
	Close()
}
