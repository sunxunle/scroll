package app

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"

	"bridge-history-api/config"
	"bridge-history-api/db"
	"bridge-history-api/db/migrate"
	"bridge-history-api/utils"
)

func getConfig(ctx *cli.Context) (*config.Config, error) {
	file := ctx.String(utils.ConfigFileFlag.Name)
	dbCfg, err := config.NewConfig(file)
	if err != nil {
		return nil, err
	}
	return dbCfg, nil
}

func initDB(dbCfg *config.Config) (*sqlx.DB, error) {
	factory, err := db.NewOrmFactory(dbCfg)
	if err != nil {
		return nil, err
	}
	log.Debug("Got db config from env", "driver name", dbCfg.DB.DriverName, "dsn", dbCfg.DB.DSN)

	return factory.GetDB(), nil
}

// resetDB clean or reset database.
func resetDB(ctx *cli.Context) error {
	cfg, err := getConfig(ctx)
	if err != nil {
		return err
	}
	db, err := initDB(cfg)
	if err != nil {
		return err
	}
	err = migrate.ResetDB(db.DB)
	if err != nil {
		return err
	}
	log.Info("successful to reset")
	return nil
}

// checkDBStatus check db status
func checkDBStatus(ctx *cli.Context) error {
	cfg, err := getConfig(ctx)
	if err != nil {
		return err
	}
	db, err := initDB(cfg)
	if err != nil {
		return err
	}

	return migrate.Status(db.DB)
}

// dbVersion return the latest version
func dbVersion(ctx *cli.Context) error {
	cfg, err := getConfig(ctx)
	if err != nil {
		return err
	}
	db, err := initDB(cfg)
	if err != nil {
		return err
	}

	version, err := migrate.Current(db.DB)
	log.Info("show database version", "db version", version)

	return err
}

// migrateDB migrate db
func migrateDB(ctx *cli.Context) error {
	cfg, err := getConfig(ctx)
	if err != nil {
		return err
	}
	db, err := initDB(cfg)
	if err != nil {
		return err
	}

	return migrate.Migrate(db.DB)
}

// rollbackDB rollback db by version
func rollbackDB(ctx *cli.Context) error {
	cfg, err := getConfig(ctx)
	if err != nil {
		return err
	}
	db, err := initDB(cfg)
	if err != nil {
		return err
	}
	version := ctx.Int64("version")
	return migrate.Rollback(db.DB, &version)
}