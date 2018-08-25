package vec

import (
	"math/rand"
)

// RandomInUnitSphere returns a random point in a unit sphere
func RandomInUnitSphere() Vec3 {
	var p Vec3
	for ok := true; ok; ok = p.SquaredLength() >= 1.0 {
		p = Sub(
			MulSingle(
				Make(rand.Float32(), rand.Float32(), rand.Float32()),
				2.0), Make(1, 1, 1))
	}
	return p
}

// RandomInUnitDisk returns a random point in a unit disk
func RandomInUnitDisk() Vec3 {
	var p Vec3
	for ok := true; ok; ok = Dot(p, p) >= 1.0 {
		p = MulSingle(Sub(Make(rand.Float32(), rand.Float32(), 0), Make(1, 1, 0)), 2.0)
	}
	return p
}
