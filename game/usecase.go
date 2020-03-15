package game

import (
	"github.com/dbond762/naval-battle/player"
)

type Usecase interface {
	AddFirstPlayer(player.Usecase)
	AddSecondPlayer(player.Usecase)
	Complete() bool
	Start()
}
