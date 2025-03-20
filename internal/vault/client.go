package vault

import (
	"encoding/json"
	"fmt"
	"log"
	"sheeva/config"
	vault "github.com/hashicorp/vault/api"
)

type VaultClient struct{
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
