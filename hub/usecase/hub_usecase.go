package usecase

import (
	"sync"

	"github.com/dbond762/naval-battle/domain"
)

type hubUsecase struct {
	mx    sync.Mutex
	games map[domain.GameUsecase]bool
}

func NewHub() domain.HubUsecase {
	return &hubUsecase{
		games: make(map[domain.GameUsecase]bool),
	}
}

func (h *hubUsecase) Register(game domain.GameUsecase) {
	h.mx.Lock()
	defer h.mx.Unlock()

	h.games[game] = true
}

func (h *hubUsecase) Unregister(game domain.GameUsecase) {
	h.mx.Lock()
	defer h.mx.Unlock()

	delete(h.games, game)
}
