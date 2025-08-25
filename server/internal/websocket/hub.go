package websocket

type Room struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Clients map[string]*Client
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

// hub.Run() ini mengatur pesan
func (h *Hub) Run() {
	for {
		select {
		// Memasukkan client baru ke room
		case cl := <-h.Register:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				r := h.Rooms[cl.RoomID]

				if _, ok := r.Clients[cl.ID]; !ok {
					r.Clients[cl.ID] = cl
				}
			}

		// Hapus client dari room
		case cl := <-h.Unregister:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := h.Rooms[cl.RoomID].Clients[cl.ID]; ok {
					if len(h.Rooms[cl.RoomID].Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  "User left the chat",
							RoomID:   cl.RoomID,
							Username: cl.Username,
						}
					}

					delete(h.Rooms[cl.RoomID].Clients, cl.ID)
					close(cl.Message)
				}
			}

		// Menyebarkan pesan ke semua client dalam room yang sudah dibuat
		case m := <-h.Broadcast:
			if _, ok := h.Rooms[m.RoomID]; ok {
				// Semua client di room akan mendapat pesan lewat channel cl.message
				// Goroutine writeMessage() akan ambil dari channel ini, lalu kirim via websocket
				for _, cl := range h.Rooms[m.RoomID].Clients {
					cl.Message <- m
				}
			}
		}
	}
}
