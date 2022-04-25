package agar

import (
	"fmt"
	"math"
)

type AgarPosition struct {
	X      float64
	Y      float64
	Radius int
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

// var i int = 0

type Re struct {
	Eat     bool
	Eat_key string
}

func (agar *AgarPosition) GetAgarSpace2(beads *map[string]map[string]int, RoomId string) Re {
	var dir []map[string]int = make([]map[string]int, 0)
	for angle := 1; angle <= 360; angle++ {
		for r := agar.Radius; r >= (agar.Radius - 10); r-- {
			var x float64 = float64(r) * math.Sin(math.Pi*2*float64(angle)/360)
			var y float64 = float64(r) * math.Cos(math.Pi*2*float64(angle)/360)
			var xx int = int(agar.X + math.Round(float64(x*100))/100)
			var yy int = int(agar.Y + math.Round(float64(y*100))/100)

			myMap := make(map[string]int, 0)
			myMap["x"] = xx
			myMap["y"] = yy
			dir = append(dir, myMap)

			// if (*beads)[RoomId][strconv.Itoa(xx)+"_"+strconv.Itoa(yy)] == 10 {
			// 	var str_x = fmt.Sprintf("%v", xx)
			// 	var str_y = fmt.Sprintf("%v", yy)
			// 	delete((*beads)[RoomId], str_x+"_"+str_y)
			// 	return true
			// } else {
			// 	return false
			// }

		}
	}

	for _, v := range dir {
		var x string = fmt.Sprintf("%v", v["x"])
		var y string = fmt.Sprintf("%v", v["y"])
		if (*beads)[RoomId][x+"_"+y] == 10 {
			delete((*beads)[RoomId], x+"_"+y)
			// fmt.Println("found")
			return Re{
				Eat:     true,
				Eat_key: x + "_" + y,
			}
		}
	}

	return Re{
		Eat:     false,
		Eat_key: "",
	}
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
