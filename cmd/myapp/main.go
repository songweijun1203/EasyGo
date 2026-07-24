package main

import (
	"fmt" // 替换为 fmt，用于格式化错误
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.4/glfw"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "Go Render Engine - Triangle", nil, nil)
	if err != nil {
		log.Fatalln("failed to create window:", err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to init gl:", err)
	}
	log.Println("OpenGL Version:", gl.GoStr(gl.GetString(gl.VERSION)))

	vertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}

	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

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
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		window.SwapBuffers()
		glfw.PollEvents()
	}
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
