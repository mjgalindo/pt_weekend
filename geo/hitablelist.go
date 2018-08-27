package geo

import "github.com/mjgalindo/pt_weekend/ray"

type HitableList struct {
	Hitables []Hitable
}

func (hl *HitableList) Add(hs ...Hitable) {
	hl.Hitables = append(hl.Hitables, hs...)
}

func MakeList(hitables ...Hitable) HitableList {
	return HitableList{hitables}
}

func (hl HitableList) Hit(r ray.Ray, tMin, tMax float32, rec *HitRecord) bool {
	hrec := HitRecord{}
	hasHit := false
	closest := tMax
	for _, hitable := range hl.Hitables {
		if hitable.Hit(r, tMin, closest, &hrec) {
			hasHit = true
			closest = hrec.T
			rec.Update(hrec)
		}
	}
	return hasHit
}

func (hl HitableList) BoundingBox(t0, t1 float32) (bool, AABB) {
	if len(hl.Hitables) < 1 {
		return false, AABB{}
	}
	firstTrue, box := hl.Hitables[0].BoundingBox(t0, t1)
	if !firstTrue {
		return false, AABB{}
	}
	for _, hitable := range hl.Hitables {
		if ok, tmpBB := hitable.BoundingBox(t0, t1); ok {
			box = SurroundingBox(box, tmpBB)
		} else {
			return false, AABB{}
		}
	}
	return true, box
}
