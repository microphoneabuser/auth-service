package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/microphoneabuser/auth-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) getTokens(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()

	id, err := primitive.ObjectIDFromHex(keys.Get("id"))
	if err != nil {
		newErrorResponse(w, "invalid id param", http.StatusBadRequest)
		return
	}

	tokens, err := h.service.Authorization.GenerateTokens(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrorUserNotFound) {
			newErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res := tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}

func (h *Handler) refreshTokens(w http.ResponseWriter, r *http.Request) {
	var input refreshInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, "invalid request body", http.StatusBadRequest)
		return
	}
	id, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		newErrorResponse(w, "invalid id field", http.StatusBadRequest)
		return
	}
	tokens, err := h.service.Authorization.RefreshTokens(r.Context(), id, input.Token)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res := tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.Authorization.CreateUser(r.Context())
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res := createResponse{
		Id: user.ID.Hex(),
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type refreshInput struct {
	Id    string `json:"id" binding:"required"`
	Token string `json:"token" binding:"required"`
}

type createResponse struct {
	Id string `json:"id"`
}
