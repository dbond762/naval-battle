package usecase

import (
	"encoding/json"
	"github.com/dbond762/naval-battle/domain"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type realPlayerUsecase struct {
	conn   *websocket.Conn
	enemy  domain.PlayerUsecase
	myGame domain.GameUsecase

	readMsg  chan []byte
	writeMsg chan []byte
}

func NewRealPlayerUsecase(conn *websocket.Conn, myGame domain.GameUsecase) domain.PlayerUsecase {
	return &realPlayerUsecase{
		conn:     conn,
		myGame:   myGame,
		readMsg:  make(chan []byte),
		writeMsg: make(chan []byte),
	}
}

func (p *realPlayerUsecase) DoAction(action domain.Action) error {
	msg, err := json.Marshal(domain.ActionMsg{Action: action})
	if err != nil {
		return err
	}

	p.writeMsg <- msg

	return nil
}

func (p *realPlayerUsecase) DoStep(coords chan domain.Coords, result chan domain.Result) {
	p.DoAction(domain.ActionDoStep)

	var coordsMsg domain.CoordsMsg
	for {
		if err := json.Unmarshal(<-p.readMsg, &coordsMsg); err != nil {
			log.Print(err)
			continue
		}
		break
	}
	coords <- coordsMsg.Coords

	msg, err := json.Marshal(domain.ResultMsg{Result: <-result})
	if err != nil {
		//
	}
	p.writeMsg <- msg
}

func (p *realPlayerUsecase) EnemyStep(coords chan domain.Coords, result chan domain.Result) {
	msg, err := json.Marshal(domain.CoordsMsg{Coords: <-coords})
	if err != nil {
		//
	}
	p.writeMsg <- msg

	var resultMsg domain.ResultMsg
	for {
		if err := json.Unmarshal(<-p.readMsg, &resultMsg); err != nil {
			log.Print(err)
			continue
		}
		break
	}
	result <- resultMsg.Result
}

func (p *realPlayerUsecase) readConn() {
	defer func() {
		p.myGame.Finish()
		if err := p.conn.Close(); err != nil {
			log.Print(err)
		}
	}()

	p.conn.SetReadLimit(maxMessageSize)
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	p.conn.SetPongHandler(func(string) error {
		p.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		p.readMsg <- msg
	}
}

func (p *realPlayerUsecase) writeConn() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := p.conn.Close(); err != nil {
			log.Print(err)
		}
	}()

	for {
		select {
		case msg, ok := <-p.writeMsg:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := p.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
