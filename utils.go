package goggles

import (
	"math"

	"github.com/gopherjs/gopherjs/js"
)

func CancelAnimationFrame(id int) {
	js.Global.Call("cancelAnimationFrame")
}

func DegToRad(degrees float32) float32 {
	return float32(degrees * float32(math.Pi) / 180.0)
}

func RequestAnimationFrame(callback func(float32)) int {
	return js.Global.Call("requestAnimationFrame", callback).Int()
}

func Error(message string) {
	js.Global.Call("alert", "Error: "+message)
}
