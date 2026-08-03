package xrender

import (
	xutils "EasyGo/internal/utils"
	"EasyGo/internal/xbaserender"
	"fmt"
)

type XCameraNavigationHandler struct {
	XInteractionEventHandler
}

func (i *XCameraNavigationHandler) LeftButtonPressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//处理事件，转发信号
	fmt.Printf("LeftButtonPressEvent\n")
	i.SigLeftButtonPress.Emit(eventData)
}

func (i *XCameraNavigationHandler) LeftButtonReleaseEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	fmt.Printf("LeftButtonReleaseEvent\n")
	i.SigLeftButtonRelease.Emit(eventData)
}

func (i *XCameraNavigationHandler) LeftButtonDoublePressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	fmt.Printf("LeftButtonDoublePressEvent\n")
	i.SigLeftButtonDoublePress.Emit(eventData)
}

func (i *XCameraNavigationHandler) RightButtonPressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	fmt.Printf("RightButtonPressEvent\n")
	i.SigRightButtonPress.Emit(eventData)
}

func (i *XCameraNavigationHandler) RightButtonReleaseEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	fmt.Printf("RightButtonReleaseEvent\n")
	i.SigRightButtonRelease.Emit(eventData)
}

func (i *XCameraNavigationHandler) RightButtonDoublePressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	fmt.Printf("RightButtonDoublePressEvent\n")
	i.SigRightButtonDoublePress.Emit(eventData)
}

func (i *XCameraNavigationHandler) MiddleButtonPressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	fmt.Printf("MiddleButtonPressEvent\n")
	i.SigMiddleButtonPress.Emit(eventData)
}

func (i *XCameraNavigationHandler) MiddleButtonReleaseEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	fmt.Printf("MiddleButtonReleaseEvent\n")
	i.SigMiddleButtonRelease.Emit(eventData)
}

func (i *XCameraNavigationHandler) MiddleButtonDoublePressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	fmt.Printf("MiddleButtonDoublePressEvent\n")
	i.SigMiddleButtonDoublePress.Emit(eventData)
}

func (i *XCameraNavigationHandler) ResizeEvent(pos xutils.Vec2i, e *xbaserender.XEvent) {
	fmt.Printf("ResizeEvent\n")
}

func (i *XCameraNavigationHandler) KeyPressEvent(eventData xbaserender.KeyEvent, e *xbaserender.XEvent) {
	fmt.Printf("KeyPressEvent\n")
	i.SigKeyPress.Emit(eventData)
}

func (i *XCameraNavigationHandler) KeyReleaseEvent(eventData xbaserender.KeyEvent, e *xbaserender.XEvent) {
	fmt.Printf("KeyReleaseEvent\n")
	i.SigKeyRelease.Emit(eventData)
}

func (i *XCameraNavigationHandler) MouseMoveEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	fmt.Printf("MouseMoveEvent\n")
	i.SigMouseMove.Emit(eventData)
}

func (i *XCameraNavigationHandler) MouseWheelForwardEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	fmt.Printf("MouseWheelForwardEvent\n")
}

func (i *XCameraNavigationHandler) MouseWheelBackwardEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	fmt.Printf("MouseWheelBackwardEvent\n")
}
