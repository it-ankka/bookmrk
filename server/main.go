package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /bookmarks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello world!")
		json.NewEncoder(w).Encode(struct{ bookmarks []string }{bookmarks: []string{}})
	})
}
