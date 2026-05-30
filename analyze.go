package main

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"math/bits"
	"os"
	"strconv"

	"github.com/staleread/aquila/sym"
)

func runAnalyzeBoundaryZeros(outputPath string) error {
	block, err := sym.New(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate block cipher: %w", err)
	}

	blockSizeBytes := block.BlockSize()
	blockSizeBits := blockSizeBytes * 8

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	if err := writer.Write([]string{"OneInputPosition", "OnesOutputCount"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	countOnes := func(b []byte) int {
		ones := 0
		for _, val := range b {
			ones += bits.OnesCount8(val)
		}
		return ones
	}

	src := make([]byte, blockSizeBytes)
	dst := make([]byte, blockSizeBytes)
	block.Encrypt(dst, src)
	onesCount := countOnes(dst)

	if err := writer.Write([]string{"-1", strconv.Itoa(onesCount)}); err != nil {
		return fmt.Errorf("failed to write zeroed block result: %w", err)
	}

	for i := range blockSizeBits {
		for j := range src {
			src[j] = 0
		}
		byteIndex := i / 8
		bitIndex := i % 8
		src[byteIndex] = 1 << bitIndex

		block.Encrypt(dst, src)
		onesCount = countOnes(dst)

		if err := writer.Write([]string{strconv.Itoa(i), strconv.Itoa(onesCount)}); err != nil {
			return fmt.Errorf("failed to write result for bit %d: %w", i, err)
		}
	}

	return nil
}

func runAnalyzeBoundaryOnes(outputPath string) error {
	block, err := sym.New(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate block cipher: %w", err)
	}

	blockSizeBytes := block.BlockSize()
	blockSizeBits := blockSizeBytes * 8

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	if err := writer.Write([]string{"ZeroInputPosition", "OnesOutputCount"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	countOnes := func(b []byte) int {
		ones := 0
		for _, val := range b {
			ones += bits.OnesCount8(val)
		}
		return ones
	}

	src := make([]byte, blockSizeBytes)
	for j := range src {
		src[j] = 0xFF
	}
	dst := make([]byte, blockSizeBytes)
	block.Encrypt(dst, src)
	onesCount := countOnes(dst)

	if err := writer.Write([]string{"-1", strconv.Itoa(onesCount)}); err != nil {
		return fmt.Errorf("failed to write all-ones block result: %w", err)
	}

	for i := range blockSizeBits {
		for j := range src {
			src[j] = 0xFF
		}
		byteIndex := i / 8
		bitIndex := i % 8
		src[byteIndex] &^= (1 << bitIndex)

		block.Encrypt(dst, src)
		onesCount = countOnes(dst)

		if err := writer.Write([]string{strconv.Itoa(i), strconv.Itoa(onesCount)}); err != nil {
			return fmt.Errorf("failed to write result for bit %d: %w", i, err)
		}
	}

	return nil
}
