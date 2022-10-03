package material

import (
	"github.com/deeean/go-vector/vector2"
	"github.com/deeean/go-vector/vector3"
)

type Material struct {
	Color     *vector3.Vector3
	SpecColor *vector3.Vector3
	Albedo    *vector2.Vector2
	SpecExpo  float64
}
