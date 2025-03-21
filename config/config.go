package config

import (
	"log"
	"os" 
	"github.com/joho/godotenv"
)

type Config struct {
	VaultAddress		string `mapstructure:"VAULT_ADDRESS"`
	VaultToken		string `mapstructure:"VAULT_TOKEN"`
	ServerPort		string `mapstructure:"SERVER_PORT"`
	InventoryPath	string `mapstructure:"INVENTORY_PATH"`
	Group			string `mapstructure:"GROUP"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found used system env")
	}

	return &Config{
		VaultAddress: os.Getenv("VAULT_ADDRESS"),
		VaultToken: os.Getenv("VAULT_TOKEN"),
		ServerPort: os.Getenv("SERVER_PORT"),
		InventoryPath: os.Getenv("INVENTORY_PATH"),
		Group: os.Getenv("GROUP"),
	}, nil
}
