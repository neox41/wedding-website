package db

import (
	"github.com/evilsocket/islazy/log"
	"mrwebsites/config"
	"os"
	"path/filepath"
)

func Setup() error {
	if _, errExist := os.Stat(filepath.FromSlash(config.DB)); os.IsNotExist(errExist) {
		if err := Create(); err != nil {
			log.Error("Error creating database")
			return err
		}
		return nil
	}
	return nil
}

func Create() error {
	var createTableStatements = []string{
		`CREATE TABLE attendees (
    id           INTEGER,
    name         TEXT,
    IDNucleo     INTEGER,
    surname      TEXT,
    vegan        INTEGER DEFAULT 0,
    vegetarian   INTEGER DEFAULT 0,
    transport    INTEGER,
    attendance   TEXT,
    category     TEXT,
    abroad       INTEGER,
    bambino      INTEGER,
    beve         INTEGER,
    tavolo       TEXT,
    requirements TEXT,
    PRIMARY KEY (
        id
    )
);`,
		`CREATE TABLE families (
    IDNucleo  INTEGER,
    Email     TEXT,
    link      TEXT,
    english   INTEGER DEFAULT 0,
    confirmed INTEGER DEFAULT 0,
    name      TEXT,
    status    TEXT,
    PRIMARY KEY (
        IDNucleo
    )
);`,
		`CREATE TABLE requests (
    id     INTEGER PRIMARY KEY AUTOINCREMENT,
    name   TEXT,
    email  TEXT,
    status TEXT
);`,
	}

	db = connect()
	defer close(db)
	for _, stmt := range createTableStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}
