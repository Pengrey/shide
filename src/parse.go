package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func readBinaryFile(filename string) []byte {
	// Open file for reading
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Get file size
	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	// Read file into byte array
	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func getShellCodeArrayC(data []byte, cols int) string {
	// Get shellcode array in C
	shellcodeString := "unsigned char Payload[] =\n{\n\t\""
	for i, b := range data {
		// If the byte is only one digit, add a 0 in front of it for the format \x0x
		if b < 16 {
			shellcodeString += fmt.Sprintf("\\x0%x", b)
		} else {
			shellcodeString += fmt.Sprintf("\\x%x", b)
		}
		// Check if we need to print a newline
		if (i+1)%cols == 0 {
			shellcodeString += "\"\n\t\""
		}
	}

	shellcodeString += "\"\n};\n"
	return shellcodeString
}

func getShellCodeArrayGo(data []byte, cols int) string {
	// Get shellcode array in Go
	shellcodeString := "var payload = []byte{\n"
	for i, b := range data {
		// If the byte is only one digit, add a 0 in front of it
		if b < 16 {
			shellcodeString += fmt.Sprintf("0x0%x,", b)
		} else {
			shellcodeString += fmt.Sprintf("0x%x,", b)
		}
		// Check if we need to print a newline
		if (i+1)%cols == 0 {
			shellcodeString += "\n"
		}
	}
	// Remove last comma
	shellcodeString = shellcodeString[:len(shellcodeString)-1]

	shellcodeString += "\n}\n"
	return shellcodeString
}

func getShellCodeArrayRust(data []byte, cols int) string {
	// Get shellcode array in Rust
	shellcodeString := "let payload: [u8; " + strconv.Itoa(len(data)) + "] = [\n"
	for i, b := range data {
		// If the byte is only one digit, add a 0 in front of it
		if b < 16 {
			shellcodeString += fmt.Sprintf("0x0%x,", b)
		} else {
			shellcodeString += fmt.Sprintf("0x%x,", b)
		}
		// Check if we need to print a newline
		if (i+1)%cols == 0 {
			shellcodeString += "\n"
		}
	}
	// Remove last comma
	shellcodeString = shellcodeString[:len(shellcodeString)-1]

	shellcodeString += "];\n"
	return shellcodeString
}

func getShellCodeArray(language string, data []byte, cols int) string {
	// Check if cols is specified
	if cols == 0 {
		// Set default cols
		cols = 14
	}

	// Get shellcode array
	switch language {
	case "C":
		return getShellCodeArrayC(data, cols)
	case "Go":
		return getShellCodeArrayGo(data, cols)
	case "Rust":
		return getShellCodeArrayRust(data, cols)
	default:
		log.Fatal("Language not supported")
		return ""
	}
}

func writeStringToFile(filename string, data string) {
	// Open file for writing
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Write string to file
	_, err = file.WriteString(data)
	if err != nil {
		log.Fatal(err)
	}
}
