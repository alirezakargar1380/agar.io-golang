package agar

import (
	"fmt"
	"math"
)

type AgarPosition struct {
	X float64
	Y float64
}

func (agar *AgarPosition) GetAgarSpace() []map[string]float64 {
	var dir []map[string]float64
	for angle := 1; angle <= 360; angle++ {
		for r := 1; r <= 10; r++ {
			var x float64 = float64(r) * math.Sin(math.Pi*2*float64(angle)/360)
			var y float64 = float64(r) * math.Cos(math.Pi*2*float64(angle)/360)
			x = (agar.X + math.Round(float64(x*100))/100)
			y = (agar.Y + math.Round(float64(y*100))/100)
			myMap := make(map[string]float64, 0)
			myMap["x"] = x
			myMap["y"] = y
			dir = append(dir, myMap)

			// fmt.Println(x, y)
		}
	}
	return dir
}

var dfg int = 0

func CheckAgarSpace(dir []map[string]float64) {
	// fmt.Println(math.Round(dir[0]["x"]))
	// isThere := false

	for _, v := range dir {
		var x int = int(math.Round(v["x"]))
		// fmt.Println(x)
		var y int = int(math.Round(v["y"]))
		if x == int(250) && y == int(250) {
			dfg++
			fmt.Println("is there", dfg)
		}
	}
}
