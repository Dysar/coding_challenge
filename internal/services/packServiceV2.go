package services

import (
	"challenge/internal/model"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"slices"
)

type (
	PackServiceV2 interface {
		CalculatePacks(orderQuantity int) ([]model.PackDetails, error)
	}

	PackServiceImplV2 struct {
		packSizes []int
	}
)

func NewPackServiceV2() *PackServiceImplV2 {
	initialPackSizes := []int{23, 31, 53}
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

	var x1, x2, x3, sum int
	for size, count := range packMap {
		sum += count.Count * size
		switch size {
		case 53:
			x1 = count.Count
		case 31:
			x2 = count.Count
		case 23:
			x3 = count.Count
		}
	}
	//x1, x2, x3 = calculate(x1, x2, x3, orderQuantity, sum)
	//result := []model.PackDetails{
	//	{PackSize: 53, PacksCount: x1},
	//	{PackSize: 31, PacksCount: x2},
	//	{PackSize: 23, PacksCount: x3},
	//}

	c := newCalculator([]int{x1, x2, x3}, []int{53, 31, 23}, orderQuantity, sum)
	values := c.flexibleCalculate()
	result := []model.PackDetails{
		{PackSize: 53, PacksCount: values[0]},
		{PackSize: 31, PacksCount: values[1]},
		{PackSize: 23, PacksCount: values[2]},
	}
	return result, nil
}

func calculate(x1, x2, x3, targetCount, initialCount int) (int, int, int) {

	bestSumSoFar := 10000000
	bestCountSoFarElements := [3]int{x1, x2, x3}
	sum := func(x1, x2, x3 int) int {

		sum := x1*53 + x2*31 + x3*23
		fmt.Printf("x1: %d, x2: %d, x3: %d, sum: %d\n", x1, x2, x3, sum)
		if sum == targetCount {
			return sum
		}

		if sum >= targetCount && sum < initialCount && sum < bestSumSoFar {
			fmt.Printf("target: %d, initial: %d best count so far: %d\n", targetCount, initialCount, sum)
			bestSumSoFar = sum
			bestCountSoFarElements = [3]int{x1, x2, x3}
		}
		return sum
	}

	for ; x1 >= 0; x1-- {

		fmt.Println("-----iteration-----")
		x1sum := sum(x1, x2, x3)

		//fmt.Printf("x1: %d, x2: %d, sum:%d \n", x1, x2, x1sum)
		//fmt.Println("current is not exactly target? ", x1sum != targetCount)
		if x1sum != targetCount {
			for ; x2 >= 0; x2++ {
				if x2sum := sum(x1, x2, x3); x2sum < targetCount {

					//fmt.Printf("playing with x2; x1: %d, x2: %d, sum:%d \n", x1, x2, x2sum)
					//try to adjust the sum with x3
					for ; x3 >= 0; x3++ {
						if x3sum := sum(x1, x2, x3); x3sum < targetCount {
							//fmt.Printf("playing with x3; x1: %d, x2: %d, x3: %d, sum:%d \n", x1, x2, x3, x3sum)

						} else if x3sum == targetCount {
							return x1, x2, x3
						} else {
							//fmt.Println("x3 sum", x3sum)
							x3 = 0
							break
						}
					}
				} else if x2sum == targetCount {
					return x1, x2, x3
				} else { //x2sum > target
					break
				}
			}
		} else {
			return x1, x2, x3
		}
	}

	return bestCountSoFarElements[0], bestCountSoFarElements[1], bestCountSoFarElements[2]
}

type calculator struct {
	bestSumSoFar              int
	bestCountSoFarElements    []int
	targetCount, initialCount int
	values, sizes             []int
}

func newCalculator(values, sizes []int, targetCount, initialCount int) *calculator {
	c := &calculator{
		bestSumSoFar:           10000000,
		bestCountSoFarElements: make([]int, len(values)),
		targetCount:            targetCount,
		initialCount:           initialCount,
		values:                 values,
		sizes:                  sizes,
	}
	copy(c.bestCountSoFarElements, values)
	return c
}

// flexibleCalculate accepts values, sizes (must be sorted desc)
func (s *calculator) flexibleCalculate() []int {

	for x1 := s.values[0]; x1 >= 0; x1-- {

		fmt.Println("summing up x1")
		elementsToSum := append([]int{x1}, s.values[1:]...)
		x1sum := s.sum(elementsToSum)

		//fmt.Printf("x1: %d, x2: %d, sum:%d \n", x1, x2, x1sum)
		//fmt.Println("current is not exactly target? ", x1sum != targetCount)
		if x1sum != s.targetCount {
			if result := s.iterateOver(1, x1); result != nil {
				return result
			}
		} else {
			return append([]int{x1}, s.values[1:]...)
		}
	}

	return s.bestCountSoFarElements
}

func (s *calculator) sum(values []int) int {
	var sum int
	for i, v := range values {
		sum += v * s.sizes[i]
		fmt.Printf("x%d: %d ", i+1, v)
	}
	fmt.Printf("sum:%d", sum)
	fmt.Printf("\n")

	if sum == s.targetCount {
		return sum
	}

	if sum >= s.targetCount && sum < s.initialCount && sum < s.bestSumSoFar {
		fmt.Printf("target: %d, initial: %d best count so far: %d\n", s.targetCount, s.initialCount, sum)
		s.bestSumSoFar = sum
		s.bestCountSoFarElements = values
	}
	return sum
}

// iterateOver returns the result if we got it and nil if we need to continue
func (s *calculator) iterateOver(i int, prevValues ...int) []int {
	val := s.values[i] //TODO: check if the element exists

	for x := val; x < 1000000 && x >= 0; x++ {

		fmt.Printf("summing up x%d\n", i+1)
		elementsToSum := append(prevValues, x)
		elementsToSum = append(elementsToSum, s.values[i+1:]...)
		//fmt.Printf("i:%d elements to sum: %v\n", i, elementsToSum)

		if xsum := s.sum(elementsToSum); xsum < s.targetCount {
			if i != len(s.values)-1 { //if is not the smallest item
				if result := s.iterateOver(i+1, append(prevValues, x)...); result != nil {
					return result
				} else {
					//break //should go back a step and do -1 for larger size
					continue
				}
			} else {

			}
		} else if xsum == s.targetCount {
			fmt.Printf("i:%d, res: %v\n", i, elementsToSum)
			return elementsToSum
		} else {
			return nil
		}
	}
	return nil
}
