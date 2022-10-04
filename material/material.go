package material

import (
	"github.com/deeean/go-vector/vector3"
)

type Material struct {
	Color      *vector3.Vector3
	SpecColor  *vector3.Vector3
	Albedo     *vector3.Vector3
	SpecExpo   float64
	Refractive float64
}
