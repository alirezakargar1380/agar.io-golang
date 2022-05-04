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
	"github.com/alirezakargar1380/agar.io-golang/app/trigonometric_circle"
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
}

type Message struct {
	roomID string
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
				if len(Gamebeads.Beads[c.RoomID]) == 200 {
					fmt.Println("---> beads are full")
				} else {
					// 	rand.Seed(time.Now().UnixNano())
					min := 500
					max := 1000
					x := rand.Intn(max-min+1) + min
					y := rand.Intn(max-min+1) + min
					var p map[string]string = make(map[string]string)
					p["Command"] = "/new_bead"
					p["x"] = fmt.Sprintf("%v", x)
					p["y"] = fmt.Sprintf("%v", y)
					key := p["x"] + "_" + p["y"]
					Gamebeads.Set(c.RoomID, key)
					json, _ := json.Marshal(p)
					c.Hub.Broadcast <- &Message{
						roomID: c.RoomID,
						Data:   []byte(json),
					}
				}
				// fmt.Println(len(beads.Beads[c.RoomID]))
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
}

func GetMaxSpeedWithRadius(Radius float64) float64 {
	speed := 7 - (Radius * 0.013)
	return math.Floor(speed*1000) / 1000
}

func (c *Client) sendResponse(beads *beads.Beads, command interface{}, data interface{}) {
	switch command {
	case "/move":
		var res map[string]string = make(map[string]string)
		d := data.(map[string]interface{})
		for i := 0; i < len(Agars[c.RoomID][c.Client_id].Agars); i++ {
			agarObject := Agars[c.RoomID][c.Client_id].Agars[i]
			// if d["opration"].(string) == "increse" {
			percent_of_speed := math.Round(float64(d["percent_of_speed"].(float64)))
			maxSpeed := GetMaxSpeedWithRadius(Agars[c.RoomID][c.Client_id].Agars[i].Radius)
			var dd float64 = float64(percent_of_speed*100) * float64(maxSpeed) / 100
			dd = dd / 100
			dd = math.Floor(dd*100) / 100
			if dd == maxSpeed || (dd+0.01) == maxSpeed {
				if maxSpeed > agarObject.Speed {
					Agars[c.RoomID][c.Client_id].Agars[i].Speed += 0.1
				}
			} else {
				if Agars[c.RoomID][c.Client_id].Agars[i].Speed > 0 {
					if Agars[c.RoomID][c.Client_id].Agars[i].Speed > dd {
						Agars[c.RoomID][c.Client_id].Agars[i].Speed -= 0.1
					}
				}
			}
			// } else {
			// 	if Agars[c.RoomID][c.Client_id].Agars[i].Speed >= 0.10 {
			// 		Agars[c.RoomID][c.Client_id].Agars[i].Speed -= 0.06
			// 	} else {
			// 		Agars[c.RoomID][c.Client_id].Agars[i].Speed = 0
			// 	}
			// }

			tri := &trigonometric_circle.AgarDetail{
				Id:     agarObject.Id,
				X:      agarObject.X,
				Y:      agarObject.Y,
				Speed:  float64(agarObject.Speed),
				Radius: float64(agarObject.Radius),
			}
			directions := tri.Test(d["angle"].(float64))

			if agarObject.Id != 1 {
				tri.CheckForEatTogether(Agars[c.RoomID][c.Client_id].Agars)
			}

			// fmt.Println(directions["x"], directions["y"])

			Agars[c.RoomID][c.Client_id].Agars[i].X = directions["x"]
			Agars[c.RoomID][c.Client_id].Agars[i].Y = directions["y"]

			dir := &agar.AgarPosition{
				X:      directions["x"],
				Y:      directions["y"],
				Radius: int(Agars[c.RoomID][c.Client_id].Agars[i].Radius),
			}

			eat := dir.GetAgarSpace4(beads, c.RoomID)

			if eat.Eat {
				eatKeys, err := json.Marshal(eat.Eat_key)
				if err != nil {
					return
				}
				res["eat_key"] = string(eatKeys)
				if Agars[c.RoomID][c.Client_id].Agars[i].Radius < 450 {
					Agars[c.RoomID][c.Client_id].Agars[i].Radius += 1
				}
				// delete(Gamebeads.Beads[c.RoomID], eat.Eat_key)
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
			roomID: c.RoomID,
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
		fmt.Println("adding new agar id:", newId)
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
			Speed: float64(Agars[c.RoomID][c.Client_id].Agars[0].Radius * 5),
		}
		directions := tri.Test(d["angle"].(float64))

		var new_agar_response map[string]string = make(map[string]string)
		new_agar_response["Command"] = "/new_agar"
		new_agar_response["x"] = fmt.Sprintf("%v", directions["x"])
		new_agar_response["y"] = fmt.Sprintf("%v", directions["y"])
		new_agar_response["radius"] = fmt.Sprintf("%v", newRadius)
		new_agar_response["id"] = fmt.Sprintf("%v", newId)

		reeee, _ := json.Marshal(new_agar_response)
		c.Hub.Broadcast <- &Message{
			roomID: c.RoomID,
			Data:   []byte(reeee),
		}

		// process
		// d := trigonometric_circle.GetDistanceBetweenTowPoint(100, 100, 300, 300)

		// if Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].Speed > 0.1 {
		// 	fmt.Println("error", Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].Speed)
		// 	Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].Speed -= 0.01
		// }

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
						fmt.Println(i)
						tri := &trigonometric_circle.AgarDetail{
							Id:    lastAgar.Id,
							X:     Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].X,
							Y:     Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].Y,
							Speed: float64(Agars[c.RoomID][c.Client_id].Agars[0].Radius),
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
								fmt.Println("eat", eat.Eat_key[i])
								eat_keys = append(eat_keys, eat.Eat_key[i])
								Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].Radius += 20

								// sending response
								movement_res["Command"] = "/move_agars"
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
									roomID: c.RoomID,
									Data:   js,
								}
							}
						}

						// fmt.Println(directions["x"], directions["y"])
						Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].X = directions["x"]
						Agars[c.RoomID][c.Client_id].Agars[lastAgarKey].Y = directions["y"]

						// dd, err := json.Marshal(Agars[c.RoomID])
						// if err != nil {
						// 	fmt.Println(err)
						// 	return
						// }
						// movement_res["agars"] = string(dd)
						// js, err := json.Marshal(movement_res)
						// if err != nil {
						// 	fmt.Println(err)
						// 	return
						// }

						// c.Hub.Broadcast <- &Message{
						// 	roomID: c.RoomID,
						// 	Data:   js,
						// }
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
		beadss, err := json.Marshal(Gamebeads.Beads[c.RoomID])
		if err != nil {
			fmt.Println(err)
			return
		}
		res["agars"] = string(dd)
		res["beads"] = string(beadss)
		js, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
			return
		}
		c.Hub.Broadcast <- &Message{
			roomID: c.RoomID,
			Data:   js,
		}
	default:
		fmt.Println("def")
	}

}
