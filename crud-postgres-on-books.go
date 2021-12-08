
package main

import (
      "database/sql"
         "fmt"
      _ "github.com/lib/pq"
	  "github.com/Bonde20/crud-postgres-on-books"
)

type Book struct {
   Id string
   Title string
}

var Db *sql.DB
func init() {
var err error
Db, err = sql.Open("postgres", "host= localhost port=5432 user=postgres dbname=book password=postgres
 sslmode=disable")
  if err != nil {
	log.Fatal(err)

    }
}

func Books(books []Book, err error) {
rows, err := Db.Query("select id, title")
	if err != nil {
	return
}
for rows.Next() {
book := Book{}
err = rows.Scan(&book.Id, &book.Title)
	if err != nil {
	return
}
books = append(books, book)
	}
	rows.Close()
	return
}

func GetBook(id string) (book Book, err error) {
	book = Book{}
	err = Db.QueryRow("select id, title where id =
	$1", id).Scan(&book.Id, &book.Title)
		return
}

func (book *Book) Create() (err error) {
statement := "insert into books (title) values ($1)
			 returning id"
stmt, err := Db.Prepare(statement)
	if err != nil {
		return
			}
	defer stmt.Close()
	err = stmt.QueryRow(book.Title).Scan(&book.Id)
		return
}

func main() {
book := Book{Title: "RobinHood"}
fmt.Println(book)
book.Create()
fmt.Println(book)

readBook, _ := GetBook(book.Id)
fmt.Println(readBook)
readBook.Title = "RobinHood"

}
