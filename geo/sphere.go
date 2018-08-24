package geo

import (
	"math"

	"github.com/mjgalindo/pt_weekend/ray"
	"github.com/mjgalindo/pt_weekend/vec"
)

type Sphere struct {
	Position vec.Vec3
	Radius   float32
}

func (s Sphere) Hit(r ray.Ray, tMin, tMax float32, rec *HitResult) bool {
	oc := vec.Sub(r.Origin, s.Position)
	a := vec.Dot(r.Direction, r.Direction)
	b := vec.Dot(oc, r.Direction)
	c := vec.Dot(oc, oc) - s.Radius*s.Radius
	discriminant := b*b - a*c
	if discriminant > 0 {
		temp := (-b - float32(math.Sqrt(float64(b*b-a*c)))) / a
		if temp < tMax && temp > tMin {
			rec.T = temp
			rec.P = r.PointAt(rec.T)
			rec.Normal = vec.DivSingle(vec.Sub(rec.P, s.Position), s.Radius)
			return true
		}
		temp = (-b + float32(math.Sqrt(float64(b*b-a*c)))) / a
		if temp < tMax && temp > tMin {
			rec.T = temp
			rec.P = r.PointAt(rec.T)
			rec.Normal = vec.DivSingle((vec.Sub(rec.P, s.Position)), s.Radius)
			return true
		}
	}
	return false
}
