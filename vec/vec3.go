package vec

import (
	"fmt"
	"math"
)

// Vec3 is a simple linalg 3 float vector type
type Vec3 struct {
	e [3]float32
}

func Make(x, y, z float32) Vec3 {
	return Vec3{[3]float32{x, y, z}}
}

func (v Vec3) X() float32 {
	return v.e[0]
}

func (v Vec3) Y() float32 {
	return v.e[1]
}

func (v Vec3) Z() float32 {
	return v.e[2]
}

func (v Vec3) XYZ() (float32, float32, float32) {
	return v.e[0], v.e[1], v.e[2]
}

func (v Vec3) R() float32 {
	return v.e[0]
}

func (v Vec3) G() float32 {
	return v.e[1]
}

func (v Vec3) B() float32 {
	return v.e[2]
}

func (v Vec3) RGB() (float32, float32, float32) {
	return v.e[0], v.e[1], v.e[2]
}

func (v Vec3) Neg() Vec3 {
	return Vec3{e: [3]float32{-v.e[0], -v.e[1], -v.e[2]}}
}

func Sum(a Vec3, bs ...Vec3) Vec3 {
	x, y, z := a.XYZ()
	for _, b := range bs {
		x, y, z = x+b.X(), y+b.Y(), z+b.Z()
	}
	return Vec3{e: [3]float32{x, y, z}}
}

func Sub(a Vec3, bs ...Vec3) Vec3 {
	x, y, z := a.XYZ()
	for _, b := range bs {
		x, y, z = x-b.X(), y-b.Y(), z-b.Z()
	}
	return Vec3{e: [3]float32{x, y, z}}
}

func Mul(a Vec3, bs ...Vec3) Vec3 {
	x, y, z := a.XYZ()
	for _, b := range bs {
		x, y, z = x*b.X(), y*b.Y(), z*b.Z()
	}
	return Vec3{e: [3]float32{x, y, z}}
}

func Div(a, b Vec3) Vec3 {
	return Vec3{e: [3]float32{a.e[0] / b.e[0], a.e[1] / b.e[1], a.e[2] / b.e[2]}}
}
func SumSingle(v Vec3, s float32) Vec3 {
	return Vec3{e: [3]float32{v.e[0] + s, v.e[1] + s, v.e[2] + s}}
}

func SubSingle(v Vec3, s float32) Vec3 {
	return Vec3{e: [3]float32{v.e[0] - s, v.e[1] - s, v.e[2] - s}}
}

func MulSingle(v Vec3, s float32) Vec3 {
	return Vec3{e: [3]float32{v.e[0] * s, v.e[1] * s, v.e[2] * s}}
}

func DivSingle(v Vec3, s float32) Vec3 {
	return Vec3{e: [3]float32{v.e[0] / s, v.e[1] / s, v.e[2] / s}}
}

func (v Vec3) SquaredLength() float32 {
	return v.e[0]*v.e[0] + v.e[1]*v.e[1] + v.e[2]*v.e[2]
}

func (v Vec3) Length() float32 {
	return float32(math.Sqrt(float64(v.SquaredLength())))
}

func (v Vec3) MakeUnit() Vec3 {
	k := 1.0 / v.Length()
	return MulSingle(v, k)
}

func Dot(a, b Vec3) float32 {
	return a.e[0]*b.e[0] + a.e[1]*b.e[1] + a.e[2]*b.e[2]
}

func Cross(a, b Vec3) Vec3 {
	return Vec3{[3]float32{
		a.e[1]*b.e[2] - a.e[2]*b.e[1],
		-(a.e[0]*b.e[2] - a.e[2]*b.e[0]),
		a.e[0]*b.e[1] - a.e[1]*b.e[0]}}
}

func Reflect(v, n Vec3) Vec3 {
	return Sub(v, MulSingle(n, 2*Dot(v, n)))
}

func Refract(v, n Vec3, niOverNt float32) (refracted *Vec3) {
	uv := v.MakeUnit()
	dt := Dot(uv, n)
	discriminant := 1.0 - niOverNt*niOverNt*(1-dt*dt)
	if discriminant > 0 {
		tmp := Sub(MulSingle(Sub(uv, MulSingle(n, dt)), niOverNt), MulSingle(n, float32(math.Sqrt(float64(discriminant)))))
		refracted = &tmp
	}
	return
}

func (v Vec3) String() string {
	return fmt.Sprintf("%.5f %.5f %.5f", v.e[0], v.e[1], v.e[2])
}
