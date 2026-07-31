package xbaserender

import (
	xutils "EasyGo/internal/utils"

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

type XBaseRenderWindowEventDispatch interface {
	SetRenderWindow(ren XBaseRenderWindow)
	GetRenderWindow() XBaseRenderWindow
	GetSigLeftButtonPress() *gsignal.Event[MouseButtonEvent]
	GetSigLeftButtonRelease() *gsignal.Event[MouseButtonEvent]
	GetSigLeftButtonDoublePress() *gsignal.Event[MouseButtonEvent]

	GetSigMiddleButtonPress() *gsignal.Event[MouseButtonEvent]
	GetSigMiddleButtonRelease() *gsignal.Event[MouseButtonEvent]
	GetSigMiddleButtonDoublePress() *gsignal.Event[MouseButtonEvent]

	GetSigRightButtonPress() *gsignal.Event[MouseButtonEvent]
	GetSigRightButtonRelease() *gsignal.Event[MouseButtonEvent]
	GetSigRightButtonDoublePress() *gsignal.Event[MouseButtonEvent]

	GetSigResize() *gsignal.Event[xutils.Vec2i]
	GetSigKeyPress() *gsignal.Event[KeyEvent]
	GetSigKeyRelease() *gsignal.Event[KeyEvent]
	GetSigMouseMove() *gsignal.Event[KeyEvent]
	GetSigMouseWheelForward() *gsignal.Event[KeyEvent]
	GetSigMouseWheelBackward() *gsignal.Event[KeyEvent]
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

	SigKeyPress           gsignal.Event[KeyEvent] //键盘摁下
	SigKeyRelease         gsignal.Event[KeyEvent] //键盘弹起
	SigMouseMove          gsignal.Event[KeyEvent] //鼠标移动
	SigMouseWheelForward  gsignal.Event[KeyEvent] //鼠标前进
	SigMouseWheelBackward gsignal.Event[KeyEvent] //鼠标后退
}
