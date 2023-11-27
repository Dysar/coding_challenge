package services

import (
	"challenge/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackService_CalculatePacks(t *testing.T) {

	testCases := []struct {
		Name          string
		OrderedQty    int
		ExpectedPacks []model.PackDetails
	}{
		{
			Name:       "1 ordered item",
			OrderedQty: 1,
			ExpectedPacks: []model.PackDetails{
				{PackSize: 250, PacksCount: 1},
			},
		},
		{
			Name:       "250 ordered items",
			OrderedQty: 250,
			ExpectedPacks: []model.PackDetails{
				{PackSize: 250, PacksCount: 1},
			},
		},
		{
			Name:       "251 ordered items",
			OrderedQty: 251,
			ExpectedPacks: []model.PackDetails{
				{PackSize: 500, PacksCount: 1},
			},
		},
		{
			Name:       "501 ordered items",
			OrderedQty: 501,
			ExpectedPacks: []model.PackDetails{
				{PackSize: 1000, PacksCount: 1},
			},
		},
		{
			Name:       "12001 ordered items",
			OrderedQty: 12001,
			ExpectedPacks: []model.PackDetails{
				{PackSize: 5000, PacksCount: 2},
				{PackSize: 2000, PacksCount: 1},
				{PackSize: 250, PacksCount: 1},
			},
		},
		{
			Name:       "12003 ordered items",
			OrderedQty: 12003,
			ExpectedPacks: []model.PackDetails{
				{PackSize: 5000, PacksCount: 2},
				{PackSize: 2000, PacksCount: 1},
				{PackSize: 250, PacksCount: 1},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			service := NewPackService([]int{250, 500, 1000, 2000, 5000})
			packs, err := service.CalculatePacks(tc.OrderedQty)
			assert.NoError(t, err)
			assert.Equal(t, tc.ExpectedPacks, packs)
		})
	}
}

func TestPackService_CalculatePacksV2(t *testing.T) {
	s := NewPackService([]int{250, 500, 1000, 2000, 5000})
	t.Log(len(s.packSizes) / 2)
}
