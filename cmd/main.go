package main

import (
	"log"
	"sheeva/internal/vault"
	"sheeva/config"
	"sheeva/api"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("Error loading config: %w", err)
		return
	}

	vaultClient, err := vault.NewVaultClient(cfg)
	if err != nil {
		log.Fatalf("Error initializing Vault client: %v", err)
	}

	router := gin.Default()

	api.SetupRoutes(router, vaultClient)

	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}
