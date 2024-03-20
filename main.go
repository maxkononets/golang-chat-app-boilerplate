package main

import (
	"fmt"
	"net/http"
)

func main() {
	// test is our app running
	fmt.Println("Hello from Go Chat App Backend")

	// Handles static files from public directory
	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.ListenAndServe(":3000", nil)
}
