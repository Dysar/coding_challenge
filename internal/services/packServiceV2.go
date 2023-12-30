package services

import (
	"challenge/internal/model"
	"errors"
	"github.com/sirupsen/logrus"
	"slices"
)

type (
	PackServiceV2 interface {
		CalculatePacks(orderQuantity int) ([]model.PackDetails, error)
	}

	PackServiceImplV2 struct {
		// pack sizes are sorted desc
		packSizes []int
	}
)

func NewPackServiceV2() *PackServiceImplV2 {
	initialPackSizes := []int{250, 500, 1000, 2000, 5000}
	slices.Reverse(initialPackSizes)
	return &PackServiceImplV2{
		packSizes: initialPackSizes,
	}
}

func (s *PackServiceImplV2) CalculatePacks(orderQuantity int) ([]model.PackDetails, error) {

	if orderQuantity <= 0 {
		return nil, errors.New("order quantity must be greater than 0")
	}
	if len(s.packSizes) == 0 {
		return nil, errors.New("no pack sizes configured")
	}

	var (
		packMap           = make(model.Packs)
		remainingQuantity = orderQuantity
		lastIndex         = len(s.packSizes) - 1
		smallestPack      = s.packSizes[lastIndex]
	)

	// Calculate the number of packs needed for each pack size
	for _, packSize := range s.packSizes {
		if packSize < 0 {
			return nil, errors.New("all pack sizes must be positive")
		}
		packsCount := remainingQuantity / packSize
		logrus.Infof("packsCount(%d) := remainingQuantity(%d) / packSize(%d)", packsCount, remainingQuantity, packSize)
		if packsCount > 0 {
			logrus.Infof("adding a pack size %d, count: %d", packSize, packsCount)
			packMap.AddPacks(packSize, packsCount, packSize*packsCount)
			remainingQuantity %= packSize
		}

		if remainingQuantity == 0 {
			continue
		}

		//if there is any remaining quantity when all the pack sizes have been checked
		if remainingQuantity <= smallestPack {

			//fill one smallest pack (250)
			logrus.Infof("adding one smallest pack size %d, count: %d", smallestPack, 1)
			packMap.AddPack(smallestPack, remainingQuantity)
			break
		}
	}

	//var x1, x2, x3,
	var sum int
	values := make([]int, len(s.packSizes))

	for i, packSize := range s.packSizes {
		if cq, ok := packMap[packSize]; ok {
			sum += packSize * cq.Count
			values[i] = cq.Count
		}
	}

	c := newCalculator(values, s.packSizes, orderQuantity, sum)
	values = c.flexibleCalculate()
	result := []model.PackDetails{}
	for i, value := range values {
		if value == 0 {
			continue
		}
		size := s.packSizes[i]
		result = append(result, model.PackDetails{PackSize: size, PacksCount: value})
	}

	//TODO: Within the constraints above, send out Rules 1 & 2 send out as few packs as possible to fulfil each order.
	return result, nil
}

func (s *PackServiceImplV2) UpdatePackSizes(packSizes []int) {
	slices.Sort(packSizes)
	slices.Reverse(packSizes)
	s.packSizes = packSizes
}
