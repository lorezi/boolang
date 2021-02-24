package book

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	h "github.com/asekhamhe/boolang/inits"
)

// Book model
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   string `json:"year"`
}

var books []Book

var db *sql.DB

func init() {
	db = h.Init()
}

// HomePage is
func HomePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// log.Println("Home page")
	fmt.Fprintln(w, "Api Page Refactoring")

}

// GetBooks is
func GetBooks(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var book Book
	books = []Book{}

	rows, err := db.Query("SELECT * FROM books")

	h.LogFatal(err)

	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		h.LogFatal(err)

		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(books)
}

// GetBook is
func GetBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var book Book

	rows := db.QueryRow("Select * From books Where id=$1", ps.ByName("id"))
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	h.LogFatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(book)

}

// AddBook is
func AddBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var book Book
	var bookID int

	// map json request to book variable
	json.NewDecoder(r.Body).Decode(&book)

	err := db.QueryRow("Insert into books (title, author, year) values($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&bookID)

	h.LogFatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)

}

// UpdateBook is
func UpdateBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var book Book
	var bookID int

	id, _ := strconv.Atoi(ps.ByName("id"))
	json.NewDecoder(r.Body).Decode(&book)

	err := db.QueryRow("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id", &book.Title, &book.Author, &book.Year, id).Scan(&bookID)

	h.LogFatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(book)

}

// DeleteBook is
func DeleteBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, _ := strconv.Atoi(ps.ByName("id"))

	_, err := db.Exec("delete from books where id=$1", id)

	h.LogFatal(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}
