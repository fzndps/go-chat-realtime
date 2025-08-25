package main

import (
	"go-chat-server/db"
	"go-chat-server/internal/user"
	"go-chat-server/internal/websocket"
	"go-chat-server/router"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbConnection, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}

	userRepo := user.NewRepository(dbConnection.GetDB())
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc)

	// Membuat hub untuk pusat koordinasi semua room
	hub := websocket.NewHub()
	wsHandler := websocket.NewHandler(hub)

	// Jalan di go routine dengan terus listen channel Register, Unregister, dan Broadcast
	go hub.Run()

	router.InitRouter(userHandler, wsHandler)
	router.Start("0.0.0.0:8080")
}
