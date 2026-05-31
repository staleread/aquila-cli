package main

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"io"
	"math/bits"
	"os"
	"path/filepath"
	"strconv"

	"github.com/staleread/aquila/sym"
)

func runAnalyze(folder string, rndSamples int) error {
	block, err := sym.New(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate block cipher: %w", err)
	}

	blockSizeBytes := block.BlockSize()
	blockSizeBits := blockSizeBytes * 8

	if err := os.MkdirAll(folder, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", folder, err)
	}

	avFile, err := os.Create(filepath.Join(folder, "avalanche.csv"))
	if err != nil {
		return fmt.Errorf("failed to create avalanche.csv: %w", err)
	}
	defer avFile.Close()

	writer := csv.NewWriter(avFile)
	defer writer.Flush()

	if err := writer.Write([]string{"bit_index", "zeros_bg", "ones_bg", "random_bg_avg"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	hammingDistance := func(a, b []byte) int {
		dist := 0
		for i := range a {
			dist += bits.OnesCount8(a[i] ^ b[i])
		}
		return dist
	}

	zeroBaseSrc := make([]byte, blockSizeBytes)
	zeroBaseDst := make([]byte, blockSizeBytes)
	block.Encrypt(zeroBaseDst, zeroBaseSrc)

	oneBaseSrc := make([]byte, blockSizeBytes)
	for j := range oneBaseSrc {
		oneBaseSrc[j] = 0xFF
	}
	oneBaseDst := make([]byte, blockSizeBytes)
	block.Encrypt(oneBaseDst, oneBaseSrc)

	randSamplesSrc := make([][]byte, rndSamples)
	randSamplesDst := make([][]byte, rndSamples)
	for s := range rndSamples {
		randSamplesSrc[s] = make([]byte, blockSizeBytes)
		if _, err := io.ReadFull(rand.Reader, randSamplesSrc[s]); err != nil {
			return fmt.Errorf("failed to read random bytes: %w", err)
		}
		randSamplesDst[s] = make([]byte, blockSizeBytes)
		block.Encrypt(randSamplesDst[s], randSamplesSrc[s])
	}

	srcBuf := make([]byte, blockSizeBytes)
	dstBuf := make([]byte, blockSizeBytes)

	for i := range blockSizeBits {
		// --- Zeros BG ---
		for j := range srcBuf {
			srcBuf[j] = 0
		}
		srcBuf[i/8] = 1 << (i % 8)
		block.Encrypt(dstBuf, srcBuf)
		zerosBgVal := hammingDistance(dstBuf, zeroBaseDst)

		// --- Ones Background ---
		for j := range srcBuf {
			srcBuf[j] = 0xFF
		}
		srcBuf[i/8] &^= (1 << (i % 8))
		block.Encrypt(dstBuf, srcBuf)
		onesBgVal := hammingDistance(dstBuf, oneBaseDst)

		// --- Random Background Avg ---
		totalRandDiff := 0
		for s := range rndSamples {
			copy(srcBuf, randSamplesSrc[s])
			srcBuf[i/8] ^= (1 << (i % 8))
			block.Encrypt(dstBuf, srcBuf)
			totalRandDiff += hammingDistance(dstBuf, randSamplesDst[s])
		}
		randBgAvgVal := float64(totalRandDiff) / float64(rndSamples)

		row := []string{
			strconv.Itoa(i),
			strconv.Itoa(zerosBgVal),
			strconv.Itoa(onesBgVal),
			fmt.Sprintf("%.4f", randBgAvgVal),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row %d to avalanche.csv: %w", i, err)
		}
	}

	// 2. Create correlation.csv
	corrFile, err := os.Create(filepath.Join(folder, "correlation.csv"))
	if err != nil {
		return fmt.Errorf("failed to create correlation.csv: %w", err)
	}
	corrFile.Close()

	return nil
}
