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
