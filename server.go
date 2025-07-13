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
	db, err = sql.Open("mysql", "root:MyStrong@1234@tcp(127.0.0.1:3306)/testdb")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// ✅ Route handlers
	http.HandleFunc("/post-user", postUserHandler)      // POST
	http.HandleFunc("/users", getAllUsersHandler)       // GET
	http.HandleFunc("/update-user/", updateUserHandler) // PUT
	http.HandleFunc("/patch-user/", patchUserHandler)   // PATCH ✅
	http.HandleFunc("/delete-user/", deleteUserHandler) // DELETE

	fmt.Println("✅ Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// ✅ POST /post-user
func postUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

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
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT id, first_name, last_name, email, password, age FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Age)
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// ✅ PUT /update-user/{id}
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT method allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/update-user/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID in URL", http.StatusBadRequest)
		return
	}

	var user User
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &user)

	query := "UPDATE users SET first_name=?, last_name=?, email=?, password=?, age=? WHERE id=?"
	_, err = db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password, user.Age, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message": "User updated successfully"}`))
}

// ✅ PATCH /patch-user/{id}
func patchUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Only PATCH method allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/patch-user/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}

	var updates map[string]interface{}
	if err := json.Unmarshal(body, &updates); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(updates) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	var setParts []string
	var values []interface{}
	for key, val := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = ?", key))
		values = append(values, val)
	}
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(setParts, ", "))
	values = append(values, id)

	_, err = db.Exec(query, values...)
	if err != nil {
		http.Error(w, "Failed to patch user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message": "User patched successfully"}`))
}

// ✅ DELETE /delete-user/{id}
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE method allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/delete-user/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID in URL", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message": "User deleted successfully"}`))
}
