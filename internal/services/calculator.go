package services

import (
	"fmt"
	"math"
)

// calculator is covered with tests in the packService code
type calculator struct {
	bestSumSoFar                 int
	bestCountSoFarElements       []int
	targetCount, initialCount    int
	values, sizes, elementsToSum []int
}

// newCalculator constructor. values, sizes (must be sorted desc)
func newCalculator(values, sizes []int, targetCount, initialCount int) *calculator {
	c := &calculator{
		bestSumSoFar:           int(math.Pow(10, 9)), //some large number
		bestCountSoFarElements: make([]int, len(values)),
		targetCount:            targetCount,
		initialCount:           initialCount,
		values:                 values,
		sizes:                  sizes,
		// to avoid unnecessary slice creation
		elementsToSum: make([]int, len(values)),
	}
	copy(c.bestCountSoFarElements, values)
	return c
}

func (s *calculator) flexibleCalculate() []int {

	for x1 := s.values[0]; x1 >= 0; x1-- {

		//fmt.Println("summing up x1")
		elementsToSum := append([]int{x1}, s.values[1:]...)
		x1sum := s.sum(elementsToSum)

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
		//fmt.Printf("x%d: %d ", i+1, v)
	}
	//fmt.Printf("sum:%d", sum)
	//fmt.Printf("\n")

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
	val := s.values[i]

	for x := val; x < 1000000 && x >= 0; x++ {

		//fmt.Printf("summing up x%d\n", i+1)
		s.elementsToSum = append(prevValues, x)
		s.elementsToSum = append(s.elementsToSum, s.values[i+1:]...)

		if xsum := s.sum(s.elementsToSum); xsum < s.targetCount {
			if i != len(s.values)-1 { //if is not the smallest item
				if result := s.iterateOver(i+1, append(prevValues, x)...); result != nil {
					return result
				} else {
					continue
				}
			}
		} else if xsum == s.targetCount {
			//fmt.Printf("i:%d, res: %v\n", i, s.elementsToSum)
			return s.elementsToSum
		} else {
			//fmt.Printf("xsum:%d > targetCount: %v\n", xsum, s.targetCount)
			return nil
		}
	}
	return nil
}
