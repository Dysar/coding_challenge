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

		//FLEXIBLE
		//-----iteration-----
		//x1: 2 x2: 0 x3: 1 sum:129
		//x1: 2 x2: 0 x3: 1 sum:129
		//-----iteration-----
		//x1: 1 x2: 0 x3: 1 sum:76
		//x1: 1 x2: 0 x3: 1 sum:76
		//x1: 1 x2: 0 x3: 1 sum:76
		//x1: 1 x2: 0 x3: 2 sum:99
		//x1: 1 x2: 0 x3: 3 sum:122
		//target: 117, initial: 129 best count so far: 122
		//-----iteration-----
		//x1: 0 x2: 0 x3: 1 sum:23
		//x1: 0 x2: 0 x3: 1 sum:23
		//x1: 0 x2: 0 x3: 1 sum:23
		//x1: 0 x2: 0 x3: 2 sum:46
		//x1: 0 x2: 0 x3: 3 sum:69
		//x1: 0 x2: 0 x3: 4 sum:92
		//x1: 0 x2: 0 x3: 5 sum:115
		//x1: 0 x2: 0 x3: 6 sum:138

		//STATIC
		//-----iteration-----
		//x1: 2, x2: 0, x3: 1, sum: 129
		//x1: 2, x2: 0, x3: 1, sum: 129
		//-----iteration-----
		//x1: 1, x2: 0, x3: 1, sum: 76
		//x1: 1, x2: 0, x3: 1, sum: 76
		//x1: 1, x2: 0, x3: 1, sum: 76
		//x1: 1, x2: 0, x3: 2, sum: 99
		//x1: 1, x2: 0, x3: 3, sum: 122
		//target: 117, initial: 129 best count so far: 122
		//x1: 1, x2: 1, x3: 0, sum: 84
		//x1: 1, x2: 1, x3: 0, sum: 84
		//x1: 1, x2: 1, x3: 1, sum: 107
		//x1: 1, x2: 1, x3: 2, sum: 130
		//x1: 1, x2: 2, x3: 0, sum: 115
		//x1: 1, x2: 2, x3: 0, sum: 115
		//x1: 1, x2: 2, x3: 1, sum: 138
		//x1: 1, x2: 3, x3: 0, sum: 146
		//-----iteration-----
		//x1: 0, x2: 3, x3: 0, sum: 93
		//x1: 0, x2: 3, x3: 0, sum: 93
		//x1: 0, x2: 3, x3: 0, sum: 93
		//x1: 0, x2: 3, x3: 1, sum: 116
		//x1: 0, x2: 3, x3: 2, sum: 139
		//x1: 0, x2: 4, x3: 0, sum: 124
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
