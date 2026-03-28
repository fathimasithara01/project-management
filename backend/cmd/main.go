package main

import (
	"log"
	"os"
	"project-management/internal/bootstrap"
)

func main() {
	router, err := bootstrap.NewRouter()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
