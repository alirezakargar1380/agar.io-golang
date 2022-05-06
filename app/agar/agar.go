package agar

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/alirezakargar1380/agar.io-golang/app/beads"
	"github.com/alirezakargar1380/agar.io-golang/app/trigonometric_circle"
)

type AgarPosition struct {
	X      float64
	Y      float64
	Radius int
}

type AllAgars struct {
	Agars    []trigonometric_circle.AgarDe
	ClientId int
	RivalId  int
	Id       int
	X        float64
	Y        float64
	Radius   int
}

type CheckForEatingResponse struct {
	Status           bool
	Eated_agar_id    int
	Eated_agar_by_id int
}

func (agars *AllAgars) CheckForAgarEatingOtherAgars() {
	for _, v := range agars.Agars {
		distance := trigonometric_circle.GetDistanceBetweenTowPoint(v.X, v.Y, agars.X, agars.Y)
		if distance < float64(agars.Radius)+v.Radius {
			if float64(agars.Radius) > v.Radius {
				// fmt.Println(agars.Id, "for", agars.ClientId)
				fmt.Println(agars.ClientId, "is eating", agars.RivalId)
			} else {
				fmt.Println(agars.RivalId, "is eating", agars.ClientId)
			}
		}
	}
}

func (agars *AllAgars) CheckForEating() CheckForEatingResponse {
	var response CheckForEatingResponse = CheckForEatingResponse{
		Status: false,
	}
	for _, v := range agars.Agars {
		if v.Id == agars.Id {
			continue
		}
		distance := trigonometric_circle.GetDistanceBetweenTowPoint(v.X, v.Y, agars.X, agars.Y)
		if distance < float64(agars.Radius)+v.Radius {
			response.Status = true
			if v.Id == 1 {
				response.Eated_agar_id = agars.Id
			}
			if agars.Id == 1 {
				response.Eated_agar_id = v.Id
			}

			if response.Eated_agar_id == agars.Id {
				response.Eated_agar_by_id = v.Id
			} else {
				response.Eated_agar_by_id = agars.Id
			}
		}
	}
	return response
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

type Re struct {
	Eat     bool
	Eat_key []string
}

type Eat struct {
	eat bool
	key string
}

func (agar *AgarPosition) getDistance(to_x float64, to_y float64) float64 {
	var x float64 = to_y - float64(agar.Y)
	var y float64 = to_x - float64(agar.X)
	return math.Sqrt(x*x + y*y)
}

func (agar *AgarPosition) GetAgarSpace4(beads *beads.Beads, RoomId string) Re {
	var eat bool = false
	var eatKey []string = make([]string, 0)
	for key := range beads.Beads[RoomId] {
		positions := strings.Split(key, "_")
		beadX, error := strconv.Atoi(positions[0])
		if error != nil {
			fmt.Println(error)
		}
		beadY, error := strconv.Atoi(positions[1])
		if error != nil {
			fmt.Println(error)
		}
		if agar.getDistance(float64(beadX), float64(beadY)) < float64(agar.Radius) {
			eat = true
			// eatKey = key
			eatKey = append(eatKey, key)
			delete(beads.Beads[RoomId], key)
		}
	}

	return Re{
		Eat:     eat,
		Eat_key: eatKey,
	}
}

// func (agar *AgarPosition) GetAgarSpace3(beads *beads.Beads, RoomId string) Re {
// 	var wg sync.WaitGroup
// 	c1 := make(chan Eat)
// 	checkBeadIsExist := func(room string) {
// 		radius14 := math.Round(float64(agar.Radius / 3))
// 		var eatKey Eat = Eat{
// 			eat: false,
// 			key: "",
// 		}
// 		for angle := 1; angle <= 360; angle++ {
// 			for r := agar.Radius; r >= (agar.Radius - int(radius14)); r-- {
// 				var x float64 = float64(r) * math.Sin(math.Pi*2*float64(angle)/360)
// 				var y float64 = float64(r) * math.Cos(math.Pi*2*float64(angle)/360)
// 				var xx int = int(agar.X + math.Round(float64(x*100))/100)
// 				var yy int = int(agar.Y + math.Round(float64(y*100))/100)
// 				sx := fmt.Sprintf("%v", xx)
// 				sy := fmt.Sprintf("%v", yy)
// 				existRes := beads.Exist(room, sx+"_"+sy)
// 				if existRes {
// 					eatKey = Eat{
// 						eat: true,
// 						key: sx + "_" + sy,
// 					}
// 					beads.DeleteWithKey(room, sx+"_"+sy)
// 				}
// 			}
// 		}
// 		wg.Done()
// 		c1 <- eatKey
// 	}

// 	wg.Add(1)
// 	go checkBeadIsExist(RoomId)
// 	wg.Wait()

// 	r := Re{
// 		Eat:     false,
// 		Eat_key: "_",
// 	}

// 	select {
// 	case eatKey := <-c1:
// 		if eatKey.eat {
// 			r = Re{
// 				Eat:     true,
// 				Eat_key: eatKey.key,
// 			}
// 		}
// 	}

// 	return r
// }

func (agar *AgarPosition) GetAgarSpace2(beads *beads.Beads, RoomId string) {
	// var dir []map[string]int = make([]map[string]int, 0)
	// for angle := 1; angle <= 360; angle++ {
	// 	for r := agar.Radius; r >= (agar.Radius - 10); r-- {
	// 		var x float64 = float64(r) * math.Sin(math.Pi*2*float64(angle)/360)
	// 		var y float64 = float64(r) * math.Cos(math.Pi*2*float64(angle)/360)
	// 		var xx int = int(agar.X + math.Round(float64(x*100))/100)
	// 		var yy int = int(agar.Y + math.Round(float64(y*100))/100)

	// 		myMap := make(map[string]int, 0)
	// 		myMap["x"] = xx
	// 		myMap["y"] = yy
	// 		dir = append(dir, myMap)
	// 	}
	// }

	// for _, v := range dir {
	// 	var x string = fmt.Sprintf("%v", v["x"])
	// 	var y string = fmt.Sprintf("%v", v["y"])

	// 	if (*beads)[RoomId][x+"_"+y] == 10 {
	// 		delete((*beads)[RoomId], x+"_"+y)
	// 		return Re{
	// 			Eat:     true,
	// 			Eat_key: x + "_" + y,
	// 		}
	// 	}
	// }

	// return Re{
	// 	Eat:     false,
	// 	Eat_key: "",
	// }
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
