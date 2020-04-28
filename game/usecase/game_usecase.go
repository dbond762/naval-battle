package usecase

import (
	"github.com/dbond762/naval-battle/domain"
)

type gameUsecase struct {
	firstPlayer  domain.PlayerUsecase
	secondPlayer domain.PlayerUsecase
	gamesHub     domain.HubUsecase
}

func NewGameUsecase(gamesHub domain.HubUsecase) domain.GameUsecase {
	newGame := new(gameUsecase)
	newGame.gamesHub = gamesHub
	return newGame
}

func (g *gameUsecase) AddPlayer(player domain.PlayerUsecase) {
	if g.firstPlayer == nil {
		g.firstPlayer = player
	} else if g.secondPlayer == nil {
		g.secondPlayer = player
		go g.start()
	}
}

func (g *gameUsecase) Finish() {
	g.gamesHub.Unregister(g)
}

func (g *gameUsecase) start() {
	coords := make(chan domain.Coords)
	result := make(chan domain.Result)

	g.firstPlayer.DoAction(domain.ActionStart)
	g.secondPlayer.DoAction(domain.ActionStart)

	for {
		go g.firstPlayer.DoStep(coords, result)
		go g.secondPlayer.EnemyStep(coords, result)

		go g.secondPlayer.DoStep(coords, result)
		go g.firstPlayer.EnemyStep(coords, result)
	}
}
