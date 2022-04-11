package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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
	RoomID string
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
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
		// fmt.Println(string([]byte(message)))
		var res Data
		json.Unmarshal([]byte(message), &res)

		// var str_message string = string([]byte(message))
		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.sendResponse(res.Command, res.Data)
		// c.Hub.Broadcast <- &Message{
		// 	roomID: c.RoomID,
		// 	Data:   message,
		// }
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

func (c *Client) sendResponse(command interface{}, data interface{}) {
	// this is a test for sending message to client
	// c.Hub.Broadcast <- &Message{
	// 	roomID: c.RoomID,
	// 	Data:   []byte("test"),
	// }
	// fmt.Println(data)
	// return
	switch command {
	case "/hello":
		// var a Agar = data
		// data := data
		// agar := Agar{
		// 	X: data["X"],
		// 	Y: data["Y"],
		// }
		aga := data.(map[string]interface{})
		agar := Agar{
			X: aga["X"].(float64),
			Y: aga["Y"].(float64),
		}
		// fmt.Println("/hello...")
		fmt.Println(agar.Y - agar.X)
		break
	case "/move":
		fmt.Println("/move...")
		d := data.(map[string]interface{})
		var a AA = AA{
			Name: d["Name"].(string),
		}
		fmt.Println(a.Name)
		break
	}
}
