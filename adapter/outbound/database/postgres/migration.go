package postgres

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/robertobff/food-service/adapter/outbound/database/postgres/migration"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var MigrationModule = fx.Module(
	"migration",
	fx.Provide(NewMigration),
	fx.Invoke(migrate),
)

type Migration struct {
	m *gormigrate.Gormigrate
}

func NewMigration(pg *Postgres) *Migration {
	m := Migration{}
	m.load(pg)
	return &m
}

func (m *Migration) load(pg *Postgres) {
	m.m = gormigrate.New(
		pg.Db,
		&gormigrate.Options{
			TableName:                 "migrations",
			IDColumnName:              "id",
			IDColumnSize:              255,
			UseTransaction:            false,
			ValidateUnknownMigrations: true,
		}, []*gormigrate.Migration{
			migration.M_202507301810,
		},
	)
}

func (m *Migration) Migrate(l *zap.SugaredLogger) error {
	if err := m.m.Migrate(); err != nil {
		l.Errorf("could not migrate: %v", err)
		return err
	}
	return nil
}

func (m *Migration) MigrateTo(migrationID string) error {
	err := m.m.MigrateTo(migrationID)
	return err
}

func (m *Migration) RollbackLast() error {
	err := m.m.RollbackLast()
	return err
}

func (m *Migration) RollbackTo(migrationID string) error {
	err := m.m.RollbackTo(migrationID)
	return err
}

func migrate(m *Migration, c *Config, l *zap.SugaredLogger) error {
	if *c.AutoMigrate {
		if err := m.Migrate(l); err != nil {
			l.Errorf("could not migrate: %v", err)
			return err
		}
	}

	l.Info("migration complete")
	return nil
}
