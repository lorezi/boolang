package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/asekhamhe/boolang/inits"
	"github.com/asekhamhe/boolang/models"
)

var db *sql.DB
var m *mongo.Client

// BookController struct
type BookController struct{}

func init() {
	db = inits.NewDB().Init()
	m = inits.NewDB().MongoConn()

}

// NewBookController instance
func NewBookController() *BookController {
	return &BookController{}
}

// HomePage is
func (bc BookController) HomePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// log.Println("Home page")
	fmt.Fprintln(w, "Api Welcome Page")

}

// GetBooks is
func (bc BookController) GetBooks(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	book := models.Book{}
	books := []models.Book{}

	// rows, err := db.Query("SELECT * FROM books")
	collection := m.Database("boolang").Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{}}
	res, err := collection.Find(ctx, filter)
	inits.LogFatal(err)

	for res.Next(ctx) {

		err := res.Decode(&book)
		inits.LogFatal(err)

		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(books)
}

// GetBook is
func (bc BookController) GetBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	book := models.Book{}

	rows := db.QueryRow("Select * From books Where id=$1", ps.ByName("id"))
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	inits.LogFatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(book)

}

// AddBook is
func (bc BookController) AddBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	book := models.Book{}

	// map json request to book variable
	json.NewDecoder(r.Body).Decode(&book)

	collection := m.Database("boolang").Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx,
		bson.D{
			bson.E{
				Key:   "title",
				Value: book.Title,
			},
			bson.E{
				Key:   "author",
				Value: book.Author,
			},
			bson.E{
				Key:   "year",
				Value: book.Year,
			},
		})
	id := res.InsertedID
	inits.LogFatal(err)
	// converts primitive objectID type to string
	book.ID = id.(primitive.ObjectID).Hex()
	// fmt.Printf("Created with the id: %T", book.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)

}

// UpdateBook is
func (bc BookController) UpdateBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	book := models.Book{}
	bookID := 0

	id, _ := strconv.Atoi(ps.ByName("id"))
	json.NewDecoder(r.Body).Decode(&book)

	err := db.QueryRow("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id", &book.Title, &book.Author, &book.Year, id).Scan(&bookID)

	inits.LogFatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(book)
	inits.LogFatal(err)

}

// DeleteBook is
func (bc BookController) DeleteBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, _ := strconv.Atoi(ps.ByName("id"))

	_, err := db.Exec("delete from books where id=$1", id)

	inits.LogFatal(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}
