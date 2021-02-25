package inits

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	// imported strictly for setup
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *sql.DB

// DB struct
type DB struct{}

// NewDB function
func NewDB() *DB {
	return &DB{}
}

// Init is ...
func (conn *DB) Init() *sql.DB {
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

// Error404 function
func Error404(err error, w http.ResponseWriter) {
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// MongoConn is connection setting
func (conn *DB) MongoConn() *mongo.Client {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URL")))

	err = client.Ping(ctx, readpref.Primary())
	LogFatal(err)

	return client
}

/*
COME ANOTHER FOR THIS BATTLE 😀😀😀
*/
// Collection for books
// func (conn *DB) Collection() (*mongo.Collection, context.Context) {
// 	m := conn.MongoConn()
// 	collection := m.Database("boolang").Collection("books")
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	return collection, ctx
// }
