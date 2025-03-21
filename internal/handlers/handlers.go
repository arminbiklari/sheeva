package handlers

import (
	"log"
	"net/http"
	"strings"
	"sheeva/internal/vault"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	vaultClient *vault.VaultClient
}

func NewHandler(vc *vault.VaultClient) *Handler {
	return &Handler{vaultClient: vc}
}


func (h *Handler) GetSecret(c *gin.Context) {
	secretPath := c.Param("path")
	log.Printf("Attempting to fetch secret from path: %s", secretPath)
	
	vaultToken := c.GetHeader("token")
	if vaultToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	h.vaultClient.SetHandlerToken(vaultToken)
	secret, err := h.vaultClient.GetSecret(secretPath)
	if err != nil {
		log.Printf("Error fetching secret: %v", err)

		switch {
		case strings.Contains(err.Error(), "permission denied") || 
			 strings.Contains(err.Error(), "invalid token"):
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		case strings.Contains(err.Error(), "not found") || 
			 strings.Contains(err.Error(), "no secret"):
			c.JSON(http.StatusNotFound, gin.H{"error": "Secret not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get secret"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"secret": secret})
}

func (h *Handler) SetSecret(c *gin.Context) {
	secretPath := c.Param("path")
	log.Printf("Attempting to set secret at path: %s", secretPath)

	// Check token authentication first
	vaultToken := c.GetHeader("token")
	if vaultToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	// Parse the secret data from request body
	var secretData map[string]interface{}
	if err := c.ShouldBindJSON(&secretData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Set token for this operation
	h.vaultClient.SetHandlerToken(vaultToken)

	// Attempt to create the secret
	err := h.vaultClient.CreateSecret(secretPath, secretData)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "permission denied") || 
			 strings.Contains(err.Error(), "invalid token"):
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		case strings.Contains(err.Error(), "already exists"):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create secret"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Secret created successfully"})
}
