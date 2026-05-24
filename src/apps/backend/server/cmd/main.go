package main

import (
	"log"

	"github.com/Nanhtu187/online-judge/src/apps/backend/server/config"
	servercmd "github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/cmd"
	commoncmd "github.com/Nanhtu187/online-judge/src/packages/cmd"
	"github.com/spf13/cobra"
)

func main() {
	command := &cobra.Command{}
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	command.AddCommand(servercmd.StartServerCommand())
	command.AddCommand(servercmd.StartResultConsumerCommand())
	command.AddCommand(commoncmd.MigrateCommand(cfg.Database.MigrationSource, cfg.Database.MigrationDSN()))

	if err := command.Execute(); err != nil {
		log.Fatalf("error executing command: %v", err)
	}
}
