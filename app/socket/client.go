package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/alirezakargar1380/agar.io-golang/app/agar"
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

var Agars map[int64]*AgarDetail = make(map[int64]*AgarDetail)

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

var beads map[string]map[string]int = make(map[string]map[string]int)

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
	go func() {
		for {
			select {
			case <-ticker.C:
				if len(beads[c.RoomID]) == 12 {
					// fmt.Println("---> beads are full")
				} else {
					rand.Seed(time.Now().UnixNano())
					min := 10
					max := 500
					x := rand.Intn(max-min+1) + min
					y := rand.Intn(max-min+1) + min
					var p map[string]string = make(map[string]string)
					p["Command"] = "/new_bead"
					p["x"] = fmt.Sprintf("%v", x)
					p["y"] = fmt.Sprintf("%v", y)
					key := p["x"] + "_" + p["y"]
					if beads[c.RoomID] == nil {
						beads[c.RoomID] = make(map[string]int)
					}
					beads[c.RoomID][key] = 10
					json, _ := json.Marshal(p)
					c.Hub.Broadcast <- &Message{
						roomID: c.RoomID,
						Data:   []byte(json),
					}
				}

			case <-quit:
				fmt.Println("stoped", c.RoomID)
				delete(beads, c.RoomID)
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
		c.sendResponse(res.Command, res.Data)
	}
	// for range time.Tick(time.Second * 5) {
	// fmt.Println("Foo", c.RoomID)
	// c.Hub.Broadcast <- &Message{
	// 	roomID: c.RoomID,
	// 	Data:   []byte(`{"Command": "/new_bead", "x": "300", "y": "300"}`),
	// }
	// }
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
	X      float64
	Y      float64
	Size   float32
	Speed  float32
	Radius float32
}

func (c *Client) sendResponse(command interface{}, data interface{}) {
	// this is a test for sending message to client
	// c.Hub.Broadcast <- &Message{
	// 	roomID: c.RoomID,
	// 	Data:   []byte(`{"Command": "/new_agar", "data": ""}`),
	// }
	// for {
	// 	fmt.Println("sending...")
	// 	c.Hub.Broadcast <- &Message{
	// 		roomID: c.RoomID,
	// 		Data:   []byte(`{"Command": "/new_agar", "data": ""}`),
	// 	}
	// }
	// fmt.Println(command)
	// return

	switch command {
	case "/hello":
		fmt.Println("hello")
		aga := data.(map[string]interface{})
		// fmt.Println(aga["X"].(float64))

		Agars[c.Client_id] = &AgarDetail{
			Size: Agars[c.Client_id].Size,
		}

		dir := &agar.AgarPosition{
			X: aga["X"].(float64),
			Y: aga["Y"].(float64),
		}
		// directions := dir.GetAgarSpace()
		// var eatIt bool = agar.CheckAgarSpace(directions, &beads)
		var res map[string]string = make(map[string]string)
		// if eatIt {
		// 	Agars[c.Client_id] = &AgarDetail{
		// 		Size: Agars[c.Client_id].Size + 0.1,
		// 	}
		// 	// fmt.Println(Agars[c.Client_id].Size)
		// 	res["size"] = fmt.Sprintf("%v", Agars[c.Client_id].Size)
		// }
		var eat bool = dir.GetAgarSpace2(&beads, c.RoomID)
		if eat {
			fmt.Println("eat")
		}

		res["Command"] = "/new_agar"
		res["x"] = fmt.Sprintf("%v", aga["X"].(float64))
		res["y"] = fmt.Sprintf("%v", aga["Y"].(float64))

		js, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
			return
		}

		var resp []byte = make([]byte, 0)
		resp = append(resp, js...)
		c.Hub.Broadcast <- &Message{
			roomID: c.RoomID,
			Data:   resp,
		}
	case "/move":
		// fmt.Println(Agars[c.Client_id].Speed)
		d := data.(map[string]interface{})
		if d["opration"].(string) == "increse" {
			if Agars[c.Client_id].Speed <= 5 {
				Agars[c.Client_id] = &AgarDetail{
					Speed: Agars[c.Client_id].Speed + 0.1,
					X:     Agars[c.Client_id].X,
					Y:     Agars[c.Client_id].Y,
				}
			}
		} else {
			if Agars[c.Client_id].Speed >= 0.10 {
				Agars[c.Client_id] = &AgarDetail{
					Speed: Agars[c.Client_id].Speed - 0.06,
					X:     Agars[c.Client_id].X,
					Y:     Agars[c.Client_id].Y,
				}
			}
		}
		tri := &trigonometric_circle.AgarDetail{
			X:      Agars[c.Client_id].X,
			Y:      Agars[c.Client_id].Y,
			Radius: float64(Agars[c.Client_id].Speed),
		}
		directions := tri.Test(d["angle"].(float64))

		Agars[c.Client_id] = &AgarDetail{
			X:     directions["x"],
			Y:     directions["y"],
			Speed: Agars[c.Client_id].Speed,
		}

		var res map[string]string = make(map[string]string)
		res["Command"] = "/m_agar"
		res["x"] = fmt.Sprintf("%v", directions["x"])
		res["y"] = fmt.Sprintf("%v", directions["y"])
		res["speed"] = fmt.Sprintf("%v", Agars[c.Client_id].Speed)

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
