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

	Build struct {
		ConfigID string `arg:"" help:"Configuration ID in format <block>c<comp>f<fold>d<deg>."`
	} `cmd:"" help:"Build the CLI with the specified cipher configuration."`

	Analyze struct {
		Folder     string `short:"o" required:"" type:"path" help:"Path to the output folder."`
		RndSamples int    `long:"rnd-samples" required:"" help:"Number of random samples to analyze."`
	} `cmd:"" name:"analyze" help:"Run cipher analysis experiments."`

	Anf struct {
		Key    string `short:"k" required:"" type:"existingfile" help:"Path to the public key file."`
		Output string `short:"o" required:"" type:"path" help:"Path to the output file."`
	} `cmd:"" help:"Export public key equations in Algebraic Normal Form (ANF)."`
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
	case "build <config-id>":
		if err := runBuild(CLI.Build.ConfigID); err != nil {
			fmt.Fprintf(os.Stderr, "Error building CLI: %v\n", err)
			os.Exit(1)
		}
	case "analyze":
		if err := runAnalyze(CLI.Analyze.Folder, CLI.Analyze.RndSamples); err != nil {
			fmt.Fprintf(os.Stderr, "Error running analysis: %v\n", err)
			os.Exit(1)
		}
	case "anf":
		if err := exportANF(CLI.Anf.Key, CLI.Anf.Output); err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting ANF: %v\n", err)
			os.Exit(1)
		}
	}
}
