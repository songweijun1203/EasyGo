package xbaserender

import xutils "EasyGo/internal/utils"

type XBaseInteractionEventHandler interface {
	LeftButtonPressEvent(eventData MouseButtonEvent, e *XEvent)
	LeftButtonReleaseEvent(eventData MouseButtonEvent, e *XEvent)
	LeftButtonDoublePressEvent(eventData MouseButtonEvent, e *XEvent)

	RightButtonPressEvent(eventData MouseButtonEvent, e *XEvent)
	RightButtonReleaseEvent(eventData MouseButtonEvent, e *XEvent)
	RightButtonDoublePressEvent(eventData MouseButtonEvent, e *XEvent)

	MiddleButtonPressEvent(eventData MouseButtonEvent, e *XEvent)
	MiddleButtonReleaseEvent(eventData MouseButtonEvent, e *XEvent)
	MiddleButtonDoublePressEvent(eventData MouseButtonEvent, e *XEvent)

	ResizeEvent(pos xutils.Vec2i, e *XEvent)
	KeyPressEvent(eventData KeyEvent, e *XEvent)
	KeyReleaseEvent(eventData KeyEvent, e *XEvent)
	MouseMoveEvent(eventData MouseButtonEvent, e *XEvent)
	MouseWheelForwardEvent(eventData MouseButtonEvent, e *XEvent)
	MouseWheelBackwardEvent(eventData MouseButtonEvent, e *XEvent)
}

type XBaseRenderInteractionEventHandler interface {
	XBaseInteractionEventHandler
	SetRender(*XBaseRender)
	GetRender() *XBaseRender
}

type XBaseMultiModeRenderInteractionEventHandler interface {
	XBaseRenderInteractionEventHandler
	SetCameraNavigationHandler(XBaseRenderInteractionEventHandler)
	SetPickHandler(XBaseRenderInteractionEventHandler)
	SetManipulatorHandler(XBaseRenderInteractionEventHandler)
	SetDrawHandler(XBaseRenderInteractionEventHandler)
	GetCameraNavigationHandler() XBaseRenderInteractionEventHandler
	GetPickHandler() XBaseRenderInteractionEventHandler
	GetManipulatorHandler() XBaseRenderInteractionEventHandler
	GetDrawHandler() XBaseRenderInteractionEventHandler
}
