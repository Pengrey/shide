package main

import (
	"fmt"
	"log"
	"math/rand"
)

func getRandomSeed() int {
	// Get random seed
	return rand.Intn(512)
}
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
func createLSFSRMap(seed int, length int) map[byte][]int {
	// Create LFSR map
	sequence := lfsr(seed, length)

	// Create a map where the key is the byte
	// and the value is an array of indexes where the byte appears
	lfsrMap := make(map[byte][]int)

	// Iterate over the sequence
	for i, value := range sequence {
		// If the byte is not in the map, add it
		if _, ok := lfsrMap[byte(value%256)]; !ok {
			lfsrMap[byte(value%256)] = []int{i}
		} else {
			// If the byte is already in the map, append the index
			lfsrMap[byte(value%256)] = append(lfsrMap[byte(value%256)], i)
		}
	}

	return lfsrMap
}

func obfuscateLSM(seed int, data []byte) []int {
	// Apply LFSR Shellcode Mapping
	var obfuscatedData []int

	// Create LFSR map
	lfsrMap := createLSFSRMap(seed, len(data))

	// Iterate over the data
	for _, b := range data {
		// Check if the byte is in the map
		if _, ok := lfsrMap[b]; ok {
			// If the byte is in the map, get the indexes where it appears
			indexes := lfsrMap[b]

			// Get a random index
			randomIndex := indexes[rand.Intn(len(indexes))]

			// Append the random index to the obfuscated data
			obfuscatedData = append(obfuscatedData, randomIndex)
		} else {
			// If the byte is not in the map
			// Append the byte to the obfuscated data in a negative form
			obfuscatedData = append(obfuscatedData, -int(b))
		}
	}

	return obfuscatedData
}

func obfuscateBinary(seed int, data []byte, otype string) []int {
	// Obfuscate binary file
	switch otype {
	case "LSM":
		return obfuscateLSM(seed, data)
	default:
		log.Fatal("Invalid obfuscation type")
		return nil
	}
}

func getDeobfuscationStubC(seed int, data []int, otype string) string {
	// Get deobfuscation stub in C
	var deobfuscationStub string
	switch otype {
	case "LSM":
		deobfuscationStub += `#include <Windows.h>`

		// Save data to a variable
		deobfuscationStub += "\n\nint data[] = {"
		for i, value := range data {
			if i == len(data)-1 {
				deobfuscationStub += fmt.Sprintf("%d", value)
			} else {
				deobfuscationStub += fmt.Sprintf("%d, ", value)
			}
		}

		// Add SEED
		deobfuscationStub += fmt.Sprintf("};\n\nunsigned int seed = %d;\n\n", seed)

		// TODO: Add deobfuscation stub

	default:
		log.Fatal("Invalid compression type")
		return ""
	}
	return deobfuscationStub
}

func getDeobfuscationStub(language string, seed int, data []int, ctype string) string {
	// Get decompression stub
	switch language {
	case "C":
		return getDeobfuscationStubC(seed, data, ctype)
	default:
		log.Fatal("Language not supported")
		return ""
	}
}
