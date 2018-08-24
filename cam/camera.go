package cam

import (
	"github.com/mjgalindo/pt_weekend/ray"
	"github.com/mjgalindo/pt_weekend/vec"
)

type Camera struct {
	Origin, LowerLeft, Horizontal, Vertical vec.Vec3
}

func Default() Camera {
	return Camera{
		Origin:     vec.Make(0.0, 0.0, 0.0),
		LowerLeft:  vec.Make(-2.0, -1.0, -1.0),
		Horizontal: vec.Make(4.0, 0.0, 0.0),
		Vertical:   vec.Make(0.0, 2.0, 0.0)}
}

func (c Camera) GetRay(u, v float32) ray.Ray {
	return ray.Ray{
		Origin: c.Origin,
		Direction: vec.Sum(c.LowerLeft, vec.MulSingle(c.Horizontal, u),
			vec.MulSingle(c.Vertical, v))}
}
