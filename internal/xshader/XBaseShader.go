package xshader

import "C"
import (
	"EasyGo/internal/xopengl"
	"EasyGo/internal/xutils"
	"fmt"
	"log"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type ShaderType uint32

const (
	VERTEX          ShaderType = gl.VERTEX_SHADER
	FRAGMENT        ShaderType = gl.FRAGMENT_SHADER
	GEOMETRY        ShaderType = gl.GEOMETRY_SHADER
	TESS_CONTROL    ShaderType = gl.TESS_CONTROL_SHADER
	TESS_EVALUATION ShaderType = gl.TESS_EVALUATION_SHADER
	COMPUTE         ShaderType = gl.COMPUTE_SHADER
	PROGRAME        ShaderType = gl.PROGRAM
)

type XShader struct {
	id uint32
}

func NewXBaseShader() *XShader {
	return &XShader{}
}

func (shader *XShader) Create() error {
	shader.id = gl.CreateProgram()
	return xopengl.CheckGLerror()
}

func (shader *XShader) CompileFromCode(t ShaderType, shaderCode string) (id uint32, err error) {
	//编译代码
	shaderCode = shaderCode + "\x00"
	fragmentShader := gl.CreateShader(uint32(t))

	cFragmentSrc, free := gl.Strs(shaderCode)
	defer free()
	//fmt.Printf("Shader source: %q\n", shaderCode)

	gl.ShaderSource(fragmentShader, 1, cFragmentSrc, nil)

	gl.CompileShader(fragmentShader)

	return fragmentShader, checkCompileErrors(fragmentShader, t)
}

func (shader *XShader) CompileFromFile(t ShaderType, shaderFile string) (id uint32, err error) {
	//读取文件内容
	data, err := os.ReadFile(shaderFile)
	if err != nil {
		log.Fatal(err)
	}

	str := string(data)
	return shader.CompileFromCode(t, str)
}

func (shader *XShader) Link(ids []uint32) error {
	for _, id := range ids {
		gl.AttachShader(shader.id, id)
	}

	gl.LinkProgram(shader.id)

	//检查编译问题
	checkCompileErrors(shader.id, PROGRAME)
	for _, id := range ids {
		gl.DeleteShader(id)
	}
	return xopengl.CheckGLerror()
}

func (shader *XShader) Use() error {
	gl.UseProgram(shader.id)
	return xopengl.CheckGLerror()
}

func (shader *XShader) UnUse() error {
	gl.UseProgram(0)
	return xopengl.CheckGLerror()
}

func (shader *XShader) Close() error {
	gl.DeleteShader(shader.id)
	return xopengl.CheckGLerror()
}

func (shader *XShader) SetBool(name string, value bool) error {
	val := 0
	if value {
		val = 1
	}

	gl.Str(name)
	gl.Uniform1i(gl.GetUniformLocation(shader.id, gl.Str(name)), int32(val))
	return xopengl.CheckGLerror()
}

func (shader *XShader) SetInt(name string, value int32) error {
	return nil
}

func (shader *XShader) SetUInt(name string, value uint32) error {
	return nil
}

func (shader *XShader) SetFloat(name string, value float32) error {
	return nil
}

func (shader *XShader) SetVec2(name string, x, y float32) error {
	return nil
}

func (shader *XShader) SetVec3(name string, x, y, z float32) error {
	return nil
}

func (shader *XShader) SetVec4(name string, x, y, z, w float32) error {
	return nil
}

func (shader *XShader) SetMat2(name string, m xutils.Mat2f) error {
	return nil
}

func (shader *XShader) SetMat3(name string, m xutils.Mat3f) error {
	return nil
}

func (shader *XShader) SetMat4(name string, m xutils.Mat4f) error {
	return nil
}

// --------------------------------------------------------------------
// 检查编译错误信息
func checkCompileErrors(id uint32, t ShaderType) error {
	mapType2String := make(map[ShaderType]string)
	mapType2String[VERTEX] = "VERTEX"
	mapType2String[FRAGMENT] = "FRAGMENT"
	mapType2String[GEOMETRY] = "GEOMETRY"
	mapType2String[TESS_CONTROL] = "TESS_CONTROL"
	mapType2String[TESS_EVALUATION] = "TESS_EVALUATION"
	mapType2String[COMPUTE] = "COMPUTE"
	mapType2String[PROGRAME] = "PROGRAME"

	var success int32
	var infoLog [1024]uint8
	var logLength int32
	if t != PROGRAME {
		gl.GetShaderiv(id, gl.COMPILE_STATUS, &success)
		if success == 0 {
			gl.GetShaderInfoLog(id, 1024, &logLength, &infoLog[0])
			logMsg := string(infoLog[:logLength])
			fmt.Printf("ERROR::SHADER_COMPILATION_ERROR of type: %s\n ", mapType2String[t])
			fmt.Println(logMsg)
			fmt.Println("-- --------------------------------------------------- -- ")
			return fmt.Errorf("logMsg")
		}
	} else {
		gl.GetProgramiv(id, gl.LINK_STATUS, &success)
		if success == 0 {
			gl.GetProgramInfoLog(id, 1024, &logLength, &infoLog[0])
			fmt.Printf("ERROR::PROGRAM_LINKING_ERROR of type: %s\n ", mapType2String[t])
			logMsg := string(infoLog[:logLength])
			fmt.Println(logMsg)
			fmt.Println("-- --------------------------------------------------- -- ")
			return fmt.Errorf("logMsg")
		}
	}
	return nil
}
