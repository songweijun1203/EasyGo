package xopengl

import (
	"fmt"
	"testing"
)

// 离屏渲染环境创建测试
func TestOffsetWindow(t *testing.T) {
	w, _ := NewXOffsetWindow()

	ctx := NewXOpenGLContext()
	ctx.Create(w.WinId())

	if ctx.MakeCurrent() {
		fmt.Println("OPenGL context created success!")
		{
			//此处进行GL相关的测试
		}
		ctx.DoneCurrent()
	}

	defer w.Close()
	defer ctx.Close() //先执行
}
