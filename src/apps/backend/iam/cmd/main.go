package main

import (
	"log"

	"github.com/Nanhtu187/online-judge/src/apps/backend/iam/config"
	"github.com/Nanhtu187/online-judge/src/apps/backend/iam/internal/cmd"
	commoncmd "github.com/Nanhtu187/online-judge/src/packages/cmd"
	"github.com/spf13/cobra"
)

func main() {
	command := &cobra.Command{}
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	command.AddCommand(cmd.StartServerCommand())
	command.AddCommand(commoncmd.MigrateCommand(cfg.Database.MigrationSource, cfg.Database.MigrationDSN("iam_schema_migrations")))

	if err := command.Execute(); err != nil {
		log.Fatalf("failed to execute command: %v", err)
	}
}
