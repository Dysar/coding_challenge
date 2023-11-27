package services

import (
	"challenge/internal/model"
	"errors"
	"github.com/sirupsen/logrus"
	"sort"
)

type (
	PackService interface {
		CalculatePacks(orderQuantity int) ([]model.PackDetails, error)
	}
	PackServiceImpl struct {
		packSizes []int
	}
)

func NewPackService(packSizes []int) *PackServiceImpl {
	sort.Sort(sort.Reverse(sort.IntSlice(packSizes)))
	return &PackServiceImpl{
		packSizes: packSizes,
	}
}

func (s *PackServiceImpl) CalculatePacks(orderQuantity int) ([]model.PackDetails, error) {
	if orderQuantity <= 0 {
		return nil, errors.New("order quantity must be greater than 0")
	}
	if len(s.packSizes) == 0 {
		return nil, errors.New("no pack sizes configured")
	}

	logrus.Info(s.packSizes)

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
		logrus.Infof("remainingQuantity: %d, packsCount: %d, packSize: %d, i: %d", remainingQuantity, packsCount, packSize, i)
		if packsCount > 0 {
			logrus.Infof("adding a pack size %d, count: %d", packSize, packsCount)
			packMap[packSize] = packsCount
			packedQuantity += packSize * packsCount
			remainingQuantity %= packSize
		}

		if remainingQuantity == 0 {
			continue
		}

		logrus.Infof("remainingQuantity:%d", remainingQuantity)

		// if the remainingQuantity is less than the smallest pack try to add the whole quantity to a larger box (if there are larger boxes)
		// but only if the box is not too large. Let's assume that the first half of boxes are too large
		if remainingQuantity < smallestPack &&
			i > len(s.packSizes)/2 && i != lastIndex && packedQuantity != 0 && packsCount != 0 {
			logrus.Infof("removing a pack size %d, count: %d", packSize, 1)
			packMap[packSize] -= 1
			logrus.Infof("adding a larger pack size %d, count: %d", s.packSizes[i-1], 1)
			packMap[s.packSizes[i-1]] = 1
			remainingQuantity = 0
		}

		//if there is any remaining quantity when all the pack sizes have been checked
		if i == lastIndex {
			if packsCount != 0 {
				//if any smallest (250) packs have been filled, we need to upgrade one 250 pack to a 500 pack
				logrus.Infof("removing a pack size %d, count: %d", packSize, 1)
				packMap[packSize] -= 1
				logrus.Infof("adding a pack size %d, count: %d", s.packSizes[i-1], 1)
				packMap[s.packSizes[i-1]] += 1
			} else {
				//fill one smallest pack (250)
				logrus.Infof("adding one smallest pack size %d, count: %d", packSize, 1)
				packMap[packSize] += 1
			}
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
