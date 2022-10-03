package sphere

import (
	"math"
	"raytracer/material"

	"github.com/deeean/go-vector/vector3"
)

type Sphere struct {
	Center   *vector3.Vector3
	Radius   float64
	Material *material.Material
}

// Returns isIntersect, distant from point to camera, and normal of the point
func (s Sphere) IntersectRay(orig, dir *vector3.Vector3) (bool, float64) {
	// if the sphere not at the side of ray direction return false
	origCenterVec := s.Center.Sub(orig)
	ray := orig.Add(dir)
	dotProd := origCenterVec.Dot(ray)
	if dotProd < 0 {
		return false, 0.0
	}
	// if the ray orig in the sphere, return false
	if origCenterVec.Magnitude() <= s.Radius {
		return false, 0.0
	}
	// get the projected center on the ray, pc
	// length of projection
	projLen := dotProd / ray.Magnitude()
	// pc = orig + normalized-ray * projectlength
	pc := orig.Add(ray.Normalize().MulScalar(projLen))
	// if |pc - c| <= radius, return true
	if res := pc.Sub(s.Center).Magnitude(); res > s.Radius {
		return false, 0.0
	}
	// caculate the distance between the closest point on this sphere and the camera
	dis := projLen - math.Sqrt(math.Pow(s.Radius, 2)-math.Pow(origCenterVec.Magnitude(), 2)+math.Pow(projLen, 2))
	return true, dis
}
