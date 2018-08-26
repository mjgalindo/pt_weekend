package geo

import (
	"math"

	"github.com/mjgalindo/pt_weekend/ray"
	"github.com/mjgalindo/pt_weekend/vec"
)

type Sphere struct {
	Position     vec.Vec3
	LastPosition vec.Vec3
	Radius       float32
	ScatterFun   func(rIn ray.Ray, rec HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray)
	Time0, Time1 float32
}

func (s Sphere) PositionAt(time float32) vec.Vec3 {
	return vec.Sum(s.Position, vec.MulSingle(vec.Sub(s.LastPosition, s.Position), ((time-s.Time0)/(s.Time1-s.Time0))))
}

func MakeSphere(position vec.Vec3, radius float32,
	scatterFun func(rIn ray.Ray, rec HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray)) Sphere {
	return Sphere{
		Position:     position,
		LastPosition: position,
		Radius:       radius,
		ScatterFun:   scatterFun,
		Time0:        0.0,
		Time1:        1.0}
}

func MakeMovingSphere(position0, position1 vec.Vec3, radius float32, time0, time1 float32, scatterFun func(rIn ray.Ray, rec HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray)) Sphere {
	return Sphere{
		Position:     position0,
		LastPosition: position1,
		Radius:       radius,
		ScatterFun:   scatterFun,
		Time0:        time0,
		Time1:        time1}

}

func (s Sphere) Hit(r ray.Ray, tMin, tMax float32, rec *HitRecord) bool {
	oc := vec.Sub(r.Origin, s.PositionAt(r.Time))
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
			rec.Normal = vec.DivSingle(vec.Sub(rec.P, s.PositionAt(r.Time)), s.Radius)
			return true
		}
		temp = (-b + float32(math.Sqrt(float64(b*b-a*c)))) / a
		if temp < tMax && temp > tMin {
			rec.T = temp
			rec.P = r.PointAt(rec.T)
			rec.Normal = vec.DivSingle((vec.Sub(rec.P, s.PositionAt(r.Time))), s.Radius)
			return true
		}
	}
	return false
}
