package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/asekhamhe/boolang/controllers"
)

func main() {

	r := httprouter.New()
	bc := controllers.NewBookController()

	r.GET("/", bc.HomePage)
	r.GET("/books", bc.GetBooks)
	r.GET("/books/:id", bc.GetBook)
	r.POST("/books", bc.AddBook)
	r.PUT("/books/:id", bc.UpdateBook)
	r.DELETE("/books/:id", bc.DeleteBook)

	log.Fatal(http.ListenAndServe(":8080", r))
}
