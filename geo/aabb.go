package geo

import (
	"github.com/mjgalindo/pt_weekend/ray"
	"github.com/mjgalindo/pt_weekend/vec"
)

// AABB is an axis aligned bounding box
type AABB struct {
	Min, Max vec.Vec3
}

func SurroundingBox(a, b AABB) AABB {
	min := vec.Min(a.Min, b.Min)
	max := vec.Max(a.Max, b.Max)
	return AABB{Min: min, Max: max}
}

func (bb AABB) HitOriginal(r ray.Ray, tMin, tMax float32, rec *HitRecord) bool {
	for a := 0; a < 3; a++ {
		t0 := min((bb.Min.I(a)-r.Origin.I(a))/r.Direction.I(a),
			(bb.Max.I(a)-r.Origin.I(a))/r.Direction.I(a))
		t1 := max((bb.Min.I(a)-r.Origin.I(a))/r.Direction.I(a),
			(bb.Max.I(a)-r.Origin.I(a))/r.Direction.I(a))
		tMin = max(t0, tMin)
		tMax = min(t1, tMax)
		if tMax <= tMin {
			return false
		}
	}
	return true
}

func (bb AABB) Hit(r ray.Ray, tMin, tMax float32, rec *HitRecord) bool {
	for a := 0; a < 3; a++ {
		invD := 1.0 / r.Direction.I(a)
		t0 := (bb.Min.I(a) - r.Origin.I(a)) * invD
		t1 := (bb.Max.I(a) - r.Origin.I(a)) * invD
		if invD < 0.0 {
			t0, t1 = t1, t0
		}
		if t0 > tMin {
			tMin = t0
		}
		if t1 < tMax {
			tMax = t1
		}
		if tMax <= tMin {
			return false
		}
	}
	return true
}

func min(a, b float32) float32 {
	if a > b {
		return b
	}
	return a
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
