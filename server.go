package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// To test the setup, start the fictive server with the following command :
// go run server.go
// If you followed along with the setup, you should see the following output in your terminal.
// Starting server at port 8080
// the next step is to create a web server.
