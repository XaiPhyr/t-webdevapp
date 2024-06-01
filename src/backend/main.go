package main

import (
	"fmt"
	"log"

	sql "app_backend/sql"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("MAIN INITIALIZE...")

	log.SetFlags(log.Llongfile | log.LstdFlags)

	sql.MigrateSQL()
}
