package xrender

import (
	xutils "EasyGo/internal/utils"
	"EasyGo/internal/xbaserender"
	"EasyGo/internal/xopengl"
	"EasyGo/internal/xsignal"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type XRender struct {
	renderWindow          xbaserender.XBaseRenderWindow
	multiModeEventHandler xbaserender.XBaseMultiModeRenderInteractionEventHandler
	connector             *xsignal.XConnector //信号管理
	viewPort              xutils.Rectf        //设置视口大小 0-1，x y w h,原点位于窗口左下角
}

func NewXRender() (*XRender, error) {
	render := &XRender{
		connector: xsignal.NewXConnector(),
		viewPort:  xutils.Rectf{X: 0, Y: 0, Z: 0.5, W: 1},
	}
	render.init()
	return render, nil
}
func (r *XRender) init() {
	r.getOrCreateMultiModeEventHandler()
}

func (r *XRender) SetViewPort(x float32, y float32, w float32, h float32) {
	r.viewPort = xutils.Rectf{X: x, Y: y, Z: w, W: h}
}
func (r *XRender) GetViewPort() (x float32, y float32, w float32, h float32) {
	return r.viewPort.X, r.viewPort.Y, r.viewPort.Z, r.viewPort.W
}

// 获取窗口坐标系下转换后的视口
func (r *XRender) GetConvertViewPort() (x, y, z, w int32) {
	renderWindow := r.GetRenderWindow()
	width := renderWindow.GetWindowWidth()
	height := renderWindow.GetWindowHeigth()

	nx, ny, nw, nh := r.GetViewPort()
	x1 := int32(nx * float32(width))
	y1 := int32(ny * float32(height))
	w1 := int32(nw * float32(width))
	h1 := int32(nh * float32(height))
	return x1, y1, w1, h1
}

func (r *XRender) SetRenderWindow(w xbaserender.XBaseRenderWindow) {
	r.renderWindow = w
}

func (r *XRender) GetRenderWindow() xbaserender.XBaseRenderWindow {
	return r.renderWindow
}

func (r *XRender) Render() {
	r.MakeCurrent()
	gl.Enable(gl.SCISSOR_TEST)
	xopengl.XglViewPort(r.GetConvertViewPort())
	xopengl.XglScissor(r.GetConvertViewPort())
	gl.ClearColor(0.1, 0.1, 0.2, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	r.DoneCurrent()
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

	r.disConnectSignalSlot() //关系信号槽连接

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

	//对象内信号槽绑定
}

func (r *XRender) disConnectSignalSlot() {
	r.connector.DisconnectAll()
}

func (r *XRender) Close() {
	//析构的清理工作
	r.disConnectSignalSlot()

	//gl相关资源关闭
}
