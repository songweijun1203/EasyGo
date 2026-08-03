package xglwidget

import (
	"EasyGo/internal/xbaserender"

	"github.com/go-gl/glfw/v3.4/glfw"
)

type GlfwKeyData struct {
	Key      glfw.Key
	Scancode int //状态 摁下还是抬起
}

func convertMods(m glfw.ModifierKey) xbaserender.KeyboardModifier {
	switch m {
	case glfw.ModShift:
		return xbaserender.ShiftModifier
	case glfw.ModControl:
		return xbaserender.ControlModifier
	case glfw.ModAlt:
		return xbaserender.AltModifier
	case glfw.ModSuper:
		return xbaserender.MetaModifier
	default:
		return xbaserender.NoModifier
	}
}

func convertKey(key glfw.Key) xbaserender.KeyType {
	switch key {
	case glfw.KeyA:
		return xbaserender.Key_A
	case glfw.KeyB:
		return xbaserender.Key_B
	case glfw.KeyC:
		return xbaserender.Key_C
	case glfw.KeyD:
		return xbaserender.Key_D
	case glfw.KeyE:
		return xbaserender.Key_E
	case glfw.KeyF:
		return xbaserender.Key_F
	default:
		return xbaserender.Key_Unknown
	}
}
