package geo

import (
	"github.com/mjgalindo/pt_weekend/ray"
	"github.com/mjgalindo/pt_weekend/vec"
)

type HitResult struct {
	T         float32
	P, Normal vec.Vec3
}

func (hr *HitResult) Update(res HitResult) {
	hr.T = res.T
	hr.P = res.P
	hr.Normal = res.Normal
}

type Hitable interface {
	// Hit fills the hit result and returns true when hit
	Hit(r ray.Ray, tMin, tMax float32, rec *HitResult) bool
}
