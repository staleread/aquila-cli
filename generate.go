package main

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/staleread/aquila/asym"
)

func generateKeyPair(name string) error {
	cfgID := getConfigID()
	privFile := fmt.Sprintf("id_aquila%s", cfgID)
	pubFile := fmt.Sprintf("id_aquila%s.pub", cfgID)

	if name != "" {
		privFile = fmt.Sprintf("id_aquila%s_%s", cfgID, name)
		pubFile = fmt.Sprintf("id_aquila%s_%s.pub", cfgID, name)
	}

	priv, pub, err := asym.GenerateKeyPair(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	privF, err := os.OpenFile(privFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open private key file: %w", err)
	}
	defer privF.Close()

	if err := priv.Encode(privF); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	pubF, err := os.OpenFile(pubFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open public key file: %w", err)
	}
	defer pubF.Close()

	if err := pub.Encode(pubF); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	fmt.Printf("Successfully generated key pair:\n  Private key: %s\n  Public key:  %s\n", privFile, pubFile)
	return nil
}
