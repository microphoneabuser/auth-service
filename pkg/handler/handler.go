package handler

import (
	"net/http"

	"github.com/microphoneabuser/auth-service/pkg/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SetupRoutes() {
	http.HandleFunc("/auth/get-tokens", h.getTokens)
	http.HandleFunc("/auth/refresh", h.refreshTokens)
	http.HandleFunc("/auth/create-user", h.createUser)
}
