package geo

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/mjgalindo/pt_weekend/ray"
)

type BVHNode struct {
	Left, Right Hitable
	Box         AABB
}

func ConstructBVH(hitables []Hitable) BVHNode {
	axis := int(3.0 * rand.Float32())
	lessFunc := func(i, j int) bool {
		ah := hitables[i]
		bh := hitables[j]
		oka, boxLeft := ah.BoundingBox(0, 0)
		okb, boxRight := bh.BoundingBox(0, 0)
		if !oka || !okb {
			fmt.Println("Oopsies, no bounding box in bvh node")
		}
		if boxLeft.Min.I(axis)-boxRight.Min.I(axis) < 0.0 {
			return true
		}
		return false
	}
	var bvh BVHNode
	sort.Slice(hitables, lessFunc)
	if len(hitables) == 1 {
		bvh.Left, bvh.Right = hitables[0], hitables[0]
	} else if len(hitables) == 2 {
		bvh.Left, bvh.Right = hitables[0], hitables[1]
	} else {
		bvh.Left = ConstructBVH(hitables[:len(hitables)/2])
		bvh.Right = ConstructBVH(hitables[len(hitables)/2:])
	}
	okl, bbLeft := bvh.Left.BoundingBox(0, 0)
	okr, bbRight := bvh.Right.BoundingBox(0, 0)
	if !okl || !okr {
		fmt.Println("no bounding box in bvh node constructor")
	}
	bvh.Box = SurroundingBox(bbLeft, bbRight)
	return bvh
}

func (b BVHNode) Hit(r ray.Ray, tMin, tMax float32, rec *HitRecord) bool {
	if b.Box.Hit(r, tMin, tMax, rec) {
		var leftRec, rightRec HitRecord
		hitLeft := b.Left.Hit(r, tMin, tMax, &leftRec)
		hitRight := b.Right.Hit(r, tMin, tMax, &rightRec)
		if hitLeft && hitRight {
			if leftRec.T < rightRec.T {
				rec.Update(leftRec)
			} else {
				rec.Update(rightRec)
			}
			return true
		} else if hitLeft {
			rec.Update(leftRec)
			return true
		} else if hitRight {
			rec.Update(rightRec)
			return true
		}
	}
	return false
}

func (b BVHNode) BoundingBox(t0, t1 float32) (bool, AABB) {
	return true, b.Box
}
