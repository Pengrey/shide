package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var SEED int64 = 0xdeadbeef

func main() {
	app := &cli.App{
		Name:  "shide",
		Usage: "Obfuscate, compress and parse shellcode",

		Commands: []*cli.Command{
			{
				Name:        "parse",
				Usage:       "parse [--language C|Go|Rust] [--bin BINARY_FILE] [--cols COLUMNS](optional) [--out OUTPUT_FILE](optional)",
				Description: "Parse a binary file into a language specific array",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "language",
						Aliases:  []string{"l"},
						Usage:    "Language to parse the binary file into",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "bin",
						Aliases:  []string{"b"},
						Usage:    "Binary file to parse",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "out",
						Aliases: []string{"o"},
						Usage:   "Output file",
					},
					&cli.IntFlag{
						Name:    "cols",
						Aliases: []string{"c"},
						Usage:   "Number of columns to print the shellcode array",
					},
				},
				Action: func(c *cli.Context) error {
					// Read binary file
					data := readBinaryFile(c.String("bin"))

					// Get language to parse into
					language := c.String("language")

					// Get shellcode array
					shellcodeString := getShellCodeArray(language, data, c.Int("cols"))

					// Check if output file is specified
					if c.String("out") != "" {
						// Write shellcode string to file
						writeStringToFile(c.String("out"), shellcodeString)
					} else {
						// Print shellcode string to stdout
						fmt.Println(shellcodeString)
					}

					return nil
				},
			},
			{
				Name:        "compress",
				Usage:       "compress [--language C|Go|Rust] [--bin BINARY_FILE] [--type RLE|BWT] [--out OUTPUT_FILE]",
				Description: "Compress a binary file and generate a decompression stub",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "language",
						Aliases:  []string{"l"},
						Usage:    "Language to build the decompression stub in",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "bin",
						Aliases:  []string{"b"},
						Usage:    "Binary file to compress",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "out",
						Aliases:  []string{"o"},
						Usage:    "Output file for the decompression stub",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "type",
						Aliases: []string{"t"},
						Usage:   "Compression type to use",
					},
				},
				Action: func(c *cli.Context) error {
					// Read binary file
					data := readBinaryFile(c.String("bin"))

					// Get language to parse into
					language := c.String("language")

					// Get compression type
					compressionType := c.String("type")

					// Compress binary file
					compressedData := compressBinary(data, compressionType)

					// Get decompression stub
					decompressionStub := getDecompressionStub(language, compressedData, compressionType)

					// Write decompression stub to file
					writeStringToFile(c.String("out"), decompressionStub)

					return nil
				},
			},
			{
				Name:        "obfuscate",
				Usage:       "obfuscate [--language C|Go] [--bin BINARY_FILE] [--type RBM] [--out OUTPUT_FILE]",
				Description: "Obfuscate a binary file and generate a deobfuscation stub",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "language",
						Aliases:  []string{"l"},
						Usage:    "Language to build the deobfuscation stub in",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "bin",
						Aliases:  []string{"b"},
						Usage:    "Binary file to obfuscate",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "type",
						Aliases:  []string{"t"},
						Usage:    "Obfuscation type to use",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "out",
						Aliases: []string{"o"},
						Usage:   "Output file for the deobfuscation stub",
					},
				},
				Action: func(c *cli.Context) error {
					// Read binary file
					data := readBinaryFile(c.String("bin"))

					// Get language to parse into
					language := c.String("language")

					// Get obfuscation type
					obfuscationType := c.String("type")

					// Obfuscate binary file
					obfuscatedData := obfuscateBinary(data, obfuscationType)

					// Get deobfuscation stub
					deobfuscationStub := getDeobfuscationStub(language, obfuscatedData, obfuscationType)

					// Check if output file is specified
					if c.String("out") != "" {
						// Write deobfuscation stub to file
						writeStringToFile(c.String("out"), deobfuscationStub)
					} else {
						// Print deobfuscation stub to stdout
						fmt.Println(deobfuscationStub)
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
