package services

import (
	"challenge/internal/model"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackServiceImplV2_CalculatePacks(t *testing.T) {
	testCases := []struct {
		name          string
		orderedQty    int
		expectedPacks []model.PackDetails
		packSizes     []int
		expectedError error
	}{
		//{
		//	name:          "no pack sizes",
		//	orderedQty:    1,
		//	expectedPacks: []model.PackDetails{},
		//	packSizes:     []int{},
		//	expectedError: errors.New("no pack sizes configured"),
		//},
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

			//time="2024-01-01T14:41:17+02:00" level=info msg="count: 2, smallerPacksQuantity: 0"
			//time="2024-01-01T14:41:17+02:00" level=info msg="totalQuantityWithSmallerPacks: 2"
			//time="2024-01-01T14:41:17+02:00" level=info msg="biggerPackSize(500) is greater or equal to the total quantity"
			//time="2024-01-01T14:41:17+02:00" level=info msg="isDiffMoreThanSmallestPack: true, biggerPackSize: 500, totalQuantityWithSmallerPacks: 2, smallestPackSize:250"
			//time="2024-01-01T14:41:17+02:00" level=info msg="checking biggest packs; smallerPacksQuantity: 2"

			//time="2024-01-01T14:41:51+02:00" level=info msg="countAndQuantity: {2 251}, smallerPacksQuantity: 0"
			//time="2024-01-01T14:41:51+02:00" level=info msg="totalQuantityWithSmallerPacks: 251"
			//time="2024-01-01T14:41:51+02:00" level=info msg="biggerPackSize(500) is greater or equal to the total quantity"
			//time="2024-01-01T14:41:51+02:00" level=info msg="isDiffMoreThanSmallestPack: false, biggerPackSize: 500, totalQuantityWithSmallerPacks: 251, smallestPackSize:250"
			//time="2024-01-01T14:41:51+02:00" level=info msg="adding pack; size 500, quantity: 251"
			//time="2024-01-01T14:41:51+02:00" level=info msg="setting count of pack sizes 250 to 0"
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
			service := NewPackServiceV2()
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
