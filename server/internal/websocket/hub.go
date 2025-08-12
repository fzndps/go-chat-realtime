package websocket

type Room struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Clients map[string]*Client
}

type Hub struct {
	Rooms map[string]*Room
}

func NewHUb() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}
