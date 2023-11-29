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
				{PackSize: 500, PacksCount: 1},
				{PackSize: 250, PacksCount: 1},
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
			service := NewPackService()
			packs, err := service.CalculatePacks(tc.OrderedQty)
			assert.NoError(t, err)
			assert.Equal(t, tc.ExpectedPacks, packs)
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
