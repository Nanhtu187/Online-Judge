package main

import (
	"fmt"
	"github.com/Nanhtu187/Online-Judge/app/iam/config"
	"github.com/Nanhtu187/Online-Judge/app/iam/pkg/migration"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	cmd := migration.MigrateCommand(conf.Database.DSN())
	err = cmd.Execute()
	if err != nil {
		fmt.Println("ERROR:", err)
	}
}
