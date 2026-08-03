package xrender

import (
	xutils "EasyGo/internal/utils"
	"EasyGo/internal/xbaserender"
)

type XInteractionEventHandler struct {
	xbaserender.XRenderWindowEventDispatchSig
	render *xbaserender.XBaseRender
}

func (i *XInteractionEventHandler) LeftButtonPressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//处理事件，转发信号
	i.SigLeftButtonPress.Emit(eventData)
}

func (i *XInteractionEventHandler) LeftButtonReleaseEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//处理事件，转发信号
	i.SigLeftButtonRelease.Emit(eventData)
}

func (i *XInteractionEventHandler) LeftButtonDoublePressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//处理事件，转发信号
	i.SigLeftButtonDoublePress.Emit(eventData)
}

func (i *XInteractionEventHandler) RightButtonPressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//处理事件，转发信号
	i.SigRightButtonPress.Emit(eventData)
}

func (i *XInteractionEventHandler) RightButtonReleaseEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//处理事件，转发信号
	i.SigRightButtonRelease.Emit(eventData)
}

func (i *XInteractionEventHandler) RightButtonDoublePressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//处理事件，转发信号
	i.SigRightButtonDoublePress.Emit(eventData)
}

func (i *XInteractionEventHandler) MiddleButtonPressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//处理事件，转发信号
	i.SigMiddleButtonPress.Emit(eventData)
}

func (i *XInteractionEventHandler) MiddleButtonReleaseEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//处理事件，转发信号
	i.SigMiddleButtonRelease.Emit(eventData)
}

func (i *XInteractionEventHandler) MiddleButtonDoublePressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//处理事件，转发信号
	i.SigMiddleButtonDoublePress.Emit(eventData)
}

func (i *XInteractionEventHandler) ResizeEvent(pos xutils.Vec2i, e *xbaserender.XEvent) {
	//处理事件，转发信号
	i.SigResize.Emit(pos)
}

func (i *XInteractionEventHandler) KeyPressEvent(eventData xbaserender.KeyEvent, e *xbaserender.XEvent) {
	//处理事件，转发信号
	i.SigKeyPress.Emit(eventData)
}

func (i *XInteractionEventHandler) KeyReleaseEvent(eventData xbaserender.KeyEvent, e *xbaserender.XEvent) {
	//处理事件，转发信号
	i.SigKeyRelease.Emit(eventData)
}

func (i *XInteractionEventHandler) MouseMoveEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//处理事件，转发信号
	i.SigMouseMove.Emit(eventData)
}

func (i *XInteractionEventHandler) MouseWheelForwardEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//处理事件，转发信号
	i.SigMouseWheelForward.Emit(eventData)
}

func (i *XInteractionEventHandler) MouseWheelBackwardEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//处理事件，转发信号
	i.SigMouseWheelBackward.Emit(eventData)
}

func (i *XInteractionEventHandler) SetRender(render *xbaserender.XBaseRender) {
	i.render = render
}

func (i *XInteractionEventHandler) GetRender() *xbaserender.XBaseRender {
	return i.render
}
