package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/staleread/aquila/asym"
)

func exportANF(keyPath, outputPath string) error {
	keyF, err := os.Open(keyPath)
	if err != nil {
		return fmt.Errorf("failed to open public key file: %w", err)
	}
	defer keyF.Close()

	pub, err := asym.DecodePublicKey(keyF)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	randInput := make([]byte, asym.BlockSize)
	if _, err := rand.Read(randInput); err != nil {
		return fmt.Errorf("failed to generate random input block: %w", err)
	}

	outF, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outF.Close()

	writer := bufio.NewWriter(outF)
	if err := pub.ExportToANF(writer, randInput); err != nil {
		return fmt.Errorf("failed to export ANF: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	fmt.Printf("Successfully exported ANF to %s\n", outputPath)
	return nil
}
