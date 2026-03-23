package main

import (
	"cars/pkg/handlers"
	"fmt"
	"net/http"
)

const port = ":8080"

func main() {
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/car", handlers.CarHandler)
	http.HandleFunc("/filter", handlers.FilterHandler)
	http.HandleFunc("/compare", handlers.ComparisonHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Server is running on http://localhost:8080")

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
