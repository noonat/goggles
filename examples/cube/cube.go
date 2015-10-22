package main

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"github.com/noonat/goggles"
)

const fragmentShaderSource = `
precision mediump float;

varying vec4 outVertexColor;

void main(void) {
  gl_FragColor = outVertexColor;
}
`

const vertexShaderSource = `
attribute vec3 inVertexPosition;
attribute vec3 inVertexColor;

uniform mat4 modelViewMatrix;
uniform mat4 projectionMatrix;
uniform mat4 modelViewProjectionMatrix;

varying vec4 outVertexColor;

void main(void) {
  gl_Position = projectionMatrix * modelViewMatrix * vec4(inVertexPosition, 1.0);
  outVertexColor = vec4(inVertexColor, 1.0);
}
`

var vertexPositionData = []float32{
	// near face
	-1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, -1.0, 1.0,

	// left face
	-1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0,

	// far face
	1.0, -1.0, -1.0,
	1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,

	// right face
	1.0, -1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, -1.0, -1.0,

	// top face
	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, 1.0,

	// bottom face
	-1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, 1.0,
	1.0, -1.0, -1.0,
}

var vertexColorData = []float32{
	// near face
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,

	// left face
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,

	// far face
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,

	// right face
	1.0, 1.0, 0.0,
	1.0, 1.0, 0.0,
	1.0, 1.0, 0.0,
	1.0, 1.0, 0.0,

	// top face
	1.0, 0.0, 1.0,
	1.0, 0.0, 1.0,
	1.0, 0.0, 1.0,
	1.0, 0.0, 1.0,

	// bottom face
	0.0, 1.0, 1.0,
	0.0, 1.0, 1.0,
	0.0, 1.0, 1.0,
	0.0, 1.0, 1.0,
}

var indexData = []uint16{
	0, 1, 2,
	0, 2, 3,
	4, 5, 6,
	4, 6, 7,
	8, 9, 10,
	8, 10, 11,
	12, 13, 14,
	12, 14, 15,
	16, 17, 18,
	16, 18, 19,
	20, 21, 22,
	20, 22, 23,
}

const viewportWidth float32 = 500
const viewportHeight float32 = 500

func main() {
	document := js.Global.Get("document")

	var viewportScale float32 = 2.0

	canvas := document.Call("createElement", "canvas")
	canvas.Set("id", "canvas")
	canvas.Set("width", viewportWidth*viewportScale)
	canvas.Set("height", viewportHeight*viewportScale)
	canvas.Get("style").Set("width", fmt.Sprintf("%dpx", int(viewportWidth)))
	canvas.Get("style").Set("height", fmt.Sprintf("%dpx", int(viewportHeight)))
	document.Get("body").Call("appendChild", canvas)

	textCanvas := document.Call("createElement", "canvas")
	textCanvas.Set("id", "text-canvas")
	textCanvas.Set("width", viewportWidth*viewportScale)
	textCanvas.Set("height", viewportHeight*viewportScale)
	textCanvas.Get("style").Set("width", fmt.Sprintf("%dpx", int(viewportWidth)))
	textCanvas.Get("style").Set("height", fmt.Sprintf("%dpx", int(viewportHeight)))
	document.Get("body").Call("appendChild", textCanvas)

	textContext := textCanvas.Call("getContext", "2d")
	textContext.Call("setTransform", viewportScale, 0, 0, viewportScale, 0, 0)

	attrs := webgl.DefaultAttributes()
	attrs.Alpha = false

	gl, err := webgl.NewContext(canvas, attrs)
	if err != nil {
		goggles.Error(err.Error())
		return
	}

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragmentShader, fragmentShaderSource)
	gl.CompileShader(fragmentShader)
	if !gl.GetShaderParameterb(fragmentShader, gl.COMPILE_STATUS) {
		goggles.Error(gl.GetShaderInfoLog(fragmentShader))
		return
	}

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertexShader, vertexShaderSource)
	gl.CompileShader(vertexShader)
	if !gl.GetShaderParameterb(vertexShader, gl.COMPILE_STATUS) {
		goggles.Error(gl.GetShaderInfoLog(vertexShader))
		return
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	if !gl.GetProgramParameterb(program, gl.LINK_STATUS) {
		goggles.Error(gl.GetProgramInfoLog(program))
		return
	}

	gl.UseProgram(program)

	inVertexPositionLocation := gl.GetAttribLocation(program, "inVertexPosition")
	inVertexColorLocation := gl.GetAttribLocation(program, "inVertexColor")
	modelViewMatrixLocation := gl.GetUniformLocation(program, "modelViewMatrix")
	projectionMatrixLocation := gl.GetUniformLocation(program, "projectionMatrix")

	gl.EnableVertexAttribArray(inVertexPositionLocation)
	gl.EnableVertexAttribArray(inVertexColorLocation)

	vertexPositionBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexPositionBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, vertexPositionData, gl.STATIC_DRAW)

	vertexColorBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexColorBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, vertexColorData, gl.STATIC_DRAW)

	indexBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indexData, gl.STATIC_DRAW)

	gl.Viewport(0, 0, int(viewportWidth*viewportScale), int(viewportHeight*viewportScale))

	var pitch float32 = 0.0
	var yaw float32 = 0.0
	var tweenTime float32 = 0.0

	var tick func(float32)

	tick = func(timeMs float32) {
		textContext.Call("clearRect", 0, 0, viewportWidth, viewportHeight)
		textContext.Set("font", "12px sans-serif")
		textContext.Set("fillStyle", "rgb(255, 255, 255)")
		textContext.Call("fillText", "Hello, world!", 10, 20)

		time := timeMs / 1000

		if tweenTime == 0.0 {
			tweenTime = time + 1.0
		}

		for tweenTime < time {
			tweenTime += 1.0
			pitch = float32(math.Mod(float64(pitch)+60.0, 360.0))
			yaw = float32(math.Mod(float64(yaw)+40.0, 360.0))
		}

		factor := tweenTime - time
		if factor < 0.0 {
			factor = 0.0
		} else if factor > 1.0 {
			factor = 1.0
		}
		factor = 1.0 - float32(math.Pow(float64(factor), 4.0))

		tweenPitch := pitch + (60.0 * factor)
		tweenYaw := yaw + (40.0 * factor)

		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Enable(gl.CULL_FACE)
		gl.FrontFace(gl.CW)

		modelViewMatrix := mgl32.Ident4().
			Mul4(mgl32.Translate3D(0, 0, -4)).
			Mul4(mgl32.HomogRotate3DX(goggles.DegToRad(tweenPitch))).
			Mul4(mgl32.HomogRotate3DY(goggles.DegToRad(tweenYaw)))
		gl.UniformMatrix4fv(modelViewMatrixLocation, false, modelViewMatrix[:])

		projectionMatrix := mgl32.Perspective(goggles.DegToRad(60), viewportWidth/viewportHeight, 0.1, 2048)
		gl.UniformMatrix4fv(projectionMatrixLocation, false, projectionMatrix[:])

		gl.BindBuffer(gl.ARRAY_BUFFER, vertexPositionBuffer)
		gl.VertexAttribPointer(inVertexPositionLocation, 3, gl.FLOAT, false, 0, 0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertexColorBuffer)
		gl.VertexAttribPointer(inVertexColorLocation, 3, gl.FLOAT, false, 0, 0)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
		gl.DrawElements(gl.TRIANGLES, len(indexData), gl.UNSIGNED_SHORT, 0)

		goggles.RequestAnimationFrame(tick)
	}

	goggles.RequestAnimationFrame(tick)
}
