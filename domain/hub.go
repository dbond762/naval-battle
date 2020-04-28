package domain

type HubUsecase interface {
	Register(GameUsecase)
	Unregister(GameUsecase)
}
