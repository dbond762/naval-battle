package usecase

import (
	"github.com/dbond762/naval-battle/player"

	"github.com/gorilla/websocket"
)

type realPlayerUsecase struct {
	conn *websocket.Conn
}

func NewRealPlayerUsecase(conn *websocket.Conn) player.Usecase {
	return &realPlayerUsecase{conn}
}

func (p *realPlayerUsecase) Step(coords chan<- []byte, result <-chan []byte) error {
	_, msg, err := p.conn.ReadMessage()
	if err != nil {
		return err
	}
	coords <- msg

	if err := p.conn.WriteMessage(websocket.TextMessage, <-result); err != nil {
		return err
	}

	return nil
}

func (p *realPlayerUsecase) EnemyStep(coords <-chan []byte, result chan<- []byte) error {
	if err := p.conn.WriteMessage(websocket.TextMessage, <-coords); err != nil {
		return err
	}

	_, msg, err := p.conn.ReadMessage()
	if err != nil {
		return err
	}
	result <- msg

	return nil
}
