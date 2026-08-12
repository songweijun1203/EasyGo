package xrendernode

import (
	"EasyGo/internal/xbaserendernode"
	"EasyGo/internal/xutils"
)

type XRenderNode struct {
	xbaserendernode.XBaseRenderNode
	transform xutils.XAffine3f
	id        uint32
	//纹理数据 TODO
}
