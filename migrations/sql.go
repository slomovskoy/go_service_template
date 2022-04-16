// Package migrations implements migrations interface
package migrations

import (
	"context"
	"github.com/jackc/pgx/v4"
	migrate "github.com/rubenv/sql-migrate"
	"go_service_template/internal/sql"
)

const (
	defaultSchemaName = "public"
	defaultSrc        = "./"
)

type SQL struct {
	connConf *pgx.ConnConfig
	src      string
	set      *migrate.MigrationSet
}

func New(ctx context.Context, connString string, src string) (*SQL, error) {
	dbConf, err := pgx.ParseConfig(connString)
	if err != nil {
		applog.FromContext(ctx).WithError(err).Error("parse postgres connection tring failed")
		return nil, err
	}
	schemaName, ok := dbConf.RuntimeParams["search_path"]
	if !ok {
		schemaName = defaultSchemaName
	}
	if src == "" {
		src = defaultSrc
	}
	s := &SQL{
		connConf: dbConf,
		set: &migrate.MigrationSet{
			SchemaName: schemaName,
		},
		src: src,
	}
	return s, nil
}

func (s *SQL) Up(ctx context.Context, count int) (int, error) {
	return s.run(ctx, count, migrate.Up)
}

func (s *SQL) Down(ctx context.Context, count int) (int, error) {
	return s.run(ctx, count, migrate.Up)
}

func (s *SQL) run(ctx context.Context, count int, dir migrate.MigrationDirection) (int, error) {
	logger := applog.FromContext(ctx)
	dbConn, err := sql.OpenDatabase(ctx, &sql.Config{ConnString: s.conf.ConnString()}, sql.Reg())
	if err != nil {
		logger.WithError(err).Fatal("db connect failed")
		return 0, err
	}
	defer func() {
		clErr := dbConn.Close()
		if clErr != nil {
			logger.WithError(clErr).Error("closing db connection is failed")
		}
	}()
	migrated, err := s.set.ExecMax(dbConn.DB, "postgres", &migrate.FileMigrationSource{Dir: s.src}, dir, count)
	if err != nil {
		logger.WithError(err).Error("migrations execution failed")
		return 0, err
	}
	return migrated, nil
}
