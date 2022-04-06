package socket

type Hub struct {
	Clients    map[*Client]bool
	Rooms      map[string]map[*Client]bool
	Broadcast  chan *Message
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

// var Rooms map[string]map[*Client]bool = make(map[string]map[*Client]bool)

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			// roomTest := h.Rooms[client.RoomID]
			if h.Rooms == nil {
				h.Rooms = make(map[string]map[*Client]bool)
			}
			if h.Rooms[client.RoomID] == nil {
				r := make(map[*Client]bool)
				h.Rooms[client.RoomID] = r
			}
			h.Rooms[client.RoomID][client] = true

		case <-h.Unregister:
		case message := <-h.Broadcast:
			room := h.Rooms[message.roomID]
			for Client := range room {
				select {
				case Client.Send <- message.Data:
				default:
					close(Client.Send)
					delete(room, Client)
				}
			}
			if len(room) == 0 {
				delete(h.Rooms, message.roomID)
			}
		}
	}
}
