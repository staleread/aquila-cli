package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/staleread/aquila/asym"
)

func exportANF(keyPath string) error {
	keyF, err := os.Open(keyPath)
	if err != nil {
		return fmt.Errorf("failed to open public key file: %w", err)
	}
	defer keyF.Close()

	pub := &asym.PublicKey{}
	if err := pub.Decode(keyF); err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	randInput := make([]byte, asym.BlockSize)
	if _, err := rand.Read(randInput); err != nil {
		return fmt.Errorf("failed to generate random input block: %w", err)
	}

	writer := bufio.NewWriter(os.Stdout)
	if err := pub.ExportToANF(writer, randInput); err != nil {
		return fmt.Errorf("failed to export ANF: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	return nil
}
