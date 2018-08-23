package ray

import "github.com/mjgalindo/pt_weekend/vec"

// Ray represents a single 3d space ray to be cast
type Ray struct {
	Origin, Direction vec.Vec3
}

func Make(p, d vec.Vec3) Ray {
	return Ray{Origin: p, Direction: d}
}

func MakeUnit(p, d vec.Vec3) Ray {
	return Ray{Origin: p, Direction: d.MakeUnit()}
}
