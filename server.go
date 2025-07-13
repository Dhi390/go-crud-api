package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// User model struct
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Age       int    `json:"age"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:MyStrong@1234@tcp(127.0.0.1:3306)/testdb") // ✅ DB credentials here
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// ✅ Route registrations
	http.HandleFunc("/post-user", postUserHandler)
	http.HandleFunc("/users", getAllUsersHandler)
	http.HandleFunc("/update-user", updateUserHandler)
	http.HandleFunc("/delete-user/", deleteUserHandler)

	fmt.Println("✅ Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// ✅ POST /post-user
func postUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &user)

	query := "INSERT INTO users (first_name, last_name, email, password, age) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password, user.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message": "User created successfully"}`))
}

// ✅ GET /users
func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, first_name, last_name, email, password, age FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Age)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// ✅ POST /update-user
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &user)

	query := "UPDATE users SET first_name=?, last_name=?, email=?, password=?, age=? WHERE id=?"
	_, err := db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password, user.Age, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message": "User updated successfully"}`))
}

// ✅ DELETE /delete-user/{id}
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/delete-user/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message": "User deleted successfully"}`))
}
