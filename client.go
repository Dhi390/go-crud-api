package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Age       int    `json:"age"`
}

func main() {
	var action string

	fmt.Print("🔧 Enter action (create / read / update / patch / delete): ")
	fmt.Scanln(&action)

	switch strings.ToLower(action) {

	case "create":
		var user User
		fmt.Print("First Name: ")
		fmt.Scanln(&user.FirstName)
		fmt.Print("Last Name: ")
		fmt.Scanln(&user.LastName)
		fmt.Print("Email: ")
		fmt.Scanln(&user.Email)
		fmt.Print("Password: ")
		fmt.Scanln(&user.Password)
		fmt.Print("Age: ")
		fmt.Scanln(&user.Age)

		jsonData, _ := json.Marshal(user)
		resp, err := http.Post("http://localhost:8080/post-user", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("📦 CREATE Response:")
		fmt.Println(string(body))

	case "read":
		resp, err := http.Get("http://localhost:8080/users")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("📦 READ Response:")
		fmt.Println(string(body))

	case "update":
		var user User
		fmt.Print("User ID to update: ")
		fmt.Scanln(&user.ID)
		fmt.Print("New First Name: ")
		fmt.Scanln(&user.FirstName)
		fmt.Print("New Last Name: ")
		fmt.Scanln(&user.LastName)
		fmt.Print("New Email: ")
		fmt.Scanln(&user.Email)
		fmt.Print("New Password: ")
		fmt.Scanln(&user.Password)
		fmt.Print("New Age: ")
		fmt.Scanln(&user.Age)

		jsonData, _ := json.Marshal(user)
		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodPut, "http://localhost:8080/update-user/"+strconv.Itoa(user.ID), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("📦 UPDATE Response:")
		fmt.Println(string(body))

	case "patch":
		var id, field, value string
		fmt.Print("User ID to PATCH: ")
		fmt.Scanln(&id)
		fmt.Print("Field to update (e.g., email, first_name): ")
		fmt.Scanln(&field)
		fmt.Print("New value: ")
		fmt.Scanln(&value)

		patchData := fmt.Sprintf(`{"%s":"%s"}`, field, value)
		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodPatch, "http://localhost:8080/patch-user/"+id, strings.NewReader(patchData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("📦 PATCH Response:")
		fmt.Println(string(body))

	case "delete":
		var id string
		fmt.Print("User ID to delete: ")
		fmt.Scanln(&id)

		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodDelete, "http://localhost:8080/delete-user/"+id, nil)
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("📦 DELETE Response:")
		fmt.Println(string(body))

	default:
		fmt.Println("❌ Invalid action! Use: create | read | update | patch | delete")
		os.Exit(1)
	}
}
