package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/alirezakargar1380/agar.io-golang/app/agar"
	"github.com/alirezakargar1380/agar.io-golang/app/beads"
	redis_db "github.com/alirezakargar1380/agar.io-golang/app/service"
	"github.com/alirezakargar1380/agar.io-golang/app/trigonometric_circle"
	agar_arrays "github.com/alirezakargar1380/agar.io-golang/app/utils"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var Agars map[string]map[int64]*AgarDetail = make(map[string]map[int64]*AgarDetail)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{}

type Client struct {
	RoomID    string
	Client_id int64
	Hub       *Hub
	Conn      *websocket.Conn
	Send      chan []byte
	Color     string
	Loose     bool
}

type Message struct {
	RoomID string
	Data   []byte
}

type Data struct {
	Command string
	Data    interface{}
}

var Gamebeads *beads.Beads = &beads.Beads{
	Beads: make(map[string]map[string]int),
}

func (c *Client) ReadPump() {
	quit := make(chan struct{})
	defer func() {
		fmt.Println("Client disconnected...")
		quit <- struct{}{}
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	ticker := time.NewTicker(1 * time.Second)
	if Gamebeads.Beads[c.RoomID] == nil {
		fmt.Println("Beads is nil")
		Gamebeads.Beads[c.RoomID] = make(map[string]int)
	}
	go func() {
		for {
			select {
			case <-ticker.C:
				if redis_db.Client.CountStars(c.RoomID) >= 110 {
					continue
				} else {
					min := 500
					max := 1000
					x := rand.Intn(max-min+1) + min
					y := rand.Intn(max-min+1) + min
					var p map[string]string = make(map[string]string)
					p["Command"] = "/new_bead"
					p["x"] = fmt.Sprintf("%v", x)
					p["y"] = fmt.Sprintf("%v", y)
					key := p["x"] + "_" + p["y"]
					redis_db.Client.AddStar(key, c.RoomID)
					// Gamebeads.Set(c.RoomID, key)
					json, _ := json.Marshal(p)
					c.Hub.Broadcast <- &Message{
						RoomID: c.RoomID,
						Data:   []byte(json),
					}
				}
			case <-quit:
				fmt.Println("stoped", c.RoomID)
				// delete(beads, c.RoomID)
				ticker.Stop()
				return
			}
		}
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var res Data
		json.Unmarshal([]byte(message), &res)
		if len(Agars[c.RoomID][c.Client_id].Agars) == 0 {
			if !c.Loose {
				c.Loose = true
				var res map[string]string = make(map[string]string)
				res["Command"] = "/finish_game"
				res["LooserId"] = fmt.Sprintf("%v", c.Client_id)
				jData, _ := json.Marshal(res)
				c.Hub.Broadcast <- &Message{
					RoomID: c.RoomID,
					Data:   jData,
				}
				fmt.Println("you have been loose", c.Client_id)
			}
			continue
		}
		c.sendResponse(Gamebeads, res.Command, res.Data)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type Agar struct {
	X float64
	Y float64
}

type AgarDetail struct {
	Agars     []trigonometric_circle.AgarDe
	Client_id int
	Color     string
}

func GetMaxSpeedWithRadius(Radius float64) float64 {
	speed := 7 - (Radius * 0.013)
	return math.Floor(speed*1000) / 1000
}

func (c *Client) sendResponse(beads *beads.Beads, command interface{}, data interface{}) {
	switch command {
	case "/create_battle":
		var res map[string]string = make(map[string]string)
		res["Command"] = "/start_game"
		jData, _ := json.Marshal(res)
		c.Hub.Broadcast <- &Message{
			RoomID: c.RoomID,
			Data:   jData,
		}
	case "/move":
		// fmt.Println("move...")
		var res map[string]string = make(map[string]string)
		d := data.(map[string]interface{})
		for i := 0; i < len(Agars[c.RoomID][c.Client_id].Agars); i++ {
			agarObject := Agars[c.RoomID][c.Client_id].Agars[i]
			if agarObject.Lock {
				continue
			}
			// if d["opration"].(string) == "increse" {
			percent_of_speed := math.Round(float64(d["percent_of_speed"].(float64)))

			var speed float64
			if i == 0 {
				speed = Agars[c.RoomID][c.Client_id].Agars[i].Radius
			} else {
				speed = Agars[c.RoomID][c.Client_id].Agars[0].Radius - 1
			}

			maxSpeed := GetMaxSpeedWithRadius(speed)
			var speedThatUserWant float64 = float64(percent_of_speed*100) * float64(maxSpeed) / 100
			speedThatUserWant = speedThatUserWant / 100
			speedThatUserWant = math.Floor(speedThatUserWant*100) / 100
			// fmt.Println("id", agarObject.Id, "speed", speedThatUserWant, "now speed ", Agars[c.RoomID][c.Client_id].Agars[i].Speed, "max speed", maxSpeed)

			if Agars[c.RoomID][c.Client_id].Agars[i].Speed != maxSpeed {
				if maxSpeed >= speedThatUserWant {
					if speedThatUserWant > Agars[c.RoomID][c.Client_id].Agars[i].Speed {
						Agars[c.RoomID][c.Client_id].Agars[i].Speed += 0.1
						Agars[c.RoomID][c.Client_id].Agars[i].Speed = math.Floor(Agars[c.RoomID][c.Client_id].Agars[i].Speed*100) / 100
					} else {
						Agars[c.RoomID][c.Client_id].Agars[i].Speed -= 0.1
						Agars[c.RoomID][c.Client_id].Agars[i].Speed = math.Floor(Agars[c.RoomID][c.Client_id].Agars[i].Speed*100) / 100
					}
				}
			}

			// get new movement (x,y)
			tri := &trigonometric_circle.AgarDetail{
				Id:     agarObject.Id,
				X:      agarObject.X,
				Y:      agarObject.Y,
				Speed:  float64(Agars[c.RoomID][c.Client_id].Agars[i].Speed),
				Radius: float64(agarObject.Radius),
			}
			directions := tri.Test(d["angle"].(float64))

			if agarObject.Id != 1 {
				tri.CheckForEatTogether(Agars[c.RoomID][c.Client_id].Agars)
			}

			if directions["y"]-agarObject.Radius < 0 {
				directions["y"] = agarObject.Radius
			}
			if directions["y"] > 3000 {
				directions["y"] = 3000
			}

			if directions["x"]-agarObject.Radius < 0 {
				directions["x"] = agarObject.Radius
			}
			if directions["x"] > 3000 {
				directions["x"] = 3000
			}

			Agars[c.RoomID][c.Client_id].Agars[i].X = directions["x"]
			Agars[c.RoomID][c.Client_id].Agars[i].Y = directions["y"]

			dir := &agar.AgarPosition{
				X:      directions["x"],
				Y:      directions["y"],
				Radius: int(Agars[c.RoomID][c.Client_id].Agars[i].Radius),
			}

			sss := redis_db.Client.GetStars(c.RoomID)
			// fmt.Println(sss)
			// check agar eat bead
			eat := dir.GetAgarSpace5(sss, c.RoomID)
			if eat.Eat {
				// fmt.Println("eat", eat.Eat_key)
				eatKeys, err := json.Marshal(eat.Eat_key)
				if err != nil {
					return
				}
				res["eat_key"] = string(eatKeys)
				redis_db.Client.DeleteStart(c.RoomID, eat.Eat_key)
				if Agars[c.RoomID][c.Client_id].Agars[i].Radius < 450 {
					Agars[c.RoomID][c.Client_id].Agars[i].Radius += 2
				}
			}

			// check for user agars if they eat together
			checkAgars := agar.AllAgars{
				Agars:  Agars[c.RoomID][c.Client_id].Agars,
				Id:     Agars[c.RoomID][c.Client_id].Agars[i].Id,
				X:      Agars[c.RoomID][c.Client_id].Agars[i].X,
				Y:      Agars[c.RoomID][c.Client_id].Agars[i].Y,
				Radius: int(Agars[c.RoomID][c.Client_id].Agars[i].Radius),
			}
			eatTogetherResult := checkAgars.CheckForEating()
			if eatTogetherResult.Status {
				agarsArrayHandler := &agar_arrays.Agars{
					Agars: Agars[c.RoomID][c.Client_id].Agars,
				}
				Eated_agar_by_index := agarsArrayHandler.GETAgarIndexWithId(eatTogetherResult.Eated_agar_by_id)
				Eated_agar_index := agarsArrayHandler.GETAgarIndexWithId(eatTogetherResult.Eated_agar_id)
				Agars[c.RoomID][c.Client_id].Agars[Eated_agar_by_index].Radius += Agars[c.RoomID][c.Client_id].Agars[Eated_agar_index].Radius
				fmt.Println("REmove an agar", eatTogetherResult.Eated_agar_by_id, eatTogetherResult.Eated_agar_id)
				Agars[c.RoomID][c.Client_id].Agars = agarsArrayHandler.RemoveAgarFromArrayWithIndex(Eated_agar_index)
			}

			// check for other user agars if they eat together
			for _, v := range Agars[c.RoomID] {
				if int64(v.Client_id) != c.Client_id {
					if len(Agars[c.RoomID][int64(v.Client_id)].Agars) == 0 {
						delete(Agars[c.RoomID], int64(v.Client_id))
						continue
					}
					checkForOtherAgars := agar.AllAgars{
						ClientId: int(c.Client_id),
						RivalId:  v.Client_id,
						Agars:    v.Agars,
						Id:       Agars[c.RoomID][c.Client_id].Agars[i].Id,
						X:        Agars[c.RoomID][c.Client_id].Agars[i].X,
						Y:        Agars[c.RoomID][c.Client_id].Agars[i].Y,
						Radius:   int(Agars[c.RoomID][c.Client_id].Agars[i].Radius),
					}
					res := checkForOtherAgars.CheckForAgarEatingOtherAgars()
					if res.Status {
						EatAgarsArrayHandler := &agar_arrays.Agars{
							Agars: Agars[c.RoomID][int64(res.EatClientId)].Agars,
						}
						EatenAgarsArrayHandler := &agar_arrays.Agars{
							Agars: Agars[c.RoomID][int64(res.EatenClientId)].Agars,
						}
						EatAgarIndex := EatAgarsArrayHandler.GETAgarIndexWithId(res.EatAgarId)
						EatenAgarIndex := EatenAgarsArrayHandler.GETAgarIndexWithId(res.EatenAgarId)
						Agars[c.RoomID][int64(res.EatClientId)].Agars[EatAgarIndex].Radius += Agars[c.RoomID][int64(res.EatenClientId)].Agars[EatenAgarIndex].Radius

						if res.EatenAgarId == 1 {
							fmt.Println("----------------------------------------------> ")
							Agars[c.RoomID][int64(res.EatenClientId)].Agars = make([]trigonometric_circle.AgarDe, 0)
							var resp map[string]string = make(map[string]string)
							resp["Command"] = "/finish_game"
							resp["LooserId"] = fmt.Sprintf("%v", int(res.EatenClientId))
							jData, _ := json.Marshal(resp)
							c.Hub.Broadcast <- &Message{
								RoomID: c.RoomID,
								Data:   jData,
							}
						} else {
							agarsArrayHandler := &agar_arrays.Agars{
								Agars: Agars[c.RoomID][int64(res.EatenClientId)].Agars,
							}
							EatenAgarIndex := agarsArrayHandler.GETAgarIndexWithId(res.EatenAgarId)
							Agars[c.RoomID][int64(res.EatenClientId)].Agars = agarsArrayHandler.RemoveAgarFromArrayWithIndex(EatenAgarIndex)
						}
					}
				}
			}
		}

		res["Command"] = "/move_agars"
		dd, err := json.Marshal(Agars[c.RoomID])
		if err != nil {
			fmt.Println(err)
			return
		}
		res["agars"] = string(dd)
		js, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
			return
		}

		c.Hub.Broadcast <- &Message{
			RoomID: c.RoomID,
			Data:   js,
		}
	case "/halfagar":
		d := data.(map[string]interface{})
		var lastId int = 0
		var newId int = 0
		var newRadius float64 = 0
		for i := 0; i < len(Agars[c.RoomID][c.Client_id].Agars); i++ {
			agar := Agars[c.RoomID][c.Client_id].Agars[i]
			lastId = agar.Id
		}
		newId = lastId + 1
		newRadius = 20

		Agars[c.RoomID][c.Client_id].Agars[0].Radius = Agars[c.RoomID][c.Client_id].Agars[0].Radius - 20
		Agars[c.RoomID][c.Client_id].Agars[0].Max_speed = GetMaxSpeedWithRadius(Agars[c.RoomID][c.Client_id].Agars[0].Radius)

		Agars[c.RoomID][c.Client_id].Agars = append(Agars[c.RoomID][c.Client_id].Agars, trigonometric_circle.AgarDe{
			Id:        newId,
			X:         Agars[c.RoomID][c.Client_id].Agars[0].X,
			Y:         Agars[c.RoomID][c.Client_id].Agars[0].Y,
			Radius:    newRadius,
			Max_speed: GetMaxSpeedWithRadius(newRadius),
			Speed:     2.5,
		})

		lastAgarKey := len(Agars[c.RoomID][c.Client_id].Agars) - 1
		lastAgar := Agars[c.RoomID][c.Client_id].Agars[lastAgarKey]
		var movement_res map[string]string = make(map[string]string)

		tri := &trigonometric_circle.AgarDetail{
			Id:    lastAgar.Id,
			X:     Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].X,
			Y:     Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].Y,
			Speed: float64(Agars[c.RoomID][c.Client_id].Agars[0].Radius) + (4 * 20),
		}
		directions := tri.Test(d["angle"].(float64))

		var new_agar_response map[string]string = make(map[string]string)
		new_agar_response["Command"] = "/new_agar"
		new_agar_response["start_x"] = fmt.Sprintf("%v", Agars[c.RoomID][c.Client_id].Agars[0].X)
		new_agar_response["start_y"] = fmt.Sprintf("%v", Agars[c.RoomID][c.Client_id].Agars[0].Y)
		new_agar_response["x"] = fmt.Sprintf("%v", directions["x"])
		new_agar_response["y"] = fmt.Sprintf("%v", directions["y"])
		new_agar_response["radius"] = fmt.Sprintf("%v", newRadius)
		new_agar_response["id"] = fmt.Sprintf("%v", newId)
		new_agar_response["color"] = fmt.Sprintf("%v", c.Color)
		new_agar_response["client_id"] = fmt.Sprintf("%v", c.Client_id)

		reeee, _ := json.Marshal(new_agar_response)
		c.Hub.Broadcast <- &Message{
			RoomID: c.RoomID,
			Data:   []byte(reeee),
		}

		quit := make(chan struct{})
		var i int = 0
		ticker := time.NewTicker(100 * time.Millisecond)
		eat_keys := make([]string, 0)
		go func() {
			for {
				select {
				case <-ticker.C:
					i++
					if i > 5 {
						quit <- struct{}{}
					} else {
						agarsNum := len(Agars[c.RoomID][c.Client_id].Agars) - 1
						if lastAgarKey > agarsNum {
							continue
						}

						var speed float64 = 0
						if i == 1 {
							speed = Agars[c.RoomID][c.Client_id].Agars[0].Radius
						} else {
							speed = 20
						}

						tri := &trigonometric_circle.AgarDetail{
							Id:    lastAgar.Id,
							X:     Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].X,
							Y:     Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].Y,
							Speed: float64(speed),
						}
						directions := tri.Test(d["angle"].(float64))

						dir := &agar.AgarPosition{
							X:      directions["x"],
							Y:      directions["y"],
							Radius: int(Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].Radius),
						}
						eat := dir.GetAgarSpace4(beads, c.RoomID)
						if eat.Eat {
							for i := 0; i < len(eat.Eat_key); i++ {
								eat_keys = append(eat_keys, eat.Eat_key[i])
								Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].Radius += 2
							}

							// sending response
							movement_res["Command"] = "/eated_agars_keys"

							dd, err := json.Marshal(Agars[c.RoomID])
							if err != nil {
								fmt.Println(err)
								return
							}
							movement_res["agars"] = string(dd)

							eatKeys, err := json.Marshal(eat_keys)
							if err != nil {
								return
							}
							movement_res["eat_key"] = string(eatKeys)
							js, err := json.Marshal(movement_res)
							if err != nil {
								fmt.Println(err)
								return
							}

							c.Hub.Broadcast <- &Message{
								RoomID: c.RoomID,
								Data:   js,
							}
						}

						Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].X = directions["x"]
						Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].Y = directions["y"]

						// all up if should check here
						// we should min the agars radius if there are more than a pecified number
						for i := 0; i < len(Agars[c.RoomID][c.Client_id].Agars); i++ {
							agarObject := Agars[c.RoomID][c.Client_id].Agars[i]
							if agarObject.Id == 1 {
								continue
							}
							checkAgars := agar.AllAgars{
								Agars:  Agars[c.RoomID][c.Client_id].Agars,
								Id:     agarObject.Id,
								X:      agarObject.X,
								Y:      agarObject.Y,
								Radius: int(agarObject.Radius),
							}
							eatTogetherResult := checkAgars.CheckForEatingWhenT()
							if eatTogetherResult.Status {
								agarsArrayHandler := &agar_arrays.Agars{
									Agars: Agars[c.RoomID][c.Client_id].Agars,
								}
								Eated_agar_by_index := agarsArrayHandler.GETAgarIndexWithId(eatTogetherResult.Eated_agar_by_id)
								Eated_agar_index := agarsArrayHandler.GETAgarIndexWithId(eatTogetherResult.Eated_agar_id)
								Agars[c.RoomID][c.Client_id].Agars[Eated_agar_by_index].Radius += Agars[c.RoomID][c.Client_id].Agars[Eated_agar_index].Radius
								fmt.Println("eaten an agar", eatTogetherResult.Eated_agar_by_id, eatTogetherResult.Eated_agar_id)
								Agars[c.RoomID][c.Client_id].Agars = agarsArrayHandler.RemoveAgarFromArrayWithIndex(Eated_agar_index)

								movement_res["Command"] = "/move_agars"
								dd, err := json.Marshal(Agars[c.RoomID])
								if err != nil {
									fmt.Println(err)
									return
								}
								movement_res["agars"] = string(dd)
								js, err := json.Marshal(movement_res)
								if err != nil {
									fmt.Println(err)
									return
								}

								c.Hub.Broadcast <- &Message{
									RoomID: c.RoomID,
									Data:   js,
								}
							}

							for _, v := range Agars[c.RoomID] {
								if int64(v.Client_id) != c.Client_id {
									checkForOtherAgars := agar.AllAgars{
										ClientId: int(c.Client_id),
										RivalId:  v.Client_id,
										Agars:    v.Agars,
										Id:       Agars[c.RoomID][c.Client_id].Agars[i].Id,
										X:        Agars[c.RoomID][c.Client_id].Agars[i].X,
										Y:        Agars[c.RoomID][c.Client_id].Agars[i].Y,
										Radius:   int(Agars[c.RoomID][c.Client_id].Agars[i].Radius),
									}
									res := checkForOtherAgars.CheckForAgarEatingOtherAgars()
									if res.Status {
										if res.EatenAgarId == 1 {
											Agars[c.RoomID][int64(res.EatenClientId)].Agars = make([]trigonometric_circle.AgarDe, 0)
										} else {
											agarsArrayHandler := &agar_arrays.Agars{
												Agars: Agars[c.RoomID][int64(res.EatenClientId)].Agars,
											}
											EatenAgarIndex := agarsArrayHandler.GETAgarIndexWithId(res.EatenAgarId)
											Agars[c.RoomID][int64(res.EatenClientId)].Agars = agarsArrayHandler.RemoveAgarFromArrayWithIndex(EatenAgarIndex)
										}
									}
									movement_res["Command"] = "/move_agars"
									dd, err := json.Marshal(Agars[c.RoomID])
									if err != nil {
										fmt.Println(err)
										return
									}
									movement_res["agars"] = string(dd)
									js, err := json.Marshal(movement_res)
									if err != nil {
										fmt.Println(err)
										return
									}

									c.Hub.Broadcast <- &Message{
										RoomID: c.RoomID,
										Data:   js,
									}
								}
							}
						}

					}
				case <-quit:
					fmt.Println("quit...")
					ticker.Stop()
				}
			}
		}()

	case "/game_detail":
		var res map[string]string = make(map[string]string)
		res["Command"] = "/game_details"
		dd, err := json.Marshal(Agars[c.RoomID])
		if err != nil {
			fmt.Println(err)
			return
		}
		sss := redis_db.Client.GetStars(c.RoomID)
		beadss, err := json.Marshal(sss)
		if err != nil {
			fmt.Println(err)
			return
		}

		res["agars"] = string(dd)
		res["beads"] = string(beadss)
		res["color"] = c.Color
		js, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
			return
		}
		c.Hub.Broadcast <- &Message{
			RoomID: c.RoomID,
			Data:   js,
		}
	default:
		fmt.Println("def")
	}

}
