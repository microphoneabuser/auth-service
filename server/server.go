package server

import (
	"net/http"

	"github.com/microphoneabuser/auth-service/pkg/handler"
)

func RunServer(port string, handlers *handler.Handler) error {
	handlers.SetupRoutes()
	return http.ListenAndServe(":"+port, nil)
}
