package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Anirudh-S-Kumar/disec/cmd/server/initializer"
	"github.com/Anirudh-S-Kumar/disec/server"
	"github.com/joho/godotenv"
)

func main() {
	// setup dotenv
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env file: %+v", err)
	}

	certPath, keyPath := initializer.CertPaths()

	router := server.NewRouter(false)
	if err := router.RunTLS(fmt.Sprintf(":%v", os.Getenv("PORT")), certPath, keyPath); err != nil {
		fmt.Println("Error starting server with TLS:", err)
	}
}
