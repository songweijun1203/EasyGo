package xbaserendernode

import (
	"EasyGo/internal/xbaserender"
	"EasyGo/internal/xutils"
)

// 渲染节点数据源接口
type XBaseRenderNodeSource interface {
	Update()
	UpdateVertextCoordArray() //更新vbo，用于将CPU端顶点法线数据更新后，将顶点法线信息更新到底层vbo的buffer中
	UpdateFaceIndexArray()    //更新ebo，用于将CPU端三角形索引据更新，同时修改时间戳

	GetVertextCoordArray() xutils.XFloat32Array1D
	GetFaceIndexArray() xutils.XUInt32Array1D
}

type XBaseRenderNode interface {
	Draw(render *xbaserender.XBaseRender, parentMatrix xutils.Mat4f)
	SetNodeSource(source XBaseRenderNodeSource)
	GetNodeSource() XBaseRenderNodeSource
}
