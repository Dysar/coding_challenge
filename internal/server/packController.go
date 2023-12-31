package server

import (
	"challenge/internal/model"
	"challenge/internal/services"
	"encoding/json"
	"net/http"
)

type PackController struct {
	PackService services.PackService
}

func NewPackController(packService services.PackService) *PackController {
	return &PackController{
		PackService: packService,
	}
}

func (c *PackController) calculatePacks(w http.ResponseWriter, r *http.Request) {

	var orderRequest model.OrderRequest
	err := json.NewDecoder(r.Body).Decode(&orderRequest)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	packsNeeded, err := c.PackService.CalculatePacks(orderRequest.OrderQuantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := model.PacksNeeded{
		OrderQuantity: orderRequest.OrderQuantity,
		Packs:         packsNeeded,
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *PackController) readPackSizes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result := c.PackService.ReadPackSizes()
	response := model.PackSizes{PackSizes: result}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *PackController) updatePackSizes(w http.ResponseWriter, r *http.Request) {
	var request model.PackSizes
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	c.PackService.UpdatePackSizes(request.PackSizes)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
