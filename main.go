package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/asekhamhe/boolang/controllers"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/asekhamhe/boolang/docs"
)

// @title Boolang
// @version 1.0
// @description This is a CRUD application.
// @termsOfService http://swagger.io/terms/

// @contact.name Lawrence Onaulogho
// @contact.url https://github.com/asekhamhe/
// @contact.email lawrence[at][gmail][dot][com]

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8080
// @BasePath /

func main() {

	r := mux.NewRouter()
	bc := controllers.NewBookController()

	r.HandleFunc("/home", bc.HomePage).Methods("GET")
	r.HandleFunc("/books", bc.GetBooks).Methods("GET")
	r.HandleFunc("/books/{id}", bc.GetBook).Methods("GET")
	r.HandleFunc("/books", bc.AddBook).Methods("POST")
	r.HandleFunc("/books/{id}", bc.UpdateBook).Methods("PATCH")
	r.HandleFunc("/books/{id}", bc.DeleteBook).Methods("DELETE")

	r.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
