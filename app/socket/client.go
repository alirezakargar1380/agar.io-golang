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
	beads := &beads.Beads{
		Beads: make(map[string]map[string]int),
	}
	if beads.Beads[c.RoomID] == nil {
		beads.Beads[c.RoomID] = make(map[string]int)
	}
	go func() {
		for {
			select {
			case <-ticker.C:
				if len(beads.Beads[c.RoomID]) == 200 {
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
					beads.Set(c.RoomID, key)
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
		c.sendResponse(beads, res.Command, res.Data)
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
			if d["opration"].(string) == "increse" {
				percent_of_speed := math.Round(float64(d["percent_of_speed"].(float64)))
				maxSpeed := GetMaxSpeedWithRadius(agarObject.Radius)
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
			} else {
				if Agars[c.RoomID][c.Client_id].Agars[i].Speed >= 0.10 {
					Agars[c.RoomID][c.Client_id].Agars[i].Speed -= 0.06
				} else {
					Agars[c.RoomID][c.Client_id].Agars[i].Speed = 0
				}
			}

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

				res["eat_key"] = eat.Eat_key
				if Agars[c.RoomID][c.Client_id].Agars[i].Radius < 450 {
					Agars[c.RoomID][c.Client_id].Agars[i].Radius += 5
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
			roomID: c.RoomID,
			Data:   js,
		}

	case "/game":
		fmt.Println("adding new agar")
		var lastId int = 0
		var newId int = 0
		var newRadius float64 = 0
		for i := 0; i < len(Agars[c.RoomID][c.Client_id].Agars); i++ {
			agar := Agars[c.RoomID][c.Client_id].Agars[i]
			lastId = agar.Id
		}
		newId = lastId + 1
		if newId == 1 {
			newRadius = 60
		} else {
			newRadius = 20
			Agars[c.RoomID][c.Client_id].Agars[0].Radius = Agars[c.RoomID][c.Client_id].Agars[0].Radius - 20
			Agars[c.RoomID][c.Client_id].Agars[0].Max_speed = GetMaxSpeedWithRadius(Agars[c.RoomID][c.Client_id].Agars[0].Radius)
		}

		Agars[c.RoomID][c.Client_id].Agars = append(Agars[c.RoomID][c.Client_id].Agars, trigonometric_circle.AgarDe{
			Id:        newId,
			X:         300,
			Y:         300,
			Radius:    newRadius,
			Max_speed: GetMaxSpeedWithRadius(newRadius),
			Speed:     0,
		})

		var new_agar_response map[string]string = make(map[string]string)
		new_agar_response["Command"] = "/new_agar"
		new_agar_response["x"] = fmt.Sprintf("%v", 300)
		new_agar_response["y"] = fmt.Sprintf("%v", 300)
		new_agar_response["radius"] = fmt.Sprintf("%v", newRadius)
		new_agar_response["id"] = fmt.Sprintf("%v", newId)

		reeee, _ := json.Marshal(new_agar_response)
		c.Hub.Broadcast <- &Message{
			roomID: c.RoomID,
			Data:   []byte(reeee),
		}
	case "/game_detail":
		var res map[string]string = make(map[string]string)
		res["Command"] = "/game_details"
		fmt.Println(Agars[c.RoomID])
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
	default:
		fmt.Println("def")
	}

}
