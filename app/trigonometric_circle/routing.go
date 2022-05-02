package trigonometric_circle

import (
	"math"
)

type AgarDetail struct {
	X      float64
	Y      float64
	Speed  float64
	Radius float64
}

func (agar *AgarDetail) Test(angle float64) map[string]float64 {
	var x float64 = agar.Speed * math.Sin(math.Pi*2*angle/360)
	var y float64 = agar.Speed * math.Cos(math.Pi*2*angle/360)
	x = agar.X + x
	y = agar.Y + y
	response := make(map[string]float64)
	response["x"] = x
	response["y"] = y
	return response
}
