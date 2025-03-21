package vault

import (
	"encoding/json"
	"fmt"
	"log"
	"sheeva/config"

	vault "github.com/hashicorp/vault/api"
)

type VaultClient struct {
	Client *vault.Client
}

func NewVaultClient(LoadedConfig *config.Config) (*VaultClient, error) {

	config := vault.DefaultConfig()
	config.Address = LoadedConfig.VaultAddress

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}
	return &VaultClient{Client: client}, nil
}

func (v *VaultClient) SetHandlerToken(token string) {
	v.Client.SetToken(token)
}

func (v *VaultClient) GetSecret(path string) (string,error) {

	secret, err := v.Client.Logical().Read(path)
	if err != nil {
		log.Println("unable to read secret: %w", err)
		return "", err
	}

	if secret == nil {
		return "", fmt.Errorf("secret not found")
	}
	jsonData, err := json.Marshal(secret.Data)
	if err != nil {
		log.Println("unable to marshal secret: %w", err)
		return "", err
	}
	return string(jsonData), nil
}


func (v *VaultClient) CreateSecret(path string, secret map[string]interface{}) error {
	existing, err := v.Client.Logical().Read(path)
	if err != nil {
		return err
	}
	
	if existing == nil {
		_, err := v.Client.Logical().Write(path, secret)
		if err != nil {
			log.Printf("unable to write secret: %v in path: %s", err, path)
			return err
		}
		return nil
	}
	
	log.Printf("secret already exists in path: %s", path)
	return fmt.Errorf("secret already exists in path: %s", path)
}
