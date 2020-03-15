package http

import (
	"log"
	"net/http"

	"github.com/dbond762/naval-battle/game"
	gameUsecase "github.com/dbond762/naval-battle/game/usecase"
	playerUsecase "github.com/dbond762/naval-battle/player/usecase"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type GameHandler struct {
	games []game.Usecase
}

func NewGameHandler() {
	handler := new(GameHandler)
	http.HandleFunc("/api/connect", handler.Connect)
}

func (g *GameHandler) Connect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	newPlayer := playerUsecase.NewRealPlayerUsecase(conn)

	if len(g.games) == 0 || g.games[len(g.games)-1].Complete() {
		newGame := gameUsecase.NewGameUsecase()
		newGame.AddFirstPlayer(newPlayer)
		g.games = append(g.games, newGame)
	} else {
		g.games[len(g.games)-1].AddSecondPlayer(newPlayer)
		go g.games[len(g.games)-1].Start()
	}
}
