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
