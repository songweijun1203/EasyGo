package xopengl

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type XOpenGLEnable struct {
	depthTestState        StateDepthTest
	multisampleState      StateMultiSample
	scissorTestState      StateScissorTest
	blendState            StateBlendTest
	cullFaceState         StateCullFace
	programPointSizeState StateProgramPointSize
	stencilTestState      StateStencilTest
	polygonOffsetState    StatePolyGonOffsetFill
}

func NewXOpenGLEnable() *XOpenGLEnable {
	e := &XOpenGLEnable{}
	return e
}

func (e *XOpenGLEnable) Enable(t GLEnableTestType) error {
	gl.Enable(uint32(t))
	return CheckGLerror()
}

func (e *XOpenGLEnable) Save() {
	e.saveDepthTestState()
	e.saveMultisampleState()
	e.saveScissorTestState()
	e.saveBlendState()
	e.saveCullFaceState()
	e.saveProgramPointSizeState()
	e.saveStencilTestState()
	e.savePolygenOffsetState()
}

func (e *XOpenGLEnable) Restore() {
	e.restoreDepthTestState()
	e.restoreMultisampleState()
	e.restoreScissorTestState()
	e.restoreBlendState()
	e.restoreCullFaceState()
	e.restoreStencilTestState()
	e.restorePolygenOffsetState()
}

// 保存
func (e *XOpenGLEnable) saveDepthTestState() {
	e.depthTestState.isEnable = gl.IsEnabled(uint32(DEPTH_TEST))
	e.depthTestState.isDepthClampEnable = gl.IsEnabled(uint32(DEPTH_CLAMP))
}

func (e *XOpenGLEnable) saveMultisampleState() {
	e.multisampleState.isEnable = gl.IsEnabled(uint32(MULTISAMPLE))
}

func (e *XOpenGLEnable) saveScissorTestState() {
	e.scissorTestState.isEnable = gl.IsEnabled(uint32(SCISSOR_TEST))
	e.scissorTestState.rect, _ = XglGetScissor()
}

func (e *XOpenGLEnable) saveBlendState() {
	e.blendState.isEnable = gl.IsEnabled(uint32(BLEND))
	e.blendState.srcFactor, e.blendState.dstFactor, _ = XgetBlendFunc()
}

func (e *XOpenGLEnable) saveCullFaceState() {
	e.cullFaceState.isEnable = gl.IsEnabled(uint32(CULL_FACE))
}

func (e *XOpenGLEnable) saveProgramPointSizeState() {}

func (e *XOpenGLEnable) saveStencilTestState() {
	e.stencilTestState.isEnable = gl.IsEnabled(uint32(STENCIL_TEST))
	e.stencilTestState.compFun,
		e.stencilTestState.ref,
		e.stencilTestState.mask,
		_ = XglGetStencilFunc()

	e.stencilTestState.sfail,
		e.stencilTestState.zfail,
		e.stencilTestState.zpass,
		_ = XGetStencilOp()
}

func (e *XOpenGLEnable) savePolygenOffsetState() {
	e.polygonOffsetState.isEnable = gl.IsEnabled(uint32(POLYGON_OFFSET_FILL))
	e.polygonOffsetState.factor,
		e.polygonOffsetState.unit,
		_ = XGetPolygonOffsetParm()
}

// 恢复
func (e *XOpenGLEnable) restoreDepthTestState() {
	if e.depthTestState.isEnable {
		gl.Enable(uint32(DEPTH_TEST))
	} else {
		gl.Disable(uint32(DEPTH_TEST))
	}

	if e.depthTestState.isDepthClampEnable {
		gl.Enable(uint32(DEPTH_CLAMP))
	} else {
		gl.Disable(uint32(DEPTH_CLAMP))
	}
}

func (e *XOpenGLEnable) restoreMultisampleState() {
	if e.multisampleState.isEnable {
		gl.Enable(uint32(MULTISAMPLE))
	} else {
		gl.Disable(uint32(MULTISAMPLE))
	}
}

func (e *XOpenGLEnable) restoreScissorTestState() {
	if e.scissorTestState.isEnable {
		gl.Enable(uint32(SCISSOR_TEST))
	} else {
		gl.Disable(uint32(SCISSOR_TEST))
	}
}

func (e *XOpenGLEnable) restoreBlendState() {
	if e.blendState.isEnable {
		gl.Enable(uint32(BLEND))
	} else {
		gl.Disable(uint32(BLEND))
	}
	XglBlendFunc(e.blendState.srcFactor, e.blendState.dstFactor)
}

func (e *XOpenGLEnable) restoreCullFaceState() {
	if e.cullFaceState.isEnable {
		gl.Enable(uint32(CULL_FACE))
	} else {
		gl.Disable(uint32(CULL_FACE))
	}
}

func (e *XOpenGLEnable) restoreStencilTestState() {
	if e.stencilTestState.isEnable {
		gl.Enable(uint32(STENCIL_TEST))
	} else {
		gl.Disable(uint32(STENCIL_TEST))
	}
	XglSetStencilOp(e.stencilTestState.sfail, e.stencilTestState.zfail, e.stencilTestState.zpass)
}

func (e *XOpenGLEnable) restorePolygenOffsetState() {
	if e.polygonOffsetState.isEnable {
		gl.Enable(uint32(POLYGON_OFFSET_FILL))
	} else {
		gl.Disable(uint32(POLYGON_OFFSET_FILL))
	}
	XglPolygonOffset(e.polygonOffsetState.factor, e.polygonOffsetState.unit)
}
