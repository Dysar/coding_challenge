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

	s.adjustPacks(packMap)
	return s.packsToResponse(packMap), nil
}

// the adjustPacks method modifies the packMap so that the number of packs would be optimal
func (s *PackServiceImpl) adjustPacks(packMap model.Packs) {

	packSizesAsc := slices.Clone(s.packSizes)
	slices.Reverse(packSizesAsc)

	s.adjustMultipleSameSizePacks(packSizesAsc, packMap)
}

func (s *PackServiceImpl) adjustMultipleSameSizePacks(packSizesAsc []int, packMap model.Packs) {
	var smallerPacksQuantity int
	smallestPackSize := packSizesAsc[0]

	for i := 0; i < len(packSizesAsc); i++ {
		packSize := packSizesAsc[i]

		cq, ok := packMap[packSize]
		if !ok {
			continue
		}
		// add only positive cq items to the response
		if cq.Count <= 0 {
			delete(packMap, packSize)
		}

		if i == len(packSizesAsc)-1 {
			//if it's the biggest pack, we still might want to make adjustments
			logrus.Infof("checking biggest packs; smallerPacksQuantity: %d", smallerPacksQuantity)
			if smallerPacksQuantity > packSize {
				break
			}
			//smallerPacksQuantity <= packSize
			if isDiffLessThanSmallestPack := (packSize - smallerPacksQuantity) < smallestPackSize; isDiffLessThanSmallestPack {
				logrus.Infof("adding biggest pack %d with q:%d", packSize, smallerPacksQuantity)
				packMap.AddPack(packSize, smallerPacksQuantity)
				for k := i - 1; k >= 0; k-- {
					logrus.Infof("255: setting count of pack sizes %d to 0", packSizesAsc[k])
					packMap.SetCount(packSizesAsc[k], 0)
				}
			}
		}

		if cq.Count <= 1 || i == len(packSizesAsc)-1 {
			smallerPacksQuantity += cq.Quantity
			continue
		}

		//there is more than 1 pack they can potentially be replaced by a larger pack

		logrus.Infof("countAndQuantity: %d, smallerPacksQuantity: %d", cq, smallerPacksQuantity)

		totalQuantityWithSmallerPacks := cq.Quantity + smallerPacksQuantity
		logrus.Infof("totalQuantityWithSmallerPacks: %d", totalQuantityWithSmallerPacks)

		biggerPackSize := packSizesAsc[i+1]

		//it can happen that the total quantity is < or > or = than the bigger pack
		//if it's <, we need to make sure that the difference is less than the smallest pack

		if biggerPackSize < totalQuantityWithSmallerPacks {
			logrus.Infof("biggerPackSize(%d) is less than total quantity", biggerPackSize)
			//what if the bigger pack is 500 and total quantity of smaller packs is 2000?
			//	in that case I would look for bigger pack that could fit those items. next pack would be 1000, 1000<2000, next pack would be 2000, that is fine, we can use that one
			//	2000-2000=0; 0<250 (smallest pack), all good
			smallerPacksQuantity += cq.Quantity
			continue
		} else {
			logrus.Infof("biggerPackSize(%d) is greater or equal to the total quantity", biggerPackSize)

			//biggerPackSize => total quantity, e.g. pack size = 500, total quantity 100*4+80=480
			isDiffMoreThanSmallestPack := (biggerPackSize - totalQuantityWithSmallerPacks) > smallestPackSize
			logrus.Infof("isDiffMoreThanSmallestPack: %t, biggerPackSize: %d, totalQuantityWithSmallerPacks: %d, smallestPackSize:%d",
				isDiffMoreThanSmallestPack, biggerPackSize, totalQuantityWithSmallerPacks, smallestPackSize)

			if isDiffMoreThanSmallestPack {
				smallerPacksQuantity += cq.Quantity
				continue
			}

			//if the difference is indeed smaller, for example 20 (500-480), and smallest pack size is 80
			//	then we can rearrange.

			//add one pack of a bigger size
			logrus.Infof("adding pack; size %d, quantity: %d", biggerPackSize, totalQuantityWithSmallerPacks)
			packMap.AddPack(biggerPackSize, totalQuantityWithSmallerPacks)

			//clean up all smaller packs
			for k := i; k >= 0; k-- {
				logrus.Infof("252: setting count of pack sizes %d to 0", packSizesAsc[k])
				packMap.SetCount(packSizesAsc[k], 0)
			}
			smallerPacksQuantity = 0
		}

	}
}

// the packsToResponse method
func (s *PackServiceImpl) packsToResponse(packMap model.Packs) []model.PackDetails {
	var packsNeeded []model.PackDetails

	// iterate over the original slice to maintain the order of the response
	for _, packSize := range s.packSizes {
		if cq, ok := packMap[packSize]; ok {
			packsNeeded = append(packsNeeded, model.PackDetails{
				PackSize:   packSize,
				PacksCount: cq.Count,
			})
		}
	}

	return packsNeeded
}
