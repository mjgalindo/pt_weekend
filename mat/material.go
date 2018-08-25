package mat

import (
	"math"
	"math/rand"

	"github.com/mjgalindo/pt_weekend/geo"
	"github.com/mjgalindo/pt_weekend/ray"
	"github.com/mjgalindo/pt_weekend/vec"
)

type ScatterFun func(rIn ray.Ray, rec geo.HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray)

func Lambertian(albedo vec.Vec3) ScatterFun {
	return func(rIn ray.Ray, rec geo.HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray) {
		target := vec.Sum(rec.P, rec.Normal, vec.RandomInUnitSphere())
		scattered = ray.Make(rec.P, vec.Sub(target, rec.P))
		attenuation = albedo
		absorbed = false
		return
	}
}

func Mirror(albedo vec.Vec3) ScatterFun {
	return func(rIn ray.Ray, rec geo.HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray) {
		reflected := vec.Reflect(rIn.Direction.MakeUnit(), rec.Normal)
		scattered = ray.Make(rec.P, reflected)
		attenuation = albedo
		absorbed = vec.Dot(scattered.Direction, rec.Normal) <= 0
		return
	}
}

func Metal(albedo vec.Vec3, fuzz float32) ScatterFun {
	return func(rIn ray.Ray, rec geo.HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray) {
		reflected := vec.Reflect(rIn.Direction.MakeUnit(), rec.Normal)
		scattered = ray.Make(rec.P, vec.Sum(reflected, vec.MulSingle(vec.RandomInUnitSphere(), fuzz)))
		attenuation = albedo
		absorbed = vec.Dot(scattered.Direction, rec.Normal) <= 0
		return
	}
}

func schlick(cosine, refIndex float32) float32 {
	r0 := (1 - refIndex) / (1 + refIndex)
	r0 = r0 * r0
	return r0 + (1-r0)*float32(math.Pow(float64((1-cosine)), 5))
}

func Dielectric(refIndex float32) ScatterFun {
	return func(rIn ray.Ray, rec geo.HitRecord) (absorbed bool, attenuation vec.Vec3, scattered ray.Ray) {
		var outwardNormal vec.Vec3
		reflected := vec.Reflect(rIn.Direction, rec.Normal)
		var niOverNt float32
		attenuation = vec.Make(1.0, 1.0, 1.0) // The glass surface absorbs nothing
		var refracted *vec.Vec3
		var reflectProb, cosine float32
		if vec.Dot(rIn.Direction, rec.Normal) > 0 {
			outwardNormal = rec.Normal.Neg()
			niOverNt = refIndex
			cosine = refIndex * vec.Dot(rIn.Direction, rec.Normal) / rIn.Direction.Length()
		} else {
			outwardNormal = rec.Normal
			niOverNt = 1.0 / refIndex
			cosine = -vec.Dot(rIn.Direction, rec.Normal) / rIn.Direction.Length()
		}
		if refracted = vec.Refract(rIn.Direction, outwardNormal, niOverNt); refracted != nil {
			reflectProb = schlick(cosine, refIndex)
			scattered = ray.Make(rec.P, *refracted)
		} else {
			scattered = ray.Make(rec.P, reflected)
			reflectProb = 1.0
		}
		if rand.Float32() < reflectProb {
			scattered = ray.Make(rec.P, reflected)
		} else {
			scattered = ray.Make(rec.P, *refracted)
		}
		absorbed = false
		return
	}
}
