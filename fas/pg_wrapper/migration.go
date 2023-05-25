package pg_wrapper

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"
	"os"
	"path/filepath"
)

func Migration(host, database, user, password, port string, Filepath ...string) error {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, database)

	db, err := sql.Open("pgx/v5", connString)
	if err != nil {
		return err
	}

	pg, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+path(Filepath...), "postgres", pg)
	if err != nil {
		return err
	}

	if m.Up() != nil {
		return err
	}

	return db.Ping()
}

func path(Filepath ...string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var fp []string = []string{cwd}
	fp = append(fp, Filepath...)
	return filepath.Join(fp...)
}
