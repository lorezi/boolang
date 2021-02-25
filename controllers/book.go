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
	w.Header().Set("Content-Type", "application/json")

	book := models.Book{}

	id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
	if err != nil {
		r := models.Result{
			Status:  "fail",
			Message: "Page not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(r)
		return
	}

	collection := m.Database("boolang").Collection("books")

	filter := bson.D{bson.E{
		Key:   "_id",
		Value: id,
	}}
	err = collection.FindOne(context.Background(), filter).Decode(&book)
	if err != nil {
		r := models.Result{
			Status:  "fail",
			Message: "Page not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(r)
		return
	}

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

	res, err := collection.InsertOne(ctx, book)
	if err != nil {
		r := models.Result{
			Status:  "fail",
			Message: "Server error 😰😰😰",
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(r)
		return
	}
	id := res.InsertedID

	// converts primitive objectID type to string
	book.ID = id.(primitive.ObjectID).Hex()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)

}

// UpdateBook is
func (bc BookController) UpdateBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	book := models.Book{}

	id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
	inits.LogFatal(err)

	filter := bson.D{
		bson.E{
			Key:   "_id",
			Value: id,
		},
	}

	json.NewDecoder(r.Body).Decode(&book)

	collection := m.Database("boolang").Collection("books")

	err = collection.FindOneAndUpdate(context.TODO(), filter, book).Decode(&book)

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
