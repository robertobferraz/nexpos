package postgres

import (
	"context"
	"log"
	"net"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Module = fx.Module(
	"postgres",
	ConfigModule,
	MigrationModule,
	fx.Provide(NewPostgres),
	fx.Invoke(HookPostgres),
	fx.Invoke(enableUUIDExtension),
)

type Postgres struct {
	Db *gorm.DB
}

func NewPostgres(c *Config, l *zap.SugaredLogger) (*Postgres, error) {
	db := &Postgres{}

	err := db.connect(c, l)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (p *Postgres) connect(c *Config, l *zap.SugaredLogger) error {
	dsn := c.GetDsn(*c.ZeroThrust)

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		l.Error("pgx.ParseConfig: %w", zap.Error(err))
		return nil
	}

	if *c.ZeroThrust {
		d, err := cloudsqlconn.NewDialer(
			context.Background(),
			cloudsqlconn.WithIAMAuthN(),
		)
		if err != nil {
			return nil
		}

		config.DialFunc = func(ctx context.Context, _ string, _ string) (net.Conn, error) {
			return d.Dial(ctx, *c.Host, cloudsqlconn.WithPrivateIP())
		}
	}

	dbURI := stdlib.RegisterConnConfig(config)

	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DriverName: "pgx",
			DSN:        dbURI,
		}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		},
	)
	if err != nil {
		l.Errorf("error connecting to database: %v", err)
		return err
	}

	p.Db = db

	return nil
}

func HookPostgres(lc fx.Lifecycle, pg *Postgres, l *zap.SugaredLogger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			dbDriver, err := pg.Db.DB()
			if err != nil {
				l.Fatal("failed to get DB driver", zap.Error(err))
				return err
			}

			err = dbDriver.Ping()
			if err != nil {
				l.Fatal("failed to ping database", zap.Error(err))
				return err
			}

			l.Info("database ok!")
			return nil
		},
		OnStop: func(context.Context) error {
			dbDriver, err := pg.Db.DB()
			if err != nil {
				l.Fatal("failed to get DB driver", zap.Error(err))

			}

			err = dbDriver.Close()
			if err != nil {
				l.Fatal("failed to close database connection", zap.Error(err))

			}
			return nil
		},
	})
}

func enableUUIDExtension(pg *Postgres) {
	_, err := pg.Db.Raw(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Rows()

	if err != nil {
		log.Panicln(err)
	}
}
