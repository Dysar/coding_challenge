package server

import (
	"bytes"
	"challenge/internal/model"
	"challenge/mocks"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPackController_CalculatePacks(t *testing.T) {

	t.Run("Valid request", func(t *testing.T) {
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

func TestPackController_ReadPackSizes(t *testing.T) {
	packService := mocks.PackService{}

	packService.On("ReadPackSizes").Return([]int{1, 2, 3}).Once()

	packController := NewPackController(&packService)

	// Create a request for the readPackSizes handler
	req, err := http.NewRequest(http.MethodGet, "/api/v1/pack_sizes", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the readPackSizes handler
	http.HandlerFunc(packController.readPackSizes).ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Decode the response body
	var response model.PackSizes
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	expectedPackSizes := []int{1, 2, 3}
	assert.Equal(t, expectedPackSizes, response.PackSizes)
}

func TestPackController_UpdatePackSizes(t *testing.T) {
	packService := mocks.PackService{}

	packService.On("UpdatePackSizes", []int{1, 2, 3}).Return([]int{1, 2, 3}).Once()

	packController := NewPackController(&packService)

	requestPayload := model.PackSizes{PackSizes: []int{1, 2, 3}}

	// Marshal the request payload to JSON
	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		t.Fatal(err)
	}

	// Create a request for the readPackSizes handler
	req, err := http.NewRequest(http.MethodPut, "/api/v1/pack_sizes", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the readPackSizes handler
	http.HandlerFunc(packController.updatePackSizes).ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Decode the response body
	var response model.PackSizes
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	expectedPackSizes := []int{1, 2, 3}
	assert.Equal(t, expectedPackSizes, response.PackSizes)
}
