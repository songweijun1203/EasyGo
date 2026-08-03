package xrender

import (
	xutils "EasyGo/internal/utils"
	"EasyGo/internal/xbaserender"
	"EasyGo/internal/xsignal"
)

type XRender struct {
	renderWindow          xbaserender.XBaseRenderWindow
	multiModeEventHandler xbaserender.XBaseMultiModeRenderInteractionEventHandler
	connector             *xsignal.XConnector //信号管理
}

func NewXRender() (*XRender, error) {
	render := &XRender{
		connector: xsignal.NewXConnector(),
	}
	render.init()
	return render, nil
}
func (r *XRender) init() {
	r.getOrCreateMultiModeEventHandler()
}

func (r *XRender) SetRenderWindow(w xbaserender.XBaseRenderWindow) {
	r.renderWindow = w
}

func (r *XRender) GetRenderWindow() xbaserender.XBaseRenderWindow {
	return r.renderWindow
}

func (r *XRender) Render() {

}

func (r *XRender) MakeCurrent() bool {
	return r.GetRenderWindow().MakeCurrent()
}

func (r *XRender) DoneCurrent() {
	r.GetRenderWindow().DoneCurrent()
}

func (r *XRender) getOrCreateMultiModeEventHandler() xbaserender.XBaseMultiModeRenderInteractionEventHandler {
	if r.multiModeEventHandler == nil {
		r.multiModeEventHandler = &XRenderMultiModeInteractionHandler{}
	}
	return r.multiModeEventHandler
}

func (r *XRender) connectToRenderWindowSignals() {
	hanlder := r.getOrCreateMultiModeEventHandler()
	dispatcherSig := r.GetRenderWindow().GetEventDispatcher().GetXRenderWindowEventDispatchSig()

	r.connector.DisconnectAll()

	xsignal.ConnectSlot(r.connector, &dispatcherSig.SigLeftButtonPress, func(eventData xbaserender.MouseButtonEvent) {
		event := xbaserender.XEvent{}
		hanlder.LeftButtonPressEvent(eventData, &event)
	}, "XRenderConnectToRenderWindowSignals_SigLeftButtonPress")

	xsignal.ConnectSlot(r.connector, &dispatcherSig.SigLeftButtonRelease, func(eventData xbaserender.MouseButtonEvent) {
		event := xbaserender.XEvent{}
		hanlder.LeftButtonReleaseEvent(eventData, &event)
	}, "XRenderConnectToRenderWindowSignals_SigLeftButtonRelease")

	xsignal.ConnectSlot(r.connector, &dispatcherSig.SigLeftButtonDoublePress, func(eventData xbaserender.MouseButtonEvent) {
		event := xbaserender.XEvent{}
		hanlder.LeftButtonDoublePressEvent(eventData, &event)
	}, "XRenderConnectToRenderWindowSignals_SigLeftButtonDoublePress")

	xsignal.ConnectSlot(r.connector, &dispatcherSig.SigMiddleButtonPress, func(eventData xbaserender.MouseButtonEvent) {
		event := xbaserender.XEvent{}
		hanlder.MiddleButtonPressEvent(eventData, &event)
	}, "XRenderConnectToRenderWindowSignals_SigMiddleButtonPress")

	xsignal.ConnectSlot(r.connector, &dispatcherSig.SigMiddleButtonRelease, func(eventData xbaserender.MouseButtonEvent) {
		event := xbaserender.XEvent{}
		hanlder.MiddleButtonReleaseEvent(eventData, &event)
	}, "XRenderConnectToRenderWindowSignals_SigMiddleButtonRelease")

	xsignal.ConnectSlot(r.connector, &dispatcherSig.SigMiddleButtonDoublePress, func(eventData xbaserender.MouseButtonEvent) {
		event := xbaserender.XEvent{}
		hanlder.MiddleButtonDoublePressEvent(eventData, &event)
	}, "XRenderConnectToRenderWindowSignals_SigMiddleButtonDoublePress")

	xsignal.ConnectSlot(r.connector, &dispatcherSig.SigRightButtonPress, func(eventData xbaserender.MouseButtonEvent) {
		event := xbaserender.XEvent{}
		hanlder.RightButtonPressEvent(eventData, &event)
	}, "XRenderConnectToRenderWindowSignals_SigRightButtonPress")

	xsignal.ConnectSlot(r.connector, &dispatcherSig.SigRightButtonRelease, func(eventData xbaserender.MouseButtonEvent) {
		event := xbaserender.XEvent{}
		hanlder.RightButtonReleaseEvent(eventData, &event)
	}, "XRenderConnectToRenderWindowSignals_SigRightButtonRelease")

	xsignal.ConnectSlot(r.connector, &dispatcherSig.SigRightButtonDoublePress, func(eventData xbaserender.MouseButtonEvent) {
		event := xbaserender.XEvent{}
		hanlder.RightButtonDoublePressEvent(eventData, &event)
	}, "XRenderConnectToRenderWindowSignals_SigRightButtonDoublePress")

	xsignal.ConnectSlot(r.connector, &dispatcherSig.SigResize, func(pos xutils.Vec2i) {
		event := xbaserender.XEvent{}
		hanlder.ResizeEvent(pos, &event)
	}, "XRenderConnectToRenderWindowSignals_SigResize")

	xsignal.ConnectSlot(r.connector, &dispatcherSig.SigKeyPress, func(eventData xbaserender.KeyEvent) {
		event := xbaserender.XEvent{}
		hanlder.KeyPressEvent(eventData, &event)
	}, "XRenderConnectToRenderWindowSignals_SigKeyPress")

	xsignal.ConnectSlot(r.connector, &dispatcherSig.SigKeyRelease, func(eventData xbaserender.KeyEvent) {
		event := xbaserender.XEvent{}
		hanlder.KeyReleaseEvent(eventData, &event)
	}, "XRenderConnectToRenderWindowSignals_SigKeyRelease")

	xsignal.ConnectSlot(r.connector, &dispatcherSig.SigMouseMove, func(eventData xbaserender.MouseButtonEvent) {
		event := xbaserender.XEvent{}
		hanlder.MouseMoveEvent(eventData, &event)
	}, "XRenderConnectToRenderWindowSignals_SigMouseMove")

	xsignal.ConnectSlot(r.connector, &dispatcherSig.SigMouseWheelForward, func(eventData xbaserender.MouseButtonEvent) {
		event := xbaserender.XEvent{}
		hanlder.MouseWheelForwardEvent(eventData, &event)
	}, "XRenderConnectToRenderWindowSignals_SigMouseWheelForward")

	xsignal.ConnectSlot(r.connector, &dispatcherSig.SigMouseWheelBackward, func(eventData xbaserender.MouseButtonEvent) {
		event := xbaserender.XEvent{}
		hanlder.MouseWheelBackwardEvent(eventData, &event)
	}, "XRenderConnectToRenderWindowSignals_SigMouseWheelBackward")
}

func (r *XRender) Close() {
	//析构的清理工作
	r.connector.DisconnectAll()
}
