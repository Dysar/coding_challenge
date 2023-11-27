package server

import (
	"challenge/internal/config"
	"challenge/internal/services"
	"github.com/gorilla/mux"
)

func NewRouter(config *config.Config) *mux.Router {
	r := mux.NewRouter()

	controller := NewPackController(services.NewPackService(config.PackSizes))

	r.HandleFunc("/", viewHandler)
	r.HandleFunc("/api/v1/calculate_packs", controller.calculatePacks)
	return r
}
