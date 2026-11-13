package xshader

import (
	"EasyGo/internal/xopengl"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

func TestShader(t *testing.T) {
	GetConfigDir := func() string {
		_, filename, _, _ := runtime.Caller(0) // 当前测试文件的路径
		dir := filepath.Dir(filename)
		dir = filepath.Dir(dir) //上一级目录
		dir = filepath.Dir(dir) //上一级目录

		dir = filepath.Join(dir, "configs")
		dir = filepath.Join(dir, "glsl")
		return dir
	}
	dir := GetConfigDir()
	println(dir)

	w, _ := xopengl.NewXOffsetWindow()

	ctx := xopengl.NewXOpenGLContext()
	ctx.Create(w.WinId())

	if ctx.MakeCurrent() {
		fmt.Println("OPenGL context created success!")
		{
			//此处进行GL相关的测试
			vsPath := filepath.Join(dir, "tri.vs")
			fsPath := filepath.Join(dir, "tri.fs")

			shader := NewXBaseShader()
			shader.Create()

			id1, _ := shader.CompileFromFile(VERTEX, vsPath)
			id2, _ := shader.CompileFromFile(FRAGMENT, fsPath)

			shader.Link([]uint32{id1, id2})
			shader.Close()
		}
		ctx.DoneCurrent()
	}

	defer w.Close()
	defer ctx.Close() //先执行
}
