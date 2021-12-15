
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Book struct {
	Id    int
	Title string
}

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "host=localhost  port=5432 user=postgres dbname=book password= postgres"+
		"sslmode=disable")
	if err != nil {
		fmt.Println(err)

	}
	defer Db.Close()

	err = Db.Ping()
	if err != nil {
		fmt.Println(err)

	}

}

func GetAllBooks() (books []Book, err error) {
	rows, err := Db.Query("SELECT* from books")
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

func GetBook(id int) (book Book, err error) {
	book = Book{}
	err = Db.QueryRow("select id, title where id = $1", id).Scan(&book.Id, &book.Title)

	return
}

func (book *Book) AddBook() (err error) {
	statement := "insert into books (title) values ($1)returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {

		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(book.Title).Scan(&book.Id)
	return
}
func (book *Book) UpdateBook() (err error) {
	_, err = Db.Exec("update books set title = $2 where id = $1", book.Id, book.Title)
	return
}
func (book *Book) DeleteBook() (err error) {
	_, err = Db.Exec("delete from books where id = $1", book.Id)
	return
}

func main() {

	book := Book{Title: "Cinderella"}
	fmt.Println(book)
	book.AddBook()

	readBook, _ := GetBook(book.Id)
	fmt.Println("Book Printed:", readBook)

	readBook.Title = "Romance"
	readBook.UpdateBook()

	books, _ := GetAllBooks()
	fmt.Println("All Books Printed:", books)

	readBook.DeleteBook()
}
