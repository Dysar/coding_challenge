package services

import (
	"challenge/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackServiceV2(t *testing.T) {

	s := NewPackServiceV2()

	t.Run("107", func(t *testing.T) {
		res, err := s.CalculatePacks(107)
		assert.NoError(t, err)
		assert.Equal(t, []model.PackDetails{
			{PackSize: 53, PacksCount: 1},
			{PackSize: 31, PacksCount: 1},
			{PackSize: 23, PacksCount: 1},
		}, res)
	})

	t.Run("117", func(t *testing.T) {
		res, err := s.CalculatePacks(117)
		assert.NoError(t, err)
		t.Logf("%+v", res)
		assert.Equal(t, []model.PackDetails{
			{PackSize: 53, PacksCount: 1},
			{PackSize: 31, PacksCount: 0},
			{PackSize: 23, PacksCount: 3},
		}, res)
	})

	t.Run("500k", func(t *testing.T) {
		res, err := s.CalculatePacks(500000)
		assert.NoError(t, err)
		assert.Equal(t, []model.PackDetails{
			{PackSize: 53, PacksCount: 9429},
			{PackSize: 31, PacksCount: 7},
			{PackSize: 23, PacksCount: 2},
		}, res)
	})

	t.Run("100", func(t *testing.T) {
		res, err := s.CalculatePacks(100)
		assert.NoError(t, err)
		assert.Equal(t, []model.PackDetails{
			{PackSize: 53, PacksCount: 0},
			{PackSize: 31, PacksCount: 1},
			{PackSize: 23, PacksCount: 3},
		}, res)
	})

}
