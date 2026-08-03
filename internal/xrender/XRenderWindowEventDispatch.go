package xrender

import (
	"EasyGo/internal/xbaserender"
)

type XRenderWindowEventDispatch struct {
	xbaserender.XRenderWindowEventDispatchSig                               //信号
	renderWindow                              xbaserender.XBaseRenderWindow //渲染窗口
}

func (d *XRenderWindowEventDispatch) SetRenderWindow(ren xbaserender.XBaseRenderWindow) {
	d.renderWindow = ren
}

func (d *XRenderWindowEventDispatch) GetRenderWindow() xbaserender.XBaseRenderWindow {
	return d.renderWindow
}

func (d *XRenderWindowEventDispatch) GetXRenderWindowEventDispatchSig() *xbaserender.XRenderWindowEventDispatchSig {
	return &d.XRenderWindowEventDispatchSig
}

/*
func (d *XRenderWindowEventDispatch) GetSigLeftButtonPress() *gsignal.Event[xbaserender.MouseButtonEvent] {
	return &d.SigLeftButtonPress
}

func (d *XRenderWindowEventDispatch) GetSigLeftButtonRelease() *gsignal.Event[xbaserender.MouseButtonEvent] {
	return &d.SigLeftButtonRelease
}

func (d *XRenderWindowEventDispatch) GetSigLeftButtonDoublePress() *gsignal.Event[xbaserender.MouseButtonEvent] {
	return &d.SigLeftButtonDoublePress
}

func (d *XRenderWindowEventDispatch) GetSigMiddleButtonPress() *gsignal.Event[xbaserender.MouseButtonEvent] {
	return &d.SigMiddleButtonPress
}

func (d *XRenderWindowEventDispatch) GetSigMiddleButtonRelease() *gsignal.Event[xbaserender.MouseButtonEvent] {
	return &d.SigMiddleButtonRelease
}

func (d *XRenderWindowEventDispatch) GetSigMiddleButtonDoublePress() *gsignal.Event[xbaserender.MouseButtonEvent] {
	return &d.SigMiddleButtonDoublePress
}

func (d *XRenderWindowEventDispatch) GetSigRightButtonPress() *gsignal.Event[xbaserender.MouseButtonEvent] {
	return &d.SigRightButtonPress
}

func (d *XRenderWindowEventDispatch) GetSigRightButtonRelease() *gsignal.Event[xbaserender.MouseButtonEvent] {
	return &d.SigRightButtonRelease
}

func (d *XRenderWindowEventDispatch) GetSigRightButtonDoublePress() *gsignal.Event[xbaserender.MouseButtonEvent] {
	return &d.SigRightButtonDoublePress
}

func (d *XRenderWindowEventDispatch) GetSigResize() *gsignal.Event[xutils.Vec2i] {
	return &d.SigResize
}

func (d *XRenderWindowEventDispatch) GetSigKeyPress() *gsignal.Event[xbaserender.KeyEvent] {
	return &d.SigKeyPress
}

func (d *XRenderWindowEventDispatch) GetSigKeyRelease() *gsignal.Event[xbaserender.KeyEvent] {
	return &d.SigKeyRelease
}

func (d *XRenderWindowEventDispatch) GetSigMouseMove() *gsignal.Event[xbaserender.KeyEvent] {
	return &d.SigMouseMove
}

func (d *XRenderWindowEventDispatch) GetSigMouseWheelForward() *gsignal.Event[xbaserender.KeyEvent] {
	return &d.SigMouseWheelForward
}

func (d *XRenderWindowEventDispatch) GetSigMouseWheelBackward() *gsignal.Event[xbaserender.KeyEvent] {
	return &d.SigMouseWheelBackward
}
*/
