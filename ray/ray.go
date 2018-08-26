package ray

import "github.com/mjgalindo/pt_weekend/vec"

// Ray represents a single 3d space ray to be cast
type Ray struct {
	Origin, Direction vec.Vec3
	Time              float32
}

func Make(p, d vec.Vec3, time float32) Ray {
	return Ray{Origin: p, Direction: d, Time: time}
}

func MakeUnit(p, d vec.Vec3) Ray {
	return Ray{Origin: p, Direction: d.MakeUnit()}
}

func (r Ray) PointAt(t float32) vec.Vec3 {
	return vec.Sum(r.Origin, vec.MulSingle(r.Direction, t))
}
