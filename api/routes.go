package api

import (
	"sheeva/internal/handlers"
	"sheeva/internal/vault"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, vaultClient *vault.VaultClient) {
	h := handlers.NewHandler(vaultClient)

	router.POST("/auth/token", h.AuthWithToken)
	router.GET("/vault/*path", h.GetSecret)
}
