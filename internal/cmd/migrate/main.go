package main

import (
	"go_service_template/internal/app"
	"go_service_template/internal/version"
	"go_service_template/migrations"
	"log"
	"os"
)

func main() {
	log.Printf("Applying %v@%v migrations", version.AppName, version.AppVersion)
	conf, err := app.ReadConfigFromFile("")
	if err != nil {
		log.Fatal(err)
	}

	logger := applog.New()
	ctx := applog.WithLogger(logger)

	logger.Infof("Migrations destination is: %v", conf.Database.ConnString)
	sql, err := migrations.New(ctx, conf.Database.ConnString, "./migrations/sql")
	if err != nil {
		logger.WithError(err).Fatal("failed instantiate migrations")

	}

	var count int
	if len(os.Args) > 1 && os.Args[1] == "down" { // TODO: rewrite on cobra
		count, err = sql.Down(ctx, 1)
	} else {
		count, err = sql.Up(ctx, 0)
	}
	if err != nil {
		logger.WithError(err).Fatal("failed make migrations")
	}
	log.Printf("Applied %v migrations", count)
}
