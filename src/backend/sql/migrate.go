package sql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

func MigrateSQL() {
	DATABASE_URL := os.Getenv("DATABASE_URL")

	migrations := &migrate.FileMigrationSource{
		Dir: "sql/migrations",
	}

	db, err := sql.Open("postgres", DATABASE_URL)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	fmt.Printf("Applied %d migrations!\n", n)
}
