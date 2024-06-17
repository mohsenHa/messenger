package mysql

import (
	"database/sql"
	"fmt"
	//use mysql driver
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Config struct {
	Username        string `koanf:"username"`
	Password        string `koanf:"password"`
	Port            int    `koanf:"port"`
	Host            string `koanf:"host"`
	DBName          string `koanf:"db_name"`
	MaxIdleConns    int    `koanf:"max_idle_conns"`
	MaxOpenConns    int    `koanf:"max_open_conns"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime"`
}

type DB struct {
	config Config
	db     *sql.DB
}

func (m *DB) Conn() *sql.DB {
	return m.db
}

func New(config Config) *DB {
	// parseTime=true changes the output type of DATE and DATETIME values to time.Time
	// instead of []byte / string
	// The date or datetime like 0000-00-00 00:00:00 is converted into zero value of time.Time
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
		config.Username, config.Password, config.Host, config.Port, config.DBName))
	if err != nil {
		panic(fmt.Errorf("can't open mysql db: %w", err))
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * time.Duration(config.ConnMaxLifetime))

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	return &DB{config: config, db: db}
}

func NewDB(cfg Config) *sql.DB {
	mdb := New(cfg)

	return mdb.db
}
