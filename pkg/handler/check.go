package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (h *Handler) checkAccess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	header := r.Header.Get("Authorization")
	if header == "" {
		newErrorResponse(w, "auth header is empty", http.StatusUnauthorized)
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(w, "invalid authorization header", http.StatusUnauthorized)
		return
	}
	if len(headerParts[1]) == 0 {
		newErrorResponse(w, "token is empty", http.StatusUnauthorized)
		return
	}

	id, err := h.service.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res := checkResponse{
		Id: id,
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}

type checkResponse struct {
	Id string `json:"id"`
}
