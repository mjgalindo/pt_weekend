package geo

import (
	"math"

	"github.com/mjgalindo/pt_weekend/ray"
	"github.com/mjgalindo/pt_weekend/vec"
)

type Sphere struct {
	Position   vec.Vec3
	Radius     float32
	ScatterFun func(rIn ray.Ray, rec HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray)
}

func MakeSphere(position vec.Vec3, radius float32,
	scatterFun func(rIn ray.Ray, rec HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray)) Sphere {
	return Sphere{Position: position, Radius: radius, ScatterFun: scatterFun}
}

func (s Sphere) Hit(r ray.Ray, tMin, tMax float32, rec *HitRecord) bool {
	oc := vec.Sub(r.Origin, s.Position)
	a := vec.Dot(r.Direction, r.Direction)
	b := vec.Dot(oc, r.Direction)
	c := vec.Dot(oc, oc) - s.Radius*s.Radius
	rec.Scatter = s.ScatterFun
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

func (s Sphere) BoundingBox(t0, t1 float32) (bool, AABB) {
	rVec := vec.Make(s.Radius, s.Radius, s.Radius)
	min0 := vec.Sub(s.Position, rVec)
	max0 := vec.Sum(s.Position, rVec)
	return true, AABB{
		Min: min0,
		Max: max0}
}
