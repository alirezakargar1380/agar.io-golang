package agar

import (
	"math"
)

func GetAgarSpace(x float64, y float64) []map[string]float64 {
	var dir []map[string]float64
	for angle := 1; angle <= 360; angle++ {
		for r := 1; r <= 25; r++ {
			var x float64 = float64(r) * math.Sin(math.Pi*2*float64(angle)/360)
			var y float64 = float64(r) * math.Cos(math.Pi*2*float64(angle)/360)
			x = (300 + math.Round(float64(x*100))/100)
			y = (100 + math.Round(float64(y*100))/100)
			myMap := make(map[string]float64, 0)
			myMap["x"] = x
			myMap["y"] = y
			dir = append(dir, myMap)

			// fmt.Println(x, y)
		}
	}
	return dir
}
