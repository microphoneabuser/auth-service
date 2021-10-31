package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Message string `json:"message"`
}

func newErrorResponse(w http.ResponseWriter, message string, statuscode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statuscode)
	res := response{
		Message: message,
	}
	log.Println(message)
	json.NewEncoder(w).Encode(res)
}
