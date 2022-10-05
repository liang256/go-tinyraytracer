package plane

import (
	"raytracer/material"

	"github.com/deeean/go-vector/vector3"
)

type Disc struct {
	Center   *vector3.Vector3
	Normal   *vector3.Vector3
	Radius   float64
	Material *material.Material
}

func (d *Disc) IntersectRay(orig, dir *vector3.Vector3) (bool, float64) {
	// formula: https://www.youtube.com/watch?v=x_SEyKtCBPU
	if dir.Dot(d.Normal) == 0 {
		return false, 0.0
	}
	dis := d.Center.Sub(orig).Dot(d.Normal) / (dir.Dot(d.Normal))
	hit := orig.Add(dir.Normalize().MulScalar(dis))
	if dis <= 0 || hit.Distance(d.Center) > d.Radius {
		return false, 0.0
	}
	return true, dis
}

func (d *Disc) GetMaterial() *material.Material {
	return d.Material
}

func (d *Disc) GetCenter() *vector3.Vector3 {
	return d.Center
}
