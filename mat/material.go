package mat

import (
	"github.com/mjgalindo/pt_weekend/geo"
	"github.com/mjgalindo/pt_weekend/ray"
	"github.com/mjgalindo/pt_weekend/vec"
)

func Lambertian(albedo vec.Vec3) func(rIn ray.Ray, rec geo.HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray) {
	return func(rIn ray.Ray, rec geo.HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray) {
		target := vec.Sum(rec.P, rec.Normal, vec.RandomInUnitSphere())
		scattered = ray.Make(rec.P, vec.Sub(target, rec.P))
		attenuation = albedo
		absorbed = false
		return
	}
}

func Mirror(albedo vec.Vec3) func(rIn ray.Ray, rec geo.HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray) {
	return func(rIn ray.Ray, rec geo.HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray) {
		reflected := vec.Reflect(rIn.Direction.MakeUnit(), rec.Normal)
		scattered = ray.Make(rec.P, reflected)
		attenuation = albedo
		absorbed = vec.Dot(scattered.Direction, rec.Normal) <= 0
		return
	}
}

func Metal(albedo vec.Vec3, fuzz float32) func(rIn ray.Ray, rec geo.HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray) {
	return func(rIn ray.Ray, rec geo.HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray) {
		reflected := vec.Reflect(rIn.Direction.MakeUnit(), rec.Normal)
		scattered = ray.Make(rec.P, vec.Sum(reflected, vec.MulSingle(vec.RandomInUnitSphere(), fuzz)))
		attenuation = albedo
		absorbed = vec.Dot(scattered.Direction, rec.Normal) <= 0
		return
	}
}
