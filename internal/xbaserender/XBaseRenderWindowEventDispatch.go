package xbaserender

import (
	"EasyGo/internal/xutils"

	"github.com/quasilyte/gsignal"
)

type MouseButtonEvent struct {
	Modifier KeyboardModifier
	Pos      xutils.Vec2i
}

type KeyEvent struct {
	Modifier KeyboardModifier
	Key      KeyType
}

type XEvent struct {
	isStopPropagate bool
}

func (e *XEvent) StopPropagate() {
	e.isStopPropagate = true //不再传播
}

func (e *XEvent) IsStopPropagate() bool {
	return e.isStopPropagate
}

type XBaseRenderWindowEventDispatch interface {
	SetRenderWindow(ren XBaseRenderWindow)
	GetRenderWindow() XBaseRenderWindow
	GetXRenderWindowEventDispatchSig() *XRenderWindowEventDispatchSig
}

type XRenderWindowEventDispatchSig struct {
	SigLeftButtonPress       gsignal.Event[MouseButtonEvent] //鼠标左键摁下
	SigLeftButtonRelease     gsignal.Event[MouseButtonEvent] //鼠标左键弹起
	SigLeftButtonDoublePress gsignal.Event[MouseButtonEvent] //鼠标左键双击

	SigMiddleButtonPress       gsignal.Event[MouseButtonEvent] //鼠标中键摁下
	SigMiddleButtonRelease     gsignal.Event[MouseButtonEvent] //鼠标中键弹起
	SigMiddleButtonDoublePress gsignal.Event[MouseButtonEvent] //鼠标中键双击

	SigRightButtonPress       gsignal.Event[MouseButtonEvent] //鼠标右键摁下
	SigRightButtonRelease     gsignal.Event[MouseButtonEvent] //鼠标右键弹起
	SigRightButtonDoublePress gsignal.Event[MouseButtonEvent] //鼠标右键双击

	SigResize gsignal.Event[xutils.Vec2i] //窗口大小调整

	SigKeyPress           gsignal.Event[KeyEvent]         //键盘摁下
	SigKeyRelease         gsignal.Event[KeyEvent]         //键盘弹起
	SigMouseMove          gsignal.Event[MouseButtonEvent] //鼠标移动
	SigMouseWheelForward  gsignal.Event[MouseButtonEvent] //鼠标前进
	SigMouseWheelBackward gsignal.Event[MouseButtonEvent] //鼠标后退
}
