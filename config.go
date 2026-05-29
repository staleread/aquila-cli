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
