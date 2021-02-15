package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	b "github.com/asekhamhe/boolang/controllers"
)

func main() {

	r := httprouter.New()

	r.GET("/", b.HomePage)

	r.GET("/books", b.GetBooks)

	r.GET("/books/:id", b.GetBook)

	r.POST("/books", b.AddBook)
	r.PUT("/books/:id", b.UpdateBook)
	r.DELETE("/books/:id", b.DeleteBook)

	log.Fatal(http.ListenAndServe(":8080", r))
}
