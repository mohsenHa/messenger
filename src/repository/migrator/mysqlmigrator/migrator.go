package mysqlmigrator

import (
	"database/sql"
	"fmt"
	"github.com/mohsenHa/messenger/repository/mysql"
	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dialect    string
	dbConfig   mysql.Config
	migrations *migrate.FileMigrationSource
	db         *sql.DB
}

func New(dbConfig mysql.Config) Migrator {
	// OR: Read migrations from a folder:
	migrations := &migrate.FileMigrationSource{
		Dir: "./repository/mysql/migrations",
	}
	db := mysql.NewDB(dbConfig)

	return Migrator{dbConfig: dbConfig, dialect: "mysql", migrations: migrations, db: db}
}

func (m Migrator) Up() {

	n, err := migrate.Exec(m.db, m.dialect, m.migrations, migrate.Up)
	if err != nil {
		panic(fmt.Errorf("can't apply migrations: %w", err))
	}
	fmt.Printf("Applied %d migrations!\n", n)
}

func (m Migrator) Down() {
	n, err := migrate.Exec(m.db, m.dialect, m.migrations, migrate.Down)
	if err != nil {
		panic(fmt.Errorf("can't rollback migrations: %W", err))
	}
	fmt.Printf("Rollbacked %d migrations!\n", n)
}
