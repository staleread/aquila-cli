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

	"github.com/staleread/aquila/asym"
)

func runAnalyze(folder string, rndSamples int, correlation bool) error {
	priv, err := asym.GeneratePrivateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	blockSizeBytes := asym.BlockSize
	blockSizeBits := blockSizeBytes * 8

	changeCounts := make([]int, blockSizeBits)

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
	zeroBaseDst, err := priv.Sign(nil, zeroBaseSrc, nil)
	if err != nil {
		return fmt.Errorf("failed to encrypt zero base: %w", err)
	}

	oneBaseSrc := make([]byte, blockSizeBytes)
	for j := range oneBaseSrc {
		oneBaseSrc[j] = 0xFF
	}
	oneBaseDst, err := priv.Sign(nil, oneBaseSrc, nil)
	if err != nil {
		return fmt.Errorf("failed to encrypt one base: %w", err)
	}

	randSamplesSrc := make([][]byte, rndSamples)
	randSamplesDst := make([][]byte, rndSamples)
	for s := range rndSamples {
		randSamplesSrc[s] = make([]byte, blockSizeBytes)
		if _, err := io.ReadFull(rand.Reader, randSamplesSrc[s]); err != nil {
			return fmt.Errorf("failed to read random bytes: %w", err)
		}
		randSamplesDst[s], err = priv.Sign(nil, randSamplesSrc[s], nil)
		if err != nil {
			return fmt.Errorf("failed to encrypt random sample: %w", err)
		}
	}

	srcBuf := make([]byte, blockSizeBytes)
	var dstBuf []byte

	for i := range blockSizeBits {
		// --- Zeros BG ---
		for j := range srcBuf {
			srcBuf[j] = 0
		}
		srcBuf[i/8] = 1 << (i % 8)
		dstBuf, err = priv.Sign(nil, srcBuf, nil)
		if err != nil {
			return fmt.Errorf("failed to encrypt zeros bg: %w", err)
		}
		zerosBgVal := hammingDistance(dstBuf, zeroBaseDst)

		// --- Ones Background ---
		for j := range srcBuf {
			srcBuf[j] = 0xFF
		}
		srcBuf[i/8] &^= (1 << (i % 8))
		dstBuf, err = priv.Sign(nil, srcBuf, nil)
		if err != nil {
			return fmt.Errorf("failed to encrypt ones bg: %w", err)
		}
		onesBgVal := hammingDistance(dstBuf, oneBaseDst)

		// --- Random Background Avg ---
		totalRandDiff := 0
		for s := range rndSamples {
			copy(srcBuf, randSamplesSrc[s])
			srcBuf[i/8] ^= (1 << (i % 8))
			dstBuf, err = priv.Sign(nil, srcBuf, nil)
			if err != nil {
				return fmt.Errorf("failed to encrypt random background: %w", err)
			}
			totalRandDiff += hammingDistance(dstBuf, randSamplesDst[s])

			for o := range blockSizeBits {
				if ((randSamplesDst[s][o/8] ^ dstBuf[o/8]) & (1 << (o % 8))) != 0 {
					changeCounts[o]++
				}
			}
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

	if correlation {
		pub, err := priv.PublicKey()
		if err != nil {
			return fmt.Errorf("failed to derive public key: %w", err)
		}
		desc := pub.Describe()

		corrFile, err := os.Create(filepath.Join(folder, "monom-correlation.csv"))
		if err != nil {
			return fmt.Errorf("failed to create monom-correlation.csv: %w", err)
		}
		defer corrFile.Close()

		corrWriter := csv.NewWriter(corrFile)
		defer corrWriter.Flush()

		if err := corrWriter.Write([]string{"output_bit", "monomial_count", "avalanche_prob"}); err != nil {
			return fmt.Errorf("failed to write CSV header to monom-correlation.csv: %w", err)
		}

		for o := range blockSizeBits {
			prob := float64(changeCounts[o]) / float64(blockSizeBits*rndSamples)
			row := []string{
				strconv.Itoa(o),
				strconv.Itoa(desc.MonomialCounts[o]),
				fmt.Sprintf("%.6f", prob),
			}
			if err := corrWriter.Write(row); err != nil {
				return fmt.Errorf("failed to write row %d to monom-correlation.csv: %w", o, err)
			}
		}
	}

	return nil
}
