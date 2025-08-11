package main

import (
	"go-chat-server/db"
	"go-chat-server/internal/user"
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

	router.InitRouter(userHandler)
	router.Start("0.0.0.0:8080")
}
