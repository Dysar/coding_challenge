package services

import (
	"challenge/internal/model"
	"errors"
	"fmt"
	"slices"
)

type (
	PackService interface {
		CalculatePacks(orderQuantity int) ([]model.PackDetails, error)
		UpdatePackSizes(packSizes []int)
		ReadPackSizes() []int
	}

	PackServiceImpl struct {
		// pack sizes are sorted desc
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
		fmt.Printf("packsCount(%d) := remainingQuantity(%d) / packSize(%d)", packsCount, remainingQuantity, packSize)
		if packsCount > 0 {
			fmt.Printf("adding a pack size %d, count: %d", packSize, packsCount)
			packMap.AddPacks(packSize, packsCount)
			remainingQuantity %= packSize
		}

		if remainingQuantity == 0 {
			continue
		}

		//if there is any remaining quantity when all the pack sizes have been checked
		if remainingQuantity <= smallestPack {

			//fill one smallest pack (250)
			fmt.Printf("adding one smallest pack size %d, count: %d", smallestPack, 1)
			packMap.AddPack(smallestPack)
			break
		}
	}

	var sum int
	values := make([]int, len(s.packSizes))

	//TODO: unnecessary iteration?
	for i, packSize := range s.packSizes {
		if count, ok := packMap[packSize]; ok {
			sum += packSize * count
			values[i] = count
		}
	}

	// Within the constraints of Rule 1 above, send out no more items than necessary to fulfil the order.
	c := newCalculator(values, s.packSizes, orderQuantity, sum)
	values = c.flexibleCalculate()

	for i, v := range values {
		size := s.packSizes[i]
		packMap[size] = v
	}

	// Within the constraints above, send out Rules 1 & 2 send out as few packs as possible to fulfil each order.
	s.adjustPacks(packMap)

	return s.packsToResponse(packMap), nil
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

func (s *PackServiceImpl) adjustPacks(packMap model.Packs) {

	packSizesAsc := slices.Clone(s.packSizes)
	slices.Reverse(packSizesAsc)

	var smallerPacksQuantity int
	smallestPackSize := packSizesAsc[0]

	for i := 0; i < len(packSizesAsc); i++ {
		packSize := packSizesAsc[i]

		count, ok := packMap[packSize]
		if !ok {
			continue
		}
		// add only positive count items to the response
		if count <= 0 {
			delete(packMap, packSize)
		}

		if i == len(packSizesAsc)-1 {
			//if it's the biggest pack, we still might want to make adjustments
			fmt.Printf("checking biggest packs; smallerPacksQuantity: %d\n", smallerPacksQuantity)
			if smallerPacksQuantity > packSize {
				break
			}
			//smallerPacksQuantity <= packSize
			if isDiffLessThanSmallestPack := (packSize - smallerPacksQuantity) < smallestPackSize; isDiffLessThanSmallestPack {
				fmt.Printf("adding biggest pack %d with q:%d\n", packSize, smallerPacksQuantity)
				packMap.AddPack(packSize)
				for k := i - 1; k >= 0; k-- {
					fmt.Printf("setting count of pack sizes %d to 0\n", packSizesAsc[k])
					packMap.SetCount(packSizesAsc[k], 0)
				}
			}
		}

		if count <= 1 || i == len(packSizesAsc)-1 {
			smallerPacksQuantity += count * packSize
			continue
		}

		//there is more than 1 pack they can potentially be replaced by a larger pack

		fmt.Printf("count: %d, smallerPacksQuantity: %d\n", count, smallerPacksQuantity)

		totalQuantityWithSmallerPacks := count*packSize + smallerPacksQuantity
		fmt.Printf("totalQuantityWithSmallerPacks: %d\n", totalQuantityWithSmallerPacks)

		biggerPackSize := packSizesAsc[i+1]

		//it can happen that the total quantity is < or > or = than the bigger pack
		//if it's <, we need to make sure that the difference is less than the smallest pack

		if biggerPackSize < totalQuantityWithSmallerPacks {
			fmt.Printf("biggerPackSize(%d) is less than total quantity\n", biggerPackSize)
			//what if the bigger pack is 500 and total quantity of smaller packs is 2000?
			//	in that case I would look for bigger pack that could fit those items. next pack would be 1000, 1000<2000, next pack would be 2000, that is fine, we can use that one
			//	2000-2000=0; 0<250 (smallest pack), all good
			smallerPacksQuantity += count * packSize
			continue
		} else {
			fmt.Printf("biggerPackSize(%d) is greater or equal to the total quantity\n", biggerPackSize)

			//biggerPackSize => total quantity, e.g. pack size = 500, total quantity 100*4+80=480
			isDiffMoreThanSmallestPack := (biggerPackSize - totalQuantityWithSmallerPacks) > smallestPackSize
			fmt.Printf("isDiffMoreThanSmallestPack: %t, biggerPackSize: %d, totalQuantityWithSmallerPacks: %d, smallestPackSize:%d",
				isDiffMoreThanSmallestPack, biggerPackSize, totalQuantityWithSmallerPacks, smallestPackSize)

			if isDiffMoreThanSmallestPack {
				smallerPacksQuantity += count * packSize
				continue
			}

			//if the difference is indeed smaller, for example 20 (500-480), and smallest pack size is 80
			//	then we can rearrange.

			//add one pack of a bigger size
			fmt.Printf("adding pack; size %d, quantity: %d\n", biggerPackSize, totalQuantityWithSmallerPacks)
			packMap.AddPack(biggerPackSize)

			//clean up all smaller packs
			for k := i; k >= 0; k-- {
				fmt.Printf("setting count of pack sizes %d to 0\n", packSizesAsc[k])
				packMap.SetCount(packSizesAsc[k], 0)
			}
			smallerPacksQuantity = 0
		}

	}
}

func (s *PackServiceImpl) packsToResponse(packMap model.Packs) []model.PackDetails {
	var packsNeeded []model.PackDetails

	// iterate over the original slice to maintain the order of the response
	for _, packSize := range s.packSizes {
		if count, ok := packMap[packSize]; ok {
			packsNeeded = append(packsNeeded, model.PackDetails{
				PackSize:   packSize,
				PacksCount: count,
			})
		}
	}

	return packsNeeded
}
