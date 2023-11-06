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
		deobfuscationStub += `#include <Windows.h>
#include <stdio.h>

// LFSR function
unsigned int* lfsr(unsigned int  seed, unsigned int length) {
	unsigned int* sequence = (unsigned int*)malloc(sizeof(unsigned int) * length);
	unsigned int lfsrValue = seed;

	for (int i = 0; i < length; i++) {
		sequence[i] = lfsrValue%256;
		unsigned int feedback = lfsrValue & 1;
		lfsrValue >>= 1;
		if (feedback & 1) {
			lfsrValue ^= 0x110;
		}
	}

	return sequence;
}

// Function to deofuscate LXB payload
BOOL deofuscateLXB(IN int seed, IN INT* data, IN SIZE_T sSize, OUT PBYTE* pDeobfuscatedPayload) {
	BYTE* payload = NULL;

	// Allocate memory for deofuscated payload
	payload = (BYTE*)VirtualAlloc(NULL, sSize, MEM_COMMIT | MEM_RESERVE, PAGE_READWRITE);
	if (payload == NULL) {
		printf("[!] VirtualAlloc Failed With Error : %d \n", GetLastError());
		return;
	}

	// Deofuscate payload

	// Generate LFSR sequence with length of the payload and the key
	int* sequence = lfsr(seed, sSize);

	// Deofuscate payload
	for (int i = 0; i < sSize; i++) {
		// Check if the current integer is negative
		if (data[i] < 0) {
			// If it is negative, turn positive and add it to the payload
			payload[i] = (BYTE)(data[i] * -1);
		}
		else {
			// If it is positive, add it to the payload
			payload[i] = (BYTE)(sequence[data[i]] % 256);
		}
	}

	// Set deofuscated payload
	*pDeobfuscatedPayload = payload;

	return;
}`

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
		deobfuscationStub += fmt.Sprintf("};\n\nint seed = %d;\n\n", seed)

		deobfuscationStub += `int main() {
	PBYTE pDeobfuscatedPayload = NULL;
	SIZE_T sDeobfuscatedSize = sizeof(data) / sizeof(int);

	// deofuscate payload
	deofuscateLXB(seed, data, sDeobfuscatedSize, &pDeobfuscatedPayload);

	// Print deofuscated in hex form
	for (int i = 0; i < sDeobfuscatedSize; i++) {
		printf("%02X ", pDeobfuscatedPayload[i]);
	}

	return 0;
}`

	default:
		log.Fatal("Invalid compression type")
		return ""
	}
	return deobfuscationStub
}

func getDeobfuscationStubGo(seed int, data []int, otype string) string {
	// Get deobfuscation stub in Go
	var deobfuscationStub string
	switch otype {
	case "LSM":
		deobfuscationStub += `package main

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

func main() {`

		deobfuscationStub += fmt.Sprintf("\n\tdata := []int{")

		for i, value := range data {
			if i == len(data)-1 {
				deobfuscationStub += fmt.Sprintf("%d", value)
			} else {
				deobfuscationStub += fmt.Sprintf("%d, ", value)
			}
		}

		deobfuscationStub += fmt.Sprintf("}\n\tseed := %d\n\n", seed)

		deobfuscationStub += `	sequence := lfsr(seed, len(data))

	// xor the data with the sequence
	for i := range data {
		if data[i] > 0 {
			data[i] = sequence[data[i]] % 256
		} else {
			data[i] = data[i] * -1
		}
	}

	// Print hex values
	for _, b := range data {
		fmt.Printf("%02x ", b)
	}
}`
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
	case "Go":
		return getDeobfuscationStubGo(seed, data, ctype)
	default:
		log.Fatal("Language not supported")
		return ""
	}
}
