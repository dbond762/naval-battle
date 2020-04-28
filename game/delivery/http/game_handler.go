package http

import (
	"log"
	"net/http"

	"github.com/dbond762/naval-battle/domain"

	gameUsecase "github.com/dbond762/naval-battle/game/usecase"
	hubUsecase "github.com/dbond762/naval-battle/hub/usecase"
	playerUsecase "github.com/dbond762/naval-battle/player/usecase"

	"github.com/gorilla/websocket"
)

type GameHandler struct {
	gamesHub domain.HubUsecase
	upgrader websocket.Upgrader

	uncompletedGames chan domain.GameUsecase
}

func NewGameHandler() {
	handler := &GameHandler{
		gamesHub: hubUsecase.NewHub(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		uncompletedGames: make(chan domain.GameUsecase),
	}

	http.HandleFunc("/api/connect", handler.Connect)
}

func (g *GameHandler) Connect(w http.ResponseWriter, r *http.Request) {
	conn, err := g.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Connection error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	select {
	case uncompletedGame := <-g.uncompletedGames:
		newPlayer := playerUsecase.NewRealPlayerUsecase(conn, uncompletedGame)
		uncompletedGame.AddPlayer(newPlayer)
	default:
		newGame := gameUsecase.NewGameUsecase(g.gamesHub)
		newPlayer := playerUsecase.NewRealPlayerUsecase(conn, newGame)
		newGame.AddPlayer(newPlayer)

		g.uncompletedGames <- newGame
		g.gamesHub.Register(newGame)
	}
}
