package xrender

import "EasyGo/internal/xbaserender"

type XRenderMultiModeInteractionHandler struct {
	XInteractionEventHandler //自身的实现
	cameraNavigationHandler  xbaserender.XBaseRenderInteractionEventHandler
	pickHandler              xbaserender.XBaseRenderInteractionEventHandler
	manipulatorHandler       xbaserender.XBaseRenderInteractionEventHandler
	drawHandler              xbaserender.XBaseRenderInteractionEventHandler
}

func (h *XRenderMultiModeInteractionHandler) SetCameraNavigationHandler(handler xbaserender.XBaseRenderInteractionEventHandler) {
	// Implementation for setting camera navigation handler
	h.cameraNavigationHandler = handler
}

func (h *XRenderMultiModeInteractionHandler) SetPickHandler(handler xbaserender.XBaseRenderInteractionEventHandler) {
	// Implementation for setting pick handler
	h.pickHandler = handler
}

func (h *XRenderMultiModeInteractionHandler) SetManipulatorHandler(handler xbaserender.XBaseRenderInteractionEventHandler) {
	// Implementation for setting manipulator handler
	h.manipulatorHandler = handler
}

func (h *XRenderMultiModeInteractionHandler) SetDrawHandler(handler xbaserender.XBaseRenderInteractionEventHandler) {
	// Implementation for setting draw handler
	h.drawHandler = handler
}

func (h *XRenderMultiModeInteractionHandler) GetCameraNavigationHandler() xbaserender.XBaseRenderInteractionEventHandler {
	// Implementation for getting camera navigation handler
	if h.cameraNavigationHandler == nil {
		h.cameraNavigationHandler = &XCameraNavigationHandler{}
	}
	return h.cameraNavigationHandler
}

func (h *XRenderMultiModeInteractionHandler) GetPickHandler() xbaserender.XBaseRenderInteractionEventHandler {
	// Implementation for getting pick handler
	return h.pickHandler
}

func (h *XRenderMultiModeInteractionHandler) GetManipulatorHandler() xbaserender.XBaseRenderInteractionEventHandler {
	// Implementation for getting manipulator handler
	return h.manipulatorHandler
}

func (h *XRenderMultiModeInteractionHandler) GetDrawHandler() xbaserender.XBaseRenderInteractionEventHandler {
	// Implementation for getting draw handler
	return h.drawHandler
}

//事件转发
func (h *XRenderMultiModeInteractionHandler) LeftButtonPressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//根据当前模式选择处理器
	hanlder := h.GetCameraNavigationHandler()
	if hanlder != nil {
		hanlder.LeftButtonPressEvent(eventData, e)
	}
}

func (h *XRenderMultiModeInteractionHandler) LeftButtonReleaseEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//根据当前模式选择处理器
	hanlder := h.GetCameraNavigationHandler()
	if hanlder != nil {
		hanlder.LeftButtonReleaseEvent(eventData, e)
	}
}

func (h *XRenderMultiModeInteractionHandler) LeftButtonDoublePressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//根据当前模式选择处理器
	hanlder := h.GetCameraNavigationHandler()
	if hanlder != nil {
		hanlder.LeftButtonDoublePressEvent(eventData, e)
	}
}

func (h *XRenderMultiModeInteractionHandler) RightButtonPressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//根据当前模式选择处理器
	hanlder := h.GetCameraNavigationHandler()
	if hanlder != nil {
		hanlder.RightButtonPressEvent(eventData, e)
	}
}

func (h *XRenderMultiModeInteractionHandler) RightButtonReleaseEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//根据当前模式选择处理器
	hanlder := h.GetCameraNavigationHandler()
	if hanlder != nil {
		hanlder.RightButtonReleaseEvent(eventData, e)
	}
}

func (h *XRenderMultiModeInteractionHandler) RightButtonDoublePressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//根据当前模式选择处理器
	hanlder := h.GetCameraNavigationHandler()
	if hanlder != nil {
		hanlder.RightButtonDoublePressEvent(eventData, e)
	}
}

func (h *XRenderMultiModeInteractionHandler) MiddleButtonPressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//根据当前模式选择处理器
	hanlder := h.GetCameraNavigationHandler()
	if hanlder != nil {
		hanlder.MiddleButtonPressEvent(eventData, e)
	}
}

func (h *XRenderMultiModeInteractionHandler) MiddleButtonReleaseEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//根据当前模式选择处理器
	hanlder := h.GetCameraNavigationHandler()
	if hanlder != nil {
		hanlder.MiddleButtonReleaseEvent(eventData, e)
	}
}

func (h *XRenderMultiModeInteractionHandler) MiddleButtonDoublePressEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//根据当前模式选择处理器
	hanlder := h.GetCameraNavigationHandler()
	if hanlder != nil {
		hanlder.MiddleButtonDoublePressEvent(eventData, e)
	}
}

func (h *XRenderMultiModeInteractionHandler) MouseMoveEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//根据当前模式选择处理器
	hanlder := h.GetCameraNavigationHandler()
	if hanlder != nil {
		hanlder.MouseMoveEvent(eventData, e)
	}
}

func (h *XRenderMultiModeInteractionHandler) MouseWheelForwardEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//根据当前模式选择处理器
	hanlder := h.GetCameraNavigationHandler()
	if hanlder != nil {
		hanlder.MouseWheelForwardEvent(eventData, e)
	}
}

func (h *XRenderMultiModeInteractionHandler) MouseWheelBackwardEvent(eventData xbaserender.MouseButtonEvent, e *xbaserender.XEvent) {
	//根据当前模式选择处理器
	hanlder := h.GetCameraNavigationHandler()
	if hanlder != nil {
		hanlder.MouseWheelBackwardEvent(eventData, e)
	}
}

func (h *XRenderMultiModeInteractionHandler) KeyPressEvent(eventData xbaserender.KeyEvent, e *xbaserender.XEvent) {
	//根据当前模式选择处理器
	hanlder := h.GetCameraNavigationHandler()
	if hanlder != nil {
		hanlder.KeyPressEvent(eventData, e)
	}
}

func (h *XRenderMultiModeInteractionHandler) KeyReleaseEvent(eventData xbaserender.KeyEvent, e *xbaserender.XEvent) {
	//根据当前模式选择处理器
	hanlder := h.GetCameraNavigationHandler()
	if hanlder != nil {
		hanlder.KeyReleaseEvent(eventData, e)
	}
}
