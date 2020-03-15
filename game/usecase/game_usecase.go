package usecase

import (
	"github.com/dbond762/naval-battle/game"
	"github.com/dbond762/naval-battle/player"
)

type gameUsecase struct {
	FirstPlayer  player.Usecase
	SecondPlayer player.Usecase
	currentStep  int
}

func NewGameUsecase() game.Usecase {
	return new(gameUsecase)
}

func (g *gameUsecase) AddFirstPlayer(newPlayer player.Usecase) {
	g.FirstPlayer = newPlayer
}

func (g *gameUsecase) AddSecondPlayer(newPlayer player.Usecase) {
	g.SecondPlayer = newPlayer
}

func (g *gameUsecase) Complete() bool {
	return g.FirstPlayer != nil && g.SecondPlayer != nil
}

func (g *gameUsecase) Start() {
	g.currentStep = 1
	for {
		coords := make(chan []byte)
		result := make(chan []byte)
		switch g.currentStep {
		case 1:
			go g.FirstPlayer.Step(coords, result)
			go g.SecondPlayer.EnemyStep(coords, result)
			g.currentStep = 2
		case 2:
			go g.SecondPlayer.Step(coords, result)
			go g.FirstPlayer.EnemyStep(coords, result)
			g.currentStep = 1
		}
	}
}
