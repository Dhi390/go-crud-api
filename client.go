// client.go
package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	const url = "http://localhost:8080/update-user"

	jsonPayload := `{
		"id": 1,
		"firstName": "Dhiru",
		"lastName": "Yadav",
		"email": "updated@gmail.com",
		"password": "newpass",
		"age": 25
	}`

	requestBody := strings.NewReader(jsonPayload)

	resp, err := http.Post(url, "application/json", requestBody)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("📦 Response from server:")
	fmt.Println(string(content))
}
