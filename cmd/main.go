package main

import (
	"fmt"
	"os"

	"github.com/aogirikarma/mini-stackr-cli/pkg/docker"
	"github.com/aogirikarma/mini-stackr-cli/pkg/tui"
)

func main() {
	client, err := docker.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to docker: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	if err := tui.Run(client); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}