package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	gameHttpDelivery "github.com/dbond762/naval-battle/game/delivery/http"

	"github.com/joho/godotenv"
)

func init() {
	env := os.Getenv("env")
	envFile := fmt.Sprintf("%s.env", env)
	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("Environment file not found")
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	gameHttpDelivery.NewGameHandler()

	addr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	log.Printf("Server run on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
