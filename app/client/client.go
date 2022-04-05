package client

import (
	"github.com/alirezakargar1380/agar.io-golang/app/socket/hub"
)

type Client struct {
	// roomID string
	hub hub.Hub
	// conn *websocket.Conn
	send chan []byte
}
