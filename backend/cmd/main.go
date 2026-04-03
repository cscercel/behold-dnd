package main


import (
	"log"
	"net/http"
)


func main() {
	log.Println("Beyond DnD")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
