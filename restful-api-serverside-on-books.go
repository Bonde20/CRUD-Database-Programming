
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type Book struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

var Db *sql.DB

func initDb() *sql.DB {

	var err error

	Db, err = sql.Open("postgres", "host=localhost  port=5432 user=postgres dbname=book  password= postgres sslmode=disable")

	checkErr(err)

	return Db
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	Db := initDb()
	rows, err := Db.Query("SELECT * From books")
	checkErr(err)

	var books []Book

	for rows.Next() {
		book := Book{}
		err = rows.Scan(&book.Id, &book.Title)

		checkErr(err)

		books = append(books, book)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)

}

func AddBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vals := r.URL.Query()
	Title := vals["title"]
	Db := initDb()
	var lastInsertedID int
	err := Db.QueryRow("INSERT INTO books(Title) VALUES($1) returning Id;", Title).Scan(&lastInsertedID)
	checkErr(err)
	w.WriteHeader(http.StatusCreated)

}

func main() {

	mx := mux.NewRouter().StrictSlash(true)

	mx.HandleFunc("/books", GetAllBooks).Methods("GET")

	mx.HandleFunc("/book", AddBook).Methods("POST")

	s := &http.Server{
		Addr:         "127.0.0.1:8080",
		Handler:      mx,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server Startup Failed")
	}

}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
