package main

import (
	"fmt"
)

func main() {
	//fiveHundredK()
	hundredAndSeven()
}

func fiveHundredK() {
	// Define the sizes of the packs
	//packSizes := []int{53, 31, 23}

	// Define the target number of items to distribute
	targetItems := 500000

	// Set initial values for x1, x2, and x3
	x1 := 9434
	x2 := 0
	x3 := 0

	calculate(x1, x2, x3, targetItems)
}

func calculate(x1, x2, x3, targetItems int) {

	sum := func(x1, x2, x3 int) int {
		fmt.Printf("x1: %d, x2: %d, x3: %d\n", x1, x2, x3)
		return x1*53 + x2*31 + x3*23
	}

	possibleSolution := [3]int{}
	for ; x1 >= 0; x1-- {

		//fmt.Println("-----iteration-----")
		x1sum := sum(x1, x2, x3)
		//fmt.Printf("x1: %d, x2: %d, sum:%d \n", x1, x2, x1sum)
		//fmt.Println("current is not exactly target? ", x1sum != targetItems)
		if x1sum != targetItems {
			for ; x2 >= 0; x2++ {
				if x2sum := sum(x1, x2, x3); x2sum < targetItems {
					//fmt.Printf("playing with x2; x1: %d, x2: %d, sum:%d \n", x1, x2, x2sum)
					//try to adjust the sum with x3
					for ; x3 >= 0; x3++ {
						if x3sum := sum(x1, x2, x3); x3sum < targetItems {
							//fmt.Printf("playing with x3; x1: %d, x2: %d, x3: %d, sum:%d \n", x1, x2, x3, x3sum)

						} else if x3sum == targetItems {
							possibleSolution = [3]int{x1, x2, x3}
							goto response
						} else {
							//fmt.Println("x3 sum", x3sum)
							x3 = 0
							break
						}
					}
				} else if x2sum == targetItems {
					possibleSolution = [3]int{x1, x2, x3}
					goto response
				} else { //x2sum > target
					break
				}
			}
		} else {
			possibleSolution = [3]int{x1, x2, x3}
			break
		}
	}

response:
	totalPacks := possibleSolution[0] + possibleSolution[1] + possibleSolution[2]
	fmt.Printf("%v solution; total packs: %d\n", possibleSolution, totalPacks)
}

func hundredAndSeven() {
	targetItems := 107
	x1 := 2
	x2 := 1
	x3 := 0
	calculate(x1, x2, x3, targetItems)

}
