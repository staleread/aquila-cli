package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

var CLI struct {
	Gen struct {
		Name string `short:"n" help:"Optional name for the key pair."`
	} `cmd:"" help:"Generate a new key pair."`

	Enc struct {
		Input  string `short:"i" required:"" type:"existingfile" help:"Path to the input file."`
		Output string `short:"o" required:"" type:"path" help:"Path to the output file."`
		Key    string `short:"k" required:"" type:"existingfile" help:"Path to the public key file."`
	} `cmd:"" help:"Encrypt a file."`

	Dec struct {
		Input  string `short:"i" required:"" type:"existingfile" help:"Path to the input file."`
		Output string `short:"o" required:"" type:"path" help:"Path to the output file."`
		Key    string `short:"k" required:"" type:"existingfile" help:"Path to the private key file."`
	} `cmd:"" help:"Decrypt a file."`

	Config struct{} `cmd:"" help:"Show cipher configuration."`
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "gen":
		if err := generateKeyPair(CLI.Gen.Name); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating keys: %v\n", err)
			os.Exit(1)
		}
	case "enc":
		if err := encryptFile(CLI.Enc.Input, CLI.Enc.Output, CLI.Enc.Key); err != nil {
			fmt.Fprintf(os.Stderr, "Error encrypting file: %v\n", err)
			os.Exit(1)
		}
	case "dec":
		if err := decryptFile(CLI.Dec.Input, CLI.Dec.Output, CLI.Dec.Key); err != nil {
			fmt.Fprintf(os.Stderr, "Error decrypting file: %v\n", err)
			os.Exit(1)
		}
	case "config":
		showConfig()
	}
}
