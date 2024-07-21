package config

import (
	"context"
	"database/sql"
	"regexp"

	"github.com/alifakhimi/simple-service-go"
	"github.com/mattn/go-sqlite3"
)

type ServiceConfig struct {
	*simple.Config `mapstructure:",squash"`
	// custom meta
	// Meta *Meta `json:"meta,omitempty" mapstructure:"meta"`
	Meta map[string]any `json:"meta,omitempty" mapstructure:"meta"`
}

type Meta struct {
	MetaValue `mapstructure:",squash"`
	Mock      MetaValue `json:"mock,omitempty" mapstructure:"mock"`
}

type MetaValue struct {
}

var (
	conf = ServiceConfig{}
)

func Config() *ServiceConfig {
	return &conf
}

// addRegexpFunction adds the REGEXP function to SQLite
func AddRegexpFunction(db *sql.DB) {
	conn, err := db.Conn(context.Background())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	_, err = conn.ExecContext(context.Background(), `PRAGMA case_sensitive_like = true`)
	if err != nil {
		panic(err)
	}

	err = conn.Raw(func(driverConn interface{}) error {
		sqliteConn := driverConn.(*sqlite3.SQLiteConn)
		return sqliteConn.RegisterFunc("regexp", func(re, s string) (bool, error) {
			return regexp.MatchString(re, s)
		}, true)
	})
	if err != nil {
		panic(err)
	}
}
