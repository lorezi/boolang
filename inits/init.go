package inits

import (
	"database/sql"
	"log"
	"os"

	// imported strictly for setup
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

var db *sql.DB

// Init is ...
func Init() *sql.DB {
	gotenv.Load()

	db, err := sql.Open("postgres", os.Getenv("ELEPHANTSQL_URL"))

	LogFatal(err)

	err = db.Ping()
	LogFatal(err)

	return db
}

// LogFatal is ...
func LogFatal(err error) {
	if err != nil {
		log.Fatal(err)

	}

}
