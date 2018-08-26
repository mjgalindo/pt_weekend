package cam

import (
	"math"
	"math/rand"

	"github.com/mjgalindo/pt_weekend/ray"
	"github.com/mjgalindo/pt_weekend/vec"
)

type Camera struct {
	Origin, LowerLeft, Horizontal, Vertical, U, V, W vec.Vec3
	Time0, Time1                                     float32
	LensRadius                                       float32
}

func Make(lookFrom, lookAt, vUp vec.Vec3, vFov, aspect, aperture, focusDist, time0, time1 float32) Camera {
	var u, v, w vec.Vec3
	theta := vFov * math.Pi / 180
	halfHeight := float32(math.Tan(float64(theta) / 2.0))
	halfWidth := aspect * halfHeight
	w = vec.Sub(lookFrom, lookAt).MakeUnit()
	u = vec.Cross(vUp, w).MakeUnit()
	v = vec.Cross(w, u)
	return Camera{
		Origin:     lookFrom,
		LowerLeft:  vec.Sub(lookFrom, vec.MulSingle(u, halfWidth*focusDist), vec.MulSingle(v, halfHeight*focusDist), vec.MulSingle(w, focusDist)),
		Horizontal: vec.MulSingle(u, 2*halfWidth*focusDist),
		Vertical:   vec.MulSingle(v, 2*halfHeight*focusDist),
		U:          u,
		V:          v,
		W:          w,
		LensRadius: aperture / 2,
		Time0:      time0,
		Time1:      time1}
}

func (c Camera) GetRay(s, t float32) ray.Ray {
	rd := vec.MulSingle(vec.RandomInUnitDisk(), c.LensRadius)
	offset := vec.Sum(vec.MulSingle(c.U, rd.X()), vec.MulSingle(c.V, rd.Y()))
	time := c.Time0 + rand.Float32()*(c.Time1-c.Time0)
	return ray.Ray{
		Origin: vec.Sum(c.Origin, offset),
		Direction: vec.Sum(c.LowerLeft, vec.MulSingle(c.Horizontal, s),
			vec.MulSingle(c.Vertical, t), c.Origin.Neg(), offset.Neg()),
		Time: time}
}
