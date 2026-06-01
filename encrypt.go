package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/staleread/aquila/asym"
)

func encryptFile(inputPath, outputPath, keyPath string) error {
	keyF, err := os.Open(keyPath)
	if err != nil {
		return fmt.Errorf("failed to open public key file: %w", err)
	}
	defer keyF.Close()

	pub := &asym.PublicKey{}
	if err := pub.Decode(keyF); err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	paddedData := pkcs7Pad(inputData, asym.BlockSize)

	ciphertext, err := pub.Encrypt(rand.Reader, paddedData)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Successfully encrypted %s to %s\n", inputPath, outputPath)
	return nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}
