package services

import (
	"challenge/internal/model"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackService_CalculatePacks(t *testing.T) {

	testCases := []struct {
		name          string
		orderedQty    int
		expectedPacks []model.PackDetails
	}{
		{
			name:       "1 ordered item",
			orderedQty: 1,
			expectedPacks: []model.PackDetails{
				{PackSize: 250, PacksCount: 1},
			},
		},
		{
			name:       "250 ordered items",
			orderedQty: 250,
			expectedPacks: []model.PackDetails{
				{PackSize: 250, PacksCount: 1},
			},
		},
		{
			name:       "251 ordered items",
			orderedQty: 251,
			expectedPacks: []model.PackDetails{
				{PackSize: 500, PacksCount: 1},
			},
		},
		{
			name:       "501 ordered items",
			orderedQty: 501,
			expectedPacks: []model.PackDetails{
				{PackSize: 500, PacksCount: 1},
				{PackSize: 250, PacksCount: 1},
			},
		},
		{
			name:       "899 ordered items",
			orderedQty: 899,
			expectedPacks: []model.PackDetails{
				{PackSize: 1000, PacksCount: 1},
			},
		},
		{
			name:       "12001 ordered items",
			orderedQty: 12001,
			expectedPacks: []model.PackDetails{
				{PackSize: 5000, PacksCount: 2},
				{PackSize: 2000, PacksCount: 1},
				{PackSize: 250, PacksCount: 1},
			},
		},
		{
			name:       "12003 ordered items",
			orderedQty: 12003,
			expectedPacks: []model.PackDetails{
				{PackSize: 5000, PacksCount: 2},
				{PackSize: 2000, PacksCount: 1},
				{PackSize: 250, PacksCount: 1},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewPackService()
			packs, err := service.CalculatePacks(tc.orderedQty)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedPacks, packs)
		})
	}
}

func TestPackService_CalculatePacksForCustomSizes(t *testing.T) {

	testCases := []struct {
		name          string
		orderedQty    int
		expectedPacks []model.PackDetails
		packSizes     []int
		expectedError error
	}{
		{
			name:          "no pack sizes",
			orderedQty:    1,
			expectedPacks: []model.PackDetails{},
			packSizes:     []int{},
			expectedError: errors.New("no pack sizes configured"),
		},
		{
			name:          "negative pack sizes",
			orderedQty:    1,
			expectedPacks: []model.PackDetails{},
			packSizes:     []int{-1},
			expectedError: errors.New("all pack sizes must be positive"),
		},
		{
			name:       "100 packs",
			orderedQty: 100,
			expectedPacks: []model.PackDetails{
				{
					PackSize:   53,
					PacksCount: 2,
				},
			},
			packSizes: []int{53, 31, 23},
		},
		{
			name:       "500k packs",
			orderedQty: 500000,
			expectedPacks: []model.PackDetails{
				{
					//TODO: should be different
					PackSize:   53,
					PacksCount: 9434,
				},
			},
			packSizes: []int{53, 31, 23},
		},
		{
			name:       "107 packs",
			orderedQty: 107,
			expectedPacks: []model.PackDetails{
				{
					PackSize:   53,
					PacksCount: 1,
				},
				{
					PackSize:   31,
					PacksCount: 1,
				},
				{
					PackSize:   23,
					PacksCount: 1,
				},
			},
			packSizes: []int{53, 31, 23},
		},
		{
			name:       "100 packs",
			orderedQty: 100,
			expectedPacks: []model.PackDetails{
				{PackSize: 53, PacksCount: 2},
			},
			packSizes: []int{53, 31, 7},
		},
		{
			name:       "100 packs with packs 3x",
			orderedQty: 100,
			expectedPacks: []model.PackDetails{
				{PackSize: 90, PacksCount: 1},
				{PackSize: 10, PacksCount: 1},
			},
			packSizes: []int{90, 30, 10},
		},
		{
			name:       "15 packs simple case",
			orderedQty: 15,
			packSizes:  []int{5, 4, 3, 2, 1},
			expectedPacks: []model.PackDetails{
				{PackSize: 5, PacksCount: 3},
			},
		},
		{
			name:       "14 packs simple case",
			orderedQty: 14,
			packSizes:  []int{5, 4, 3, 2, 1},
			expectedPacks: []model.PackDetails{
				{PackSize: 5, PacksCount: 2},
				{PackSize: 4, PacksCount: 1},
			},
		},
		{
			name:       "15 packs simple case",
			orderedQty: 9,
			packSizes:  []int{5, 4, 3, 2, 1},
			expectedPacks: []model.PackDetails{
				{PackSize: 5, PacksCount: 1},
				{PackSize: 4, PacksCount: 1},
			},
		},
		{
			name:       "3 packs simple case",
			orderedQty: 3,
			packSizes:  []int{5, 4, 3, 2, 1},
			expectedPacks: []model.PackDetails{
				{PackSize: 3, PacksCount: 1},
			},
		},
		{
			name:       "300000 packs simple case",
			orderedQty: 300000,
			packSizes:  []int{5, 4, 3, 2, 1},
			expectedPacks: []model.PackDetails{
				{PackSize: 5, PacksCount: 60000},
			},
		},
		{
			name:       "910 ordered, custom packs",
			orderedQty: 910,
			packSizes:  []int{500, 100, 20},
			expectedPacks: []model.PackDetails{
				{PackSize: 500, PacksCount: 1},
				{PackSize: 100, PacksCount: 4},
				{PackSize: 20, PacksCount: 1},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewPackService()
			service.packSizes = tc.packSizes
			packs, err := service.CalculatePacks(tc.orderedQty)
			if tc.expectedError != nil {
				assert.Equal(t, err, tc.expectedError)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedPacks, packs)
		})
	}
}

func TestPackServiceImpl_UpdatePackSizes(t *testing.T) {
	svc := NewPackService()
	svc.UpdatePackSizes([]int{1, 2, 3})
	assert.Equal(t, []int{3, 2, 1}, svc.packSizes)
}

func TestPackServiceImpl_ReadPackSizes(t *testing.T) {
	svc := NewPackService()
	svc.UpdatePackSizes([]int{1, 2, 3})

	pss := svc.ReadPackSizes()
	assert.Equal(t, []int{1, 2, 3}, pss)
}

func TestPackServiceImpl_adjustPacks(t *testing.T) {
	testCases := []struct {
		name          string
		packSizes     []int
		inputPacks    model.Packs
		expectedPacks model.Packs
	}{
		{
			name:      "1",
			packSizes: []int{20, 100, 500},
			inputPacks: map[int]model.CountAndQuantity{
				20:  model.NewCountAndQuantity(1, 20),
				100: model.NewCountAndQuantity(4, 400),
				500: model.NewCountAndQuantity(1, 500),
			},
			expectedPacks: map[int]model.CountAndQuantity{
				20:  model.NewCountAndQuantity(1, 20),
				100: model.NewCountAndQuantity(4, 400),
				500: model.NewCountAndQuantity(1, 500),
			},
		},
		{
			name:      "2",
			packSizes: []int{80, 100, 500},
			inputPacks: map[int]model.CountAndQuantity{
				80:  model.NewCountAndQuantity(1, 80),
				100: model.NewCountAndQuantity(4, 400),
				500: model.NewCountAndQuantity(1, 500),
			},
			expectedPacks: map[int]model.CountAndQuantity{
				500: model.NewCountAndQuantity(2, 980),
			},
		},
		{
			name:      "3",
			packSizes: []int{250, 500},
			inputPacks: map[int]model.CountAndQuantity{
				250: model.NewCountAndQuantity(2, 500),
			},
			expectedPacks: map[int]model.CountAndQuantity{
				500: model.NewCountAndQuantity(1, 500),
			},
		},
		{
			name:      "21+31+53=105; 53*2=106",
			packSizes: []int{7, 31, 53},
			inputPacks: map[int]model.CountAndQuantity{
				53: model.NewCountAndQuantity(1, 53),
				31: model.NewCountAndQuantity(1, 31),
				7:  model.NewCountAndQuantity(3, 21),
			},
			expectedPacks: map[int]model.CountAndQuantity{
				53: model.NewCountAndQuantity(2, 105),
			},
		},
		{
			name:      "4",
			packSizes: []int{7, 31, 53},
			inputPacks: map[int]model.CountAndQuantity{
				53: model.NewCountAndQuantity(1, 53),
				31: model.NewCountAndQuantity(1, 31),
				7:  model.NewCountAndQuantity(4, 28),
			},
			expectedPacks: map[int]model.CountAndQuantity{
				53: model.NewCountAndQuantity(1, 53),
				31: model.NewCountAndQuantity(2, 59),
			},
		},
		{
			name:      "5",
			packSizes: []int{20, 100, 500},
			inputPacks: map[int]model.CountAndQuantity{
				20:  model.NewCountAndQuantity(1, 20),
				100: model.NewCountAndQuantity(4, 400),
				500: model.NewCountAndQuantity(1, 500),
			},
			expectedPacks: map[int]model.CountAndQuantity{
				20:  model.NewCountAndQuantity(1, 20),
				100: model.NewCountAndQuantity(4, 400),
				500: model.NewCountAndQuantity(1, 500),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewPackService()

			service.adjustMultipleSameSizePacks(tc.packSizes, tc.inputPacks)
			assert.Equal(t, tc.expectedPacks, tc.inputPacks)
		})
	}

}
