package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) CreateRoom(ctx *gin.Context) {
	var req CreateRoomReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	ctx.JSON(http.StatusOK, req)

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomId := ctx.Param("roomId")
	userId := ctx.Query("userId")
	username := ctx.Query("username")

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message),
		RoomID:   roomId,
		ID:       userId,
		Username: username,
	}

	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomId,
		Username: username,
	}

	// Register new client through the register channel
	h.hub.Register <- cl

	// Broadcast the message
	h.hub.Broadcast <- m
	// Write message
	go cl.writeMessage()
	// Read message
	cl.readMessage(h.hub)
}

func (h *Handler) GetRooms(ctx *gin.Context) {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	ctx.JSON(http.StatusOK, rooms)
}

func (h *Handler) GetClient(ctx *gin.Context) {
	var clients []ClientRes
	roomId := ctx.Param("roomId")

	if _, ok := h.hub.Rooms[roomId]; ok {
		clients = make([]ClientRes, 0)
		ctx.JSON(http.StatusOK, clients)
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	ctx.JSON(http.StatusOK, clients)
}
