package main

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/staleread/aquila/asym"
)

func decryptFile(inputPath, outputPath, keyPath string) error {
	keyF, err := os.Open(keyPath)
	if err != nil {
		return fmt.Errorf("failed to open private key file: %w", err)
	}
	defer keyF.Close()

	priv := &asym.PrivateKey{}
	if err := priv.Decode(keyF); err != nil {
		return fmt.Errorf("failed to decode private key: %w", err)
	}

	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	if len(ciphertext)%asym.BlockSize != 0 {
		return fmt.Errorf("ciphertext length is not a multiple of the block size")
	}

	paddedPlaintext, err := priv.Decrypt(rand.Reader, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	plaintext, err := pkcs7Unpad(paddedPlaintext, asym.BlockSize)
	if err != nil {
		return fmt.Errorf("failed to unpad decrypted data: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Successfully decrypted %s to %s\n", inputPath, outputPath)
	return nil
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("invalid padding size")
	}

	unpadding := int(data[length-1])
	if unpadding > blockSize || unpadding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	padtext := data[length-unpadding:]
	for _, b := range padtext {
		if int(b) != unpadding {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return data[:(length - unpadding)], nil
}
