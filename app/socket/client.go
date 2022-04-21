package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/alirezakargar1380/agar.io-golang/app/agar"
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

var beads map[string]int = make(map[string]int)

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
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if len(beads) == 12 {
					return
				}
				rand.Seed(time.Now().UnixNano())
				min := 1
				max := 500
				x := rand.Intn(max-min+1) + min
				// y := rand.Intn(max-min+1) + min
				// fmt.Println(x)
				var p map[string]string = make(map[string]string)
				p["Command"] = "/new_bead"
				p["x"] = fmt.Sprintf("%v", x)
				p["y"] = fmt.Sprintf("%v", 300)
				key := p["x"] + "_" + p["y"]
				// ms := strings.Split(key, "_")
				// fmt.Println(ms)
				beads[key] = 10
				// fmt.Println(key)
				json, _ := json.Marshal(p)
				c.Hub.Broadcast <- &Message{
					roomID: c.RoomID,
					Data:   []byte(json),
				}
			case <-quit:
				fmt.Println("stoped", c.RoomID)
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

			// Add queued chat messages to the current websocket message.
			// n := len(c.Send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	w.Write(<-c.Send)
			// }

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
type AA struct {
	Name string
}

type AgarDetail struct {
	// X    float64
	// Y    float64
	Size float32
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
		aga := data.(map[string]interface{})
		// fmt.Println(aga["X"].(float64))

		Agars[c.Client_id] = &AgarDetail{
			Size: Agars[c.Client_id].Size,
		}

		dir := &agar.AgarPosition{
			X: aga["X"].(float64),
			Y: aga["Y"].(float64),
		}
		directions := dir.GetAgarSpace()
		var eatIt bool = agar.CheckAgarSpace(directions, &beads)
		var res map[string]string = make(map[string]string)
		if eatIt {
			// fmt.Println(Agars[c.Client_id].Size)
			Agars[c.Client_id] = &AgarDetail{
				// X:    aga["X"].(float64),
				// Y:    aga["Y"].(float64),
				Size: Agars[c.Client_id].Size + 0.1,
			}
			// fmt.Println(Agars[c.Client_id].Size)
			res["size"] = fmt.Sprintf("%v", Agars[c.Client_id].Size)
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

		break
	case "/move":
		fmt.Println("/move...")
		// d := data.(map[string]interface{})
		// var a AA = AA{
		// 	Name: d["Name"].(string),
		// }
		// fmt.Println(a.Name)
		break
	}
}
