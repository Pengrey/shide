package main

import (
	"fmt"
	"log"
)

func obfuscateSHB(data []byte) []byte {
	// Apply JenkinsOneAtATime32Bit hash to shellcode string on each byte
	var hashedData []byte

	// Iterate over each byte
	for _, b := range data {
		// Convert byte to string
		bString := fmt.Sprintf("%d", b)

		// Iterate over each character in the byte string
		for _, c := range bString {
			// Convert character to integer
			cInt := int(c)

			// Apply JenkinsOneAtATime32Bit hash
			hashedData = append(hashedData, byte(cInt))
		}
	}
	return hashedData
}

func obfuscateBinary(data []byte, otype string) []byte {
	// Obfuscate binary file
	switch otype {
	case "SHB":
		return obfuscateSHB(data)
	default:
		log.Fatal("Invalid obfuscation type")
		return nil
	}
}

func getDeobfuscationStubC(data []byte, otype string) string {
	// Get deobfuscation stub in C
	var deobfuscationStub string
	switch otype {
	case "SHB":
		deobfuscationStub += `#include <Windows.h>`

		// Save data to a variable
		deobfuscationStub += "\n\nunsigned char data[] = {"
		for i, b := range data {
			if i == len(data)-1 {
				deobfuscationStub += fmt.Sprintf("0x%x", b)
			} else {
				deobfuscationStub += fmt.Sprintf("0x%x, ", b)
			}
		}
		deobfuscationStub += "};"

	default:
		log.Fatal("Invalid compression type")
		return ""
	}
	return deobfuscationStub
}

func getDeobfuscationStub(language string, data []byte, ctype string) string {
	// Get decompression stub
	switch language {
	case "C":
		return getDeobfuscationStubC(data, ctype)
	default:
		log.Fatal("Language not supported")
		return ""
	}
}
