package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	. "github.com/gobeam/mongo-go-pagination"

	"github.com/gorilla/mux"
	"github.com/lorezi/boolang/inits"
	"github.com/lorezi/boolang/models"
)

var db *sql.DB
var m *mongo.Client

// BookController struct
type BookController struct {
}

func init() {
	db = inits.NewDB().Init()
	m = inits.NewDB().MongoConn()

}

// NewBookController instance
func NewBookController() *BookController {
	return &BookController{}
}

// HomePage is
// @Summary HomePage
// @Description Test connection
// @Produce plain
// @Success 200 {string} string "ok"
// @response default {string} string
// @Router /home [get]
func (bc BookController) HomePage(w http.ResponseWriter, r *http.Request) {
	// log.Println("Home page")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	fmt.Fprintln(w, "Api Welcome Page")

}

// GetBooks is
// @Summary GetBooks
// @Description fetch list of books
// @Produce json
// @Success 200 {object} models.Book "ok"
// @Router /books [get]
func (bc BookController) GetBooks(w http.ResponseWriter, r *http.Request) {

	// book := models.BookResult{}
	books := []models.BookResult{}
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	inits.LogFatal(err)
	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	inits.LogFatal(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.Database("boolang").Collection("books")

	filter := bson.D{{}}

	paginatedData, err := New(collection).Context(ctx).Limit(limit).Page(page).Filter(filter).Decode(&books).Find()
	inits.LogFatal(err)

	for _, v := range paginatedData.Data {
		var book *models.BookResult

		if err := bson.Unmarshal(v, &book); err == nil {

			fmt.Printf("books: %v\n", *book)
			books = append(books, *book)
		}

	}

	// for res.Pagination.Next(ctx) {

	// 	err := res.Decode(&book)
	// 	inits.LogFatal(err)
	// 	// fmt.Println(book.Author)
	// 	// fmt.Println()

	// 	books = append(books, book)
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(books)
}

// GetBook is
// @Summary GetBook
// @Description fetch a single book
// @Param id path string true "Book ID"
// @Produce json
// @Success 200 {object} models.Book "ok"
// @Router /books/{id} [get]
func (bc BookController) GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	book := models.BookResult{}

	paramID := mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(paramID["id"])
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
// @Summary CreateBook
// @Description create a new book
// @Param id body models.CreateBook true "book model"
// @Produce json
// @Accept json
// @Success 201 {object} models.CreateBook "ok"
// @Router /books [post]
func (bc BookController) AddBook(w http.ResponseWriter, r *http.Request) {
	b := models.Book{}

	// map json request to b variable
	json.NewDecoder(r.Body).Decode(&b)

	collection := m.Database("boolang").Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, b)
	if err != nil {
		r := models.Result{
			Status:  "fail",
			Message: "Server error 😰😰😰",
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(r)
		return
	}
	// id := res.InsertedID

	// converts primitive objectID type to string
	id := res.InsertedID.(primitive.ObjectID).Hex()

	nb := models.BookResult{
		ID:   id,
		Book: b,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nb)

}

// UpdateBook is
func (bc BookController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	book := models.Book{}
	paramID := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(paramID["id"])
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
func (bc BookController) DeleteBook(w http.ResponseWriter, r *http.Request) {

	paramID := mux.Vars(r)
	id, _ := strconv.Atoi(paramID["id"])

	_, err := db.Exec("delete from books where id=$1", id)

	inits.LogFatal(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}
