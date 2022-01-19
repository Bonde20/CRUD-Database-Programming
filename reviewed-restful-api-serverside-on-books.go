
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type Book struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

func initDb() *sql.DB {

	Db, err := sql.Open("postgres", "host=localhost  port=5432 user=postgres dbname=book  password= postgres sslmode=disable")

	checkErr(err)

	return Db
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := initDb()
	rows, err := db.Query("SELECT * From books")
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

	var book Book

	db := initDb()
	reqBody, err := ioutil.ReadAll(r.Body)
	checkErr(err)
	err = json.Unmarshal(reqBody, &book)
	checkErr(err)
	err = db.QueryRow("INSERT INTO books(Title) VALUES($1) returning Id", book.Title).Scan(&book.Id)
	checkErr(err)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(book)
}

func GetSingleBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	Vars := mux.Vars(r)
	id := Vars["id"]
	var book Book

	db := initDb()

	err := db.QueryRow("SELECT * From books where id = $1", id).Scan(&book.Id, &book.Title)
	checkErr(err)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)

}
func UpdateBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	_, err := strconv.Atoi(vars["id"])

	checkErr(err)

	var book Book

	reqBody, err := ioutil.ReadAll(r.Body)
	checkErr(err)
	err = json.Unmarshal(reqBody, &book)
	checkErr(err)
	db := initDb()
	_, err = db.Exec("Update books set title = $2, where id =$1", book.Title, book.Id)
	checkErr(err)

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	Id := vars["id"]

	db := initDb()

	_, err := db.Exec("DELETE From books where id = $1", Id)

	checkErr(err)

	w.WriteHeader(http.StatusOK)
}

func main() {

	mx := mux.NewRouter().StrictSlash(true)

	mx.HandleFunc("/books", GetAllBooks).Methods("GET")
	mx.HandleFunc("/books/{id}", GetSingleBook).Methods("GET")
	mx.HandleFunc("/book", AddBook).Methods("POST")
	mx.HandleFunc("/books/{id}", UpdateBook).Methods("PUT")
	mx.HandleFunc("/books/{id}", deleteBook).Methods("Delete")

	s := &http.Server{
		Addr:         "127.0.0.1:8080",
		Handler:      mx,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal("Server Startup Failed", s.ListenAndServe())

}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
