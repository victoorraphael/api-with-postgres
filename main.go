package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Person struct {
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

const (
	host     = "0.0.0.0"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "postgres"
)

func OpenConnection() *sql.DB {
	// string com credenciais de acesso
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// abre conexao com o banco
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	// verifica se a conexao esta de pe
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM person")
	if err != nil {
		log.Fatal(err)
	}

	var people []Person

	for rows.Next() {
		var person Person
		rows.Scan(&person.Name, &person.Nickname)
		people = append(people, person)
	}

	peopleBytes, _ := json.MarshalIndent(people, "", "\t")

	w.Header().Set("Content-type", "application/json")
	w.Write(peopleBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var p Person

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO person (name, nickname) VALUES ($1, $2)`
	_, errq := db.Exec(query, p.Name, p.Nickname)
	if errq != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(errq)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte("pong"))
}

func main() {
	http.HandleFunc("/", GETHandler)
	http.HandleFunc("/pint", testHandler)
	http.HandleFunc("/insert", POSTHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
