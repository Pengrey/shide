package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
