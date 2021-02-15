package inits

import (
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"
	"github.com/subosito/gotenv"
)

var db *sql.DB

// Init is ...
func Init() *sql.DB {
	gotenv.Load()

	pgURL, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	LogFatal(err)

	db, err = sql.Open("postgres", pgURL)

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
