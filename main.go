package main

import (
	"fmt"
	"net/http"

	"github.com/Ev357/bakalari-calendar-sync/api"
)

func main() {
	http.HandleFunc("/api/sync", handler.Handler)

	fmt.Println("Server starting on http://localhost:3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
