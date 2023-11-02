package main

import "fmt"

func lfsr(seed int, length int) []int {
	sequence := make([]int, length)
	lfsrValue := seed

	for i := 0; i < length; i++ {
		sequence[i] = lfsrValue
		feedback := lfsrValue & 1
		lfsrValue >>= 1
		if feedback&1 == 1 {
			lfsrValue ^= 0x110
		}
	}

	return sequence
}

func main() {
	seed := 222 // Initial seed value (between 0 and 511)
	length := 276

	sequence := lfsr(seed, length)

	fmt.Println("Generated pseudo-random sequence with 511 numbers:")
	for _, value := range sequence {
		fmt.Printf("%d ", value%256)
	}
}
