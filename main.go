package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/lib/pq"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

// Book model
type Book struct {
	ID     int
	Title  string
	Author string
	Year   string
}

var books []Book
var db *sql.DB

func init() {
	gotenv.Load()

	pgURL, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	logFatal(err)

	db, err = sql.Open("postgres", pgURL)

	logFatal(err)

	err = db.Ping()
	logFatal(err)
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)

	}
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", homePage).Methods("GET")
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", addBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	// log.Println("Home page")
}

func getBooks(w http.ResponseWriter, r *http.Request) {

	var book Book
	books = []Book{}

	rows, err := db.Query("SELECT * FROM books")

	logFatal(err)

	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)

		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	params := mux.Vars(r)

	rows := db.QueryRow("select * from books where id=$1", params["id"])
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	logFatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(book)

}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	var bookID int

	// map json request to book variable
	json.NewDecoder(r.Body).Decode(&book)

	err := db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&bookID)

	logFatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)

}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	var bookID int
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	json.NewDecoder(r.Body).Decode(&book)

	err := db.QueryRow("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id", &book.Title, &book.Author, &book.Year, id).Scan(&bookID)

	logFatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(book)

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	_, err := db.Exec("delete from books where id=$1", id)

	logFatal(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}
