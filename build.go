package main

import (
	"fmt"
	"os"
	"os/exec"
)

func runBuild(configID string) error {
	var block, comp, fold, deg int
	n, err := fmt.Sscanf(configID, "%dc%df%dd%d", &block, &comp, &fold, &deg)
	if err != nil || n != 4 {
		return fmt.Errorf("invalid config ID format (must be <block>c<comp>f<fold>d<deg>): %w", err)
	}

	tags := fmt.Sprintf("block%d comp%d fold%d deg%d", block, comp, fold, deg)
	fmt.Printf("Building aquila-cli with tags: %s\n", tags)

	cmd := exec.Command("go", "build", "-tags", tags)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go build: %w", err)
	}

	fmt.Println("Build successful!")
	return nil
}
