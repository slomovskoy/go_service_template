package sql

import (
	"context"
	"database/sql"
	pgDriver "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/luna-duclos/instrumentedsql/opentracing"
	"io"
)

const (
	wrappedDriverName = "cloudsqlpostgres"
)

// Reg register new wrappedDriver for trace
func Reg() string {
	drv := wrappedDriverName
	if !isDriverRegistered(drv) {
		sql.Register(drv,
			instrumentedsql.WrapDriver(&pgDriver.Driver{},
				instrumentedsql.WithTracer(opentracing.NewTracer(true)),
			),
		)
	}
	return wrappedDriverName
}

// OpenDatabase open postgresql database from DbConnectString
func OpenDatabase(ctx context.Context, conf *Config, driverName string) (*sqlx.DB, error) {
	logger := applog.FromContext(ctx)

	conn, err := sqlx.ConnectContext(ctx, driverName, conf.ConnString)
	if err != nil {
		logger.WithError(err).Error("Connection to database failed")
		return nil, err
	}
	conn.DB.SetMaxOpenConns(conf.MaxOpenConns)
	conn.DB.SetConnMaxLifetime(conf.ConnMaxLifetime)
	return conn, nil
}

// CloseDatabase close db database
func CloseDatabase(ctx context.Context, db io.Closer) error {
	logger := applog.FromContext(ctx)
	err := db.Close()
	if err != nil {
		logger.WithError(err).Error("Close connection to database failed")
	}
	return err
}

func isDriverRegistered(drv string) bool {
	driversList := sql.Drivers()
	for _, driverName := range driversList {
		if driverName == drv {
			return true
		}
	}
	return false
}
