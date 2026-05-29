package main

import (
	"fmt"

	"github.com/staleread/aquila/asym"
)

func showConfig() {
	fmt.Printf("Block size: %d\n", asym.BlockSize*8)
	fmt.Printf("Compositions: %d\n", asym.Compositions)
	fmt.Printf("Fold size: %d\n", asym.FoldSize)
	fmt.Printf("Confusion Degree: %d\n", asym.Degree)
}

func getConfigID() string {
	return fmt.Sprintf("%dc%df%dd%d", asym.BlockSize*8, asym.Compositions, asym.FoldSize, asym.Degree)
}
