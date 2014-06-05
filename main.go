package main

import (
	"fmt"
	"net/http"
	"log"
	"database/sql"
	_ "github.com/lib/pq"
	"text/template"
)

var (
	db *sql.DB

	createTable = `CREATE TABLE IF NOT EXISTS users(
		name character varying(100) NOT NULL,
		email character varying(100) NOT NULL,
		description character varying(500) NOT NULL
    );`
)

const (
	DB_USER = "postgres"
	DB_PASSWORD = "ankit1234"
	DB_NAME = "go_posgres_db"
)

type User struct {
	Name string
	Email string
	Description string
}

func setupDB() *sql.DB {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME )
	db, err := sql.Open("postgres", dbInfo)
	PanicIf(err)
	return db
}

func PanicIf(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request ) {
	fmt.Println("In indexHandler")
	rows, err := db.Query("SELECT * FROM users")
	PanicIf(err)
	fmt.Println("Rows:",rows)
	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Name, &user.Email, &user.Description)
		PanicIf(err)
		users = append(users, user)
	}
	t := template.New("new.html")
	t, err = template.ParseFiles("templates/index.html")
	PanicIf(err)
	t.Execute(w, users)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In newHandler")
	t := template.New("new.index")
	t, _ = template.ParseFiles("templates/new.html")
	t.Execute(w, t)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In saveHanlder")
	res, err := db.Exec("INSERT INTO users (name, email, description) VALUES ($1, $2, $3)",
		r.FormValue("name"),
		r.FormValue("email"),
		r.FormValue("description"))
	PanicIf(err)
	fmt.Println(res)
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	db = setupDB()
	defer db.Close()

	ctble, err := db.Query(createTable)
	PanicIf(err)
	fmt.Println("Table created successfully", ctble)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/new/", newHandler)
	http.HandleFunc("/save/", saveHandler)
	fmt.Println("Listening server......")

	http.Handle("/public/css/", http.StripPrefix("/public/css/", http.FileServer(http.Dir("public/css"))))
	http.Handle("/public/images/", http.StripPrefix("/public/images/", http.FileServer(http.Dir("public/images"))))
    if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		log.Fatalf("Template Execution Error:", err)
	}
}
