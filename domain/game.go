package domain

type GameUsecase interface {
	AddPlayer(player PlayerUsecase)
	Finish()
}
