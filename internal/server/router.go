package server

import (
	"challenge/internal/services"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	controller := NewPackController(services.NewPackService())

	r.HandleFunc("/", newViewController().viewHandler)
	r.HandleFunc("/api/v1/calculate_packs", controller.calculatePacks).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/pack_sizes", controller.readPackSizes).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/pack_sizes", controller.updatePackSizes).Methods(http.MethodPut)

	return r
}
