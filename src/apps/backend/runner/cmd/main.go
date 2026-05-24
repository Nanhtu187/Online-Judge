package main

import (
	"log"

	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/internal/cmd"
	"github.com/spf13/cobra"
)

func main() {
	command := &cobra.Command{}

	command.AddCommand(cmd.StartJudgerCommand())
	if err := command.Execute(); err != nil {
		log.Fatalf("failed to execute command: %v", err)
	}
}
