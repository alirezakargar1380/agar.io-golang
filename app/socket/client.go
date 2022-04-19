package socket

import (
	"encoding/json"
	"fmt"
	"log"
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
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

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
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
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
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

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

var agars map[int64]map[string]float64 = make(map[int64]map[string]float64)

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
		agars[c.Client_id] = map[string]float64{
			"X": aga["X"].(float64),
			"Y": aga["Y"].(float64),
		}
		var p map[string]string = make(map[string]string)
		p["Command"] = "/new_agar"
		p["x"] = fmt.Sprintf("%v", aga["X"].(float64))
		p["y"] = fmt.Sprintf("%v", aga["Y"].(float64))
		js, err := json.Marshal(p)
		if err != nil {
			fmt.Println(err)
			return
		}

		dir := &agar.AgarPosition{
			X: aga["X"].(float64),
			Y: aga["Y"].(float64),
		}
		directions := dir.GetAgarSpace()
		agar.CheckAgarSpace(directions)
		// directions := agar.GetAgarSpace(aga["X"].(float64), aga["Y"].(float64))
		// fmt.Println(directions[0]["x"])
		// fmt.Println(len(directions))
		var resp []byte = make([]byte, 0)
		resp = append(resp, js...)
		c.Hub.Broadcast <- &Message{
			roomID: c.RoomID,
			// Data:   []byte(`{"Command": "/new_agar", "x": "", "y": ""}`),
			Data: resp,
		}
		// agar := Agar{
		// 	X: aga["X"].(float64),
		// 	Y: aga["Y"].(float64),
		// }
		// fmt.Println(agar.Y - agar.X)
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
