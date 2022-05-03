package trigonometric_circle

import (
	"math"
)

type AgarDetail struct {
	Id     int
	X      float64
	Y      float64
	Speed  float64
	Radius float64
}

type AgarDe struct {
	Id        int
	X         float64
	Y         float64
	Name      string
	Radius    float64
	Max_speed float64
	Speed     float64
}

type EatTogetherResponse struct {
	AgarId int
	Eat    bool
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

func (agar *AgarDetail) GetDistance(to_x float64, to_y float64) float64 {
	var x float64 = to_y - float64(agar.Y)
	var y float64 = to_x - float64(agar.X)
	return math.Sqrt(x*x + y*y)
}

func (agar *AgarDetail) CheckForEatTogether(agars []AgarDe) EatTogetherResponse {
	var response EatTogetherResponse = EatTogetherResponse{
		AgarId: 0,
		Eat:    false,
	}
	for i := 0; i < len(agars); i++ {
		agarObject := agars[i]
		if agarObject.Id != agar.Id {
			if agar.GetDistance(agarObject.X, agarObject.Y) < agarObject.Radius {
				response.AgarId = agar.Id
				response.Eat = true
			}
		}
	}

	return response
}
