package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/lorezi/boolang/controllers"
	"github.com/lorezi/boolang/middleware"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/lorezi/boolang/docs"
)

// @title Boolang
// @version 1.0
// @description This is a CRUD application.
// @termsOfService http://swagger.io/terms/

// @contact.name Lawrence Onaulogho
// @contact.url https://github.com/lorezi/
// @contact.email lawrence[at][gmail][dot][com]

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8080
// @BasePath /

func main() {

	r := mux.NewRouter()
	bc := controllers.NewBookController()
	uc := controllers.NewUserController()

	subr := r.PathPrefix("/api/v1").Subrouter()
	subr.Use(middleware.Authentication)

	subr.HandleFunc("/home", bc.HomePage).Methods("GET")
	subr.HandleFunc("/books", bc.GetBooks).Methods("GET").Queries("limit", "{limit:[0-9]+}", "page", "{page:[0-9]+}")
	subr.HandleFunc("/books/{id}", bc.GetBook).Methods("GET")
	subr.HandleFunc("/books", bc.AddBook).Methods("POST")
	subr.HandleFunc("/books/{id}", bc.UpdateBook).Methods("PATCH")
	subr.HandleFunc("/books/{id}", bc.DeleteBook).Methods("DELETE")

	subr.HandleFunc("/users", uc.GetUsers).Methods("GET").Queries("limit", "{limit:[0-9]+}", "page", "{page:[0-9]+}")

	r.HandleFunc("/users/login", uc.Login).Methods("POST")
	r.HandleFunc("/users/signup", uc.CreateUser).Methods("POST")
	r.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler)

	handler := cors.Default().Handler(r)
	srv := &http.Server{
		Handler:      handler,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
