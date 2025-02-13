package main

import (
	"database/sql"  //Interacts with the SQL database.
	"encoding/json" //Provides functionality to encode and decode JSON data.
	"fmt"
	"log"
	"net/http"

	//_ "github.com/joho/godotenv" //Loads environment variables from a .env file.
	_ "github.com/lib/pq" //A PostgreSQL driver for Go, enabling communication with PostgreSQL databases.
)

const (
	host     = "localhost"
	port     = 5432
	user     = "myuser"
	password = "mysecretpassword"
	dbname   = "mydatabase"
)

var DB *sql.DB

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			log.Println("Error scanning user", err)
			continue
		}
		users = append(users, user)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func initDB() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

}

func loadhomepage(w http.ResponseWriter, r *http.Request) {
	// Query the database to get users
	rows, err := DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, "Error fetching users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			http.Error(w, "Error scanning user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// HTML template to render the users
	html := "<html><head><title>Users List</title></head><body>"
	html += "<h1>Users List</h1>"
	html += "<table border='1'><tr><th>ID</th><th>Name</th><th>Email</th></tr>"

	// Loop through users and display them in a table
	for _, user := range users {
		html += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td></tr>", user.ID, user.Name, user.Email)
	}

	html += "</table></body></html>"

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func main() {
	initDB()
	// Set up routes and start your server here...
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/", loadhomepage)

	//start server
	fmt.Println("starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
