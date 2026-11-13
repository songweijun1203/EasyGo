package xrendernode

import (
	"EasyGo/internal/xbaserendernode"
	"EasyGo/internal/xutils"
)

type XRenderNodeSource struct {
	xbaserendernode.XBaseRenderNodeSource
	//数据
	vertexCoord *xutils.XFloat32Array1D
	FaceIndexs  *xutils.XUInt32Array1D
}
