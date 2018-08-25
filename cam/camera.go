package cam

import (
	"math"

	"github.com/mjgalindo/pt_weekend/ray"
	"github.com/mjgalindo/pt_weekend/vec"
)

type Camera struct {
	Origin, LowerLeft, Horizontal, Vertical vec.Vec3
}

func Default(vFov, aspect float32) Camera {
	theta := vFov * math.Pi / 180
	halfHeight := float32(math.Tan(float64(theta) / 2.0))
	halfWidth := aspect * halfHeight
	return Camera{
		Origin:     vec.Make(0.0, 0.0, 0.0),
		LowerLeft:  vec.Make(-halfWidth, -halfHeight, -1.0),
		Horizontal: vec.Make(2*halfWidth, 0.0, 0.0),
		Vertical:   vec.Make(0.0, 2.0*halfHeight, 0.0)}
}

func Make(lookFrom, lookAt, vUp vec.Vec3, vFov, aspect float32) Camera {
	var u, v, w vec.Vec3
	theta := vFov * math.Pi / 180
	halfHeight := float32(math.Tan(float64(theta) / 2.0))
	halfWidth := aspect * halfHeight
	w = vec.Sub(lookFrom, lookAt).MakeUnit()
	u = vec.Cross(vUp, w).MakeUnit()
	v = vec.Cross(w, u)
	return Camera{
		Origin:     lookFrom,
		LowerLeft:  vec.Sub(lookFrom, vec.MulSingle(u, halfWidth), vec.MulSingle(v, halfHeight), w),
		Horizontal: vec.MulSingle(u, 2*halfWidth),
		Vertical:   vec.MulSingle(v, 2*halfHeight)}
}

func (c Camera) GetRay(s, t float32) ray.Ray {
	return ray.Ray{
		Origin: c.Origin,
		Direction: vec.Sum(c.LowerLeft, vec.MulSingle(c.Horizontal, s),
			vec.MulSingle(c.Vertical, t), c.Origin.Neg())}
}
