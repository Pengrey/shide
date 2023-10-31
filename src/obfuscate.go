package main

import (
	"fmt"
	"log"
	"math/rand"
)

func generateRandomBytesMap(seed int64, size int) map[byte][]int {
	// Generate random bytes from a seed
	var randomBytes = make([]byte, size)

	// Set the seed
	rand.Seed(seed)

	// Generate random bytes
	for i := 0; i < size; i++ {
		randomBytes[i] = byte(rand.Intn(256))
	}

	// Generate a map of the appearance of each byte, the key is the byte and the value is an array of the indexes where the byte appears
	var randomBytesMap = make(map[byte][]int)

	// Iterate over the random bytes
	for i, b := range randomBytes {
		// Check if the byte is already in the map
		if _, ok := randomBytesMap[b]; ok {
			// Append the index to the array
			randomBytesMap[b] = append(randomBytesMap[b], i)
		} else {
			// Create a new array with the index
			randomBytesMap[b] = []int{i}
		}
	}

	return randomBytesMap
}

func obfuscateRBM(data []byte) []int {
	// Obfuscate data using Runtime Bload Mapping

	// Create obfuscated array with the same size as the original array
	var obfuscatedData = make([]int, len(data))

	// Generate an array of 200mb of random bytes from 0x00 to 0xff, originated from the hardcoded seed
	var randomBytesMap = generateRandomBytesMap(SEED, 200*1024*1024)

	// Iterate over the original array
	for i, b := range data {
		// Get the array of indexes where the byte appears
		var indexes = randomBytesMap[b]

		// Get a random index from the array
		var randomIndex = indexes[rand.Intn(len(indexes))]
		obfuscatedData[i] = randomIndex
	}

	return obfuscatedData
}

func obfuscateBinary(data []byte, otype string) []int {
	// Obfuscate binary file
	switch otype {
	case "RBM":
		return obfuscateRBM(data)
	default:
		log.Fatal("Invalid obfuscation type")
		return nil
	}
}

func getDeobfuscationStubC(data []int, otype string) string {
	// Get deobfuscation stub in C
	var deobfuscationStub string
	switch otype {
	case "RBM":
		deobfuscationStub += `#include <Windows.h>`

		// Create array of size 200mb
		var size = 200 * 1024 * 1024

		// Set the seed
		rand.Seed(SEED)

		// Save the random bytes to file in C format
		deobfuscationStub += "\n\n// Random bytes array\n"
		deobfuscationStub += "BYTE randomBytes[] = {"
		for i := 0; i < size; i++ {
			deobfuscationStub += fmt.Sprintf("0x%02x, ", byte(rand.Intn(256)))
		}
		deobfuscationStub += "};\n"

	default:
		log.Fatal("Invalid compression type")
		return ""
	}
	return deobfuscationStub
}

func getDeobfuscationStub(language string, data []int, ctype string) string {
	// Get decompression stub
	switch language {
	case "C":
		return getDeobfuscationStubC(data, ctype)
	default:
		log.Fatal("Language not supported")
		return ""
	}
}
