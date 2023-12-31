package services

import (
	"challenge/internal/model"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackServiceImpl_CalculatePacks(t *testing.T) {
	testCases := []struct {
		name          string
		orderedQty    int
		expectedPacks []model.PackDetails
		packSizes     []int
		expectedError error
	}{
		{
			name:          "negative pack sizes",
			orderedQty:    1,
			expectedPacks: []model.PackDetails{},
			packSizes:     []int{-1},
			expectedError: errors.New("all pack sizes must be positive"),
		},
		{
			name:       "107 packs",
			orderedQty: 107,
			expectedPacks: []model.PackDetails{
				{PackSize: 53, PacksCount: 1},
				{PackSize: 31, PacksCount: 1},
				{PackSize: 23, PacksCount: 1},
			},
			packSizes: []int{53, 31, 23},
		},
		{
			name:       "117 packs",
			orderedQty: 117,
			expectedPacks: []model.PackDetails{
				{PackSize: 53, PacksCount: 1},
				{PackSize: 23, PacksCount: 3},
			},
			packSizes: []int{53, 31, 23},
		},
		{
			name:       "100 packs",
			orderedQty: 100,
			expectedPacks: []model.PackDetails{
				{PackSize: 31, PacksCount: 1},
				{PackSize: 23, PacksCount: 3},
			},
			packSizes: []int{53, 31, 23},
		},
		{
			name:       "500k packs",
			orderedQty: 500000,
			expectedPacks: []model.PackDetails{
				{PackSize: 53, PacksCount: 9429},
				{PackSize: 31, PacksCount: 7},
				{PackSize: 23, PacksCount: 2},
			},
			packSizes: []int{53, 31, 23},
		},
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
		{
			name:       "100 packs",
			orderedQty: 100,
			expectedPacks: []model.PackDetails{
				{PackSize: 31, PacksCount: 1},
				{PackSize: 7, PacksCount: 10},
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
			if len(tc.packSizes) != 0 {
				service.packSizes = tc.packSizes
			}
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
