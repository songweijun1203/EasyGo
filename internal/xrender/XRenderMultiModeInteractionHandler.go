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
