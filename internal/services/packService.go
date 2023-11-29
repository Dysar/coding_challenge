package services

import (
	"challenge/internal/model"
	"errors"
	"github.com/sirupsen/logrus"
	"slices"
)

type (
	PackService interface {
		CalculatePacks(orderQuantity int) ([]model.PackDetails, error)
		UpdatePackSizes(packSizes []int)
		ReadPackSizes() []int
	}
	PackServiceImpl struct {
		packSizes []int
	}
)

func NewPackService() *PackServiceImpl {
	initialPackSizes := []int{250, 500, 1000, 2000, 5000}
	slices.Reverse(initialPackSizes)
	return &PackServiceImpl{
		packSizes: initialPackSizes,
	}
}

func (s *PackServiceImpl) UpdatePackSizes(packSizes []int) {
	slices.Sort(packSizes)
	slices.Reverse(packSizes)
	s.packSizes = packSizes
}
func (s *PackServiceImpl) ReadPackSizes() []int {
	clone := slices.Clone(s.packSizes)
	slices.Reverse(clone)
	return clone
}

func (s *PackServiceImpl) CalculatePacks(orderQuantity int) ([]model.PackDetails, error) {

	if orderQuantity <= 0 {
		return nil, errors.New("order quantity must be greater than 0")
	}
	if len(s.packSizes) == 0 {
		return nil, errors.New("no pack sizes configured")
	}

	var (
		packedQuantity    int
		packMap           = make(map[int]int)
		remainingQuantity = orderQuantity
		lastIndex         = len(s.packSizes) - 1
		smallestPack      = s.packSizes[lastIndex]
	)

	// Calculate the number of packs needed for each pack size
	for i, packSize := range s.packSizes {
		packsCount := remainingQuantity / packSize
		if packsCount > 0 {
			logrus.Infof("adding a pack size %d, count: %d", packSize, packsCount)
			packMap[packSize] = packsCount
			packedQuantity += packSize * packsCount
			remainingQuantity %= packSize
		}

		if remainingQuantity == 0 {
			continue
		}

		// if the remainingQuantity is less than the smallest pack try to add the whole quantity to a larger box (if there are larger boxes)
		// but only if the box is not too large. Let's assume that the first half of boxes are too large
		if remainingQuantity < smallestPack &&
			i > len(s.packSizes)/2 && packedQuantity != 0 && packsCount != 0 {
			logrus.Infof("removing a pack size %d, count: %d", packSize, 1)
			packMap[packSize] -= 1
			logrus.Infof("adding a larger pack size %d, count: %d", s.packSizes[i-1], 1)
			packMap[s.packSizes[i-1]] += 1
			remainingQuantity = 0
		}

		//if there is any remaining quantity when all the pack sizes have been checked
		if i == lastIndex && packsCount == 0 {
			//fill one smallest pack (250)
			logrus.Infof("adding one smallest pack size %d, count: %d", packSize, 1)
			packMap[packSize] += 1
			remainingQuantity = 0
		}
	}

	return s.packsToResponse(packMap), nil
}

func (s *PackServiceImpl) packsToResponse(packMap map[int]int) []model.PackDetails {
	var packsNeeded []model.PackDetails
	for _, packSize := range s.packSizes {
		if count, ok := packMap[packSize]; ok && count > 0 {
			packsNeeded = append(packsNeeded, model.PackDetails{
				PackSize:   packSize,
				PacksCount: count,
			})
		}
	}

	return packsNeeded
}
