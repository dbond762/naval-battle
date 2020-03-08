package main

import (
	"log"
	"net/http"

	gameHttpDelivery "github.com/dbond762/naval-battle/game/delivery/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fileServer)

	gameHttpDelivery.NewGameHandler()

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
