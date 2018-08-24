package geo

import (
	"github.com/mjgalindo/pt_weekend/ray"
	"github.com/mjgalindo/pt_weekend/vec"
)

type HitRecord struct {
	T         float32
	P, Normal vec.Vec3
	Scatter   func(rIn ray.Ray, rec HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray)
}

func (hr *HitRecord) Update(res HitRecord) {
	hr.T = res.T
	hr.P = res.P
	hr.Normal = res.Normal
	hr.Scatter = res.Scatter
}

type Hitable interface {
	// Hit fills the hit result and returns true when hit
	Hit(r ray.Ray, tMin, tMax float32, rec *HitRecord) bool
}
