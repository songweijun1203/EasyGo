package main

import (
	"EasyGo/internal/xglwidget"
	"EasyGo/internal/xrender"
	"fmt" // 替换为 fmt，用于格式化错误
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	widget, err := xglwidget.NewXGLWidget()
	if err != nil {
		log.Fatalln("failed to initialize glfw:", err)
		return
	}

	renderWindow := widget.GetRenderWindow()
	render, _ := xrender.NewXRender()
	renderWindow.AddRender(render)

	//render添加渲染节点

	widget.Run()

	/*
		if err := glfw.Init(); err != nil {
			log.Fatalln("failed to initialize glfw:", err)
		}
		defer glfw.Terminate()

		//glfw.WindowHint(glfw.ContextVersionMajor, 4)
		//glfw.WindowHint(glfw.ContextVersionMinor, 1)
		//glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		//glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
		glfw.WindowHint(glfw.Resizable, glfw.True)

		window, err := glfw.CreateWindow(800, 600, "Go Render Engine - Triangle", nil, nil)
		if err != nil {
			log.Fatalln("failed to create window:", err)
		}

		defer window.Destroy()
		//window.MakeContextCurrent()

		// 获取窗口的 HWND（Windows 平台特有）
		hwndRaw := window.GetWin32Window() // 返回 uintptr
		hwnd := uintptr(unsafe.Pointer(hwndRaw))
		if hwnd == 0 {
			log.Fatalln("failed to get HWND")
		}

		// 创建并激活 XOpenGLContext
		ctx := xopengl.NewXOpenGLContext()
		if !ctx.Create(hwnd) {
			log.Fatalln("failed to create OpenGL context with XOpenGLContext")
		}
		defer ctx.Reset()

		if !ctx.MakeCurrent() {
			log.Fatalln("failed to make context current")
		}
		defer ctx.DoneCurrent()

		if err := gl.Init(); err != nil {
			log.Fatalln("failed to init gl:", err)
		}
		log.Println("OpenGL Version:", gl.GoStr(gl.GetString(gl.VERSION)))

		vertices := []float32{
			-0.5, -0.5, 0.0,
			0.5, -0.5, 0.0,
			0.0, 0.5, 0.0,
		}

		vao := xopengl.NewXOpenGLVertexArrayObjec()
		vao.Create()

		vbo := xopengl.NewXOpenGLBuffer(xopengl.BufferTypeVertexBuffer)
		vbo.SetUsagePattern(xopengl.UsageStaticDraw)
		vbo.Create()
		vbo.Bind()
		vbo.AllocateEmpty(len(vertices) * 4)
		vbo.Write(0, len(vertices)*4, gl.Ptr(vertices))
		vbo.Release()
		//readBuffer, _ := vbo.Read(0, len(vertices)*4)
		//ss := xq_utils.As[float32](readBuffer)
		//fmt.Println(readBuffer.GetInternalData())
		//fmt.Println(ss.GetInternalData())

		vao.AddBuffer(0, vbo, 3, xopengl.GLDataType_float, uint32(3*unsafe.Sizeof(float32(0))), 0, false)
		vao.EnableAttribute(0)

		vbo.Release()
		vao.Release()

		shaderProgram, err := newShaderProgram(
			`#version 330 core
			layout (location = 0) in vec3 aPos;
			void main() {
				gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
			}`+"\x00",
			`#version 330 core
			out vec4 FragColor;
			void main() {
				FragColor = vec4(1.0, 0.5, 0.2, 1.0);
			}`+"\x00",
		)
		if err != nil {
			log.Fatalln("failed to compile shader:", err)
		}

		for !window.ShouldClose() {
			gl.ClearColor(0.1, 0.1, 0.2, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			gl.UseProgram(shaderProgram)
			//gl.BindVertexArray(VAO)
			vao.Bind()
			gl.DrawArrays(gl.TRIANGLES, 0, 3)

			//window.SwapBuffers()
			ctx.SwapBuffers()
			glfw.PollEvents()
		}
		vbo.Destory()
		vao.Destory()
	*/
}

func newShaderProgram(vertexSrc, fragmentSrc string) (uint32, error) {
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	cVertexSrc, free := gl.Strs(vertexSrc)
	gl.ShaderSource(vertexShader, 1, cVertexSrc, nil)
	free()
	gl.CompileShader(vertexShader)
	if err := checkShaderCompile(vertexShader); err != nil {
		return 0, err
	}

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	cFragmentSrc, free := gl.Strs(fragmentSrc)
	gl.ShaderSource(fragmentShader, 1, cFragmentSrc, nil)
	free()
	gl.CompileShader(fragmentShader)
	if err := checkShaderCompile(fragmentShader); err != nil {
		return 0, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var success int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		logMsg := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(logMsg))
		// 修复：使用 fmt.Errorf 替代不存在的 log.Errorf
		return 0, fmt.Errorf("failed to link shader program: %s", logMsg)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func checkShaderCompile(shader uint32) error {
	var success int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		logMsg := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logMsg))
		// 修复：使用 fmt.Errorf 替代不存在的 log.Errorf
		return fmt.Errorf("shader compilation failed: %s", logMsg)
	}
	return nil
}
