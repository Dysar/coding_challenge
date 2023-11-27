package server

import (
	"challenge/internal/model"
	"challenge/mocks"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculatePacks(t *testing.T) {

	t.Run("Valid request", func(t *testing.T) {
		// Create a PackController with a MockPackService
		packSvc := &mocks.PackService{}
		controller := NewPackController(packSvc)

		reqBody := `{"order_quantity": 10}`
		req, err := http.NewRequest("POST", "/calculate_packs", strings.NewReader(reqBody))
		assert.NoError(t, err)

		packSvc.On("CalculatePacks", 10).Return([]model.PackDetails{{PackSize: 5, PacksCount: 2}}, nil).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.calculatePacks)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response model.PacksNeeded
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, 10, response.OrderQuantity)
		assert.Equal(t, []model.PackDetails{{PackSize: 5, PacksCount: 2}}, response.Packs)
	})

	t.Run("Valid request 2", func(t *testing.T) {
		// Create a PackController with a MockPackService
		packSvc := &mocks.PackService{}
		controller := NewPackController(packSvc)

		reqBody := `{"order_quantity": 15}`
		req, err := http.NewRequest("POST", "/calculate_packs", strings.NewReader(reqBody))
		assert.NoError(t, err)

		packSvc.On("CalculatePacks", 15).Return([]model.PackDetails{{PackSize: 10, PacksCount: 1}, {PackSize: 5, PacksCount: 1}}, nil).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.calculatePacks)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response model.PacksNeeded
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, 15, response.OrderQuantity)
		assert.Equal(t, []model.PackDetails{{PackSize: 10, PacksCount: 1}, {PackSize: 5, PacksCount: 1}}, response.Packs)
	})

	t.Run("invalid request", func(t *testing.T) {
		// Create a PackController with a MockPackService
		packSvc := &mocks.PackService{}
		controller := NewPackController(packSvc)

		reqBody := `{"order_quantity": "invalid"}`
		req, err := http.NewRequest("POST", "/calculate_packs", strings.NewReader(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.calculatePacks)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
