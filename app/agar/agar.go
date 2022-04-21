package agar

import (
	"fmt"
	"math"
	"strconv"
)

type AgarPosition struct {
	X float64
	Y float64
}

func (agar *AgarPosition) GetAgarSpace() []map[string]float64 {
	var dir []map[string]float64
	for angle := 1; angle <= 360; angle++ {
		for r := 1; r <= 60; r++ {
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

func (agar *AgarPosition) GetAgarSpace2(beads *map[string]int) []map[string]float64 {
	var dir []map[string]float64
	for angle := 1; angle <= 360; angle++ {
		for r := 1; r <= 60; r++ {
			var x float64 = float64(r) * math.Sin(math.Pi*2*float64(angle)/360)
			var y float64 = float64(r) * math.Cos(math.Pi*2*float64(angle)/360)
			var xx int = int(agar.X + math.Round(float64(x*100))/100)
			var yy int = int(agar.Y + math.Round(float64(y*100))/100)

			if (*beads)[strconv.Itoa(xx)+"_"+strconv.Itoa(yy)] == 10 {
				fmt.Println("found")
			}
			// for key := range *beads {
			// 	positions := strings.Split(key, "_")
			// 	beadX, error := strconv.Atoi(positions[0])
			// 	if error != nil {
			// 		fmt.Println(error)
			// 	}
			// 	beadY, error := strconv.Atoi(positions[1])
			// 	if error != nil {
			// 		fmt.Println(error)
			// 	}
			// 	if beadX == xx && beadY == yy {
			// 		fmt.Println("found")
			// 	}
			// }
		}
	}
	return dir
}

func CheckAgarSpace(dir []map[string]float64, beads *map[string]int) bool {
	var eat bool = false
	// fmt.Println(len(dir))
	for _, v := range dir {
		if len(*beads) == 0 {
			return false
		}
		var x string = fmt.Sprintf("%v", math.Round(v["x"]))
		var y string = fmt.Sprintf("%v", math.Round(v["y"]))
		if (*beads)[x+"_"+y] == 10 {
			delete(*beads, x+"_"+y)
			// fmt.Println("found")
			return true
		}
		// for key := range *beads {
		// 	positions := strings.Split(key, "_")
		// 	beadX, error := strconv.Atoi(positions[0])
		// 	if error != nil {
		// 		fmt.Println(error)
		// 	}
		// 	beadY, error := strconv.Atoi(positions[1])
		// 	if error != nil {
		// 		fmt.Println(error)
		// 	}
		// 	var x int = int(math.Round(v["x"]))
		// 	var y int = int(math.Round(v["y"]))
		// 	if x == beadX && y == beadY {
		// 		fmt.Println("is there")
		// 	}

	}
	return eat
	// }
}
