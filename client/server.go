package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/funcionarios/", RequestsHandler)
	log.Println("Executando...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
