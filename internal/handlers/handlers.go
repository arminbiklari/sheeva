package handlers

import (
	"net/http"
	"log"
	"github.com/gin-gonic/gin"
	"sheeva/internal/vault"
)

type Handler struct {
	vaultClient *vault.VaultClient
}

func NewHandler(vc *vault.VaultClient) *Handler {
	return &Handler{vaultClient: vc}
}

func (h *Handler) AuthWithToken(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	h.vaultClient.Client.SetToken(req.Token)
	c.JSON(http.StatusOK, gin.H{"message": "Authenticated successfully"})
}

func (h *Handler) GetSecret(c *gin.Context) {
	secretPath := c.Param("path")
	log.Println("secretPath: ", secretPath)
	vaultToken := c.GetHeader("token")
	if vaultToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	log.Println(h.vaultClient)
	h.vaultClient.SetHandlerToken(vaultToken)
	secret, err := h.vaultClient.GetSecret(secretPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get secret"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"secret": secret})
}
