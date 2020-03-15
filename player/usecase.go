package player

type Usecase interface {
	Step(coords chan<- []byte, result <-chan []byte) error
	EnemyStep(coords <-chan []byte, result chan<- []byte) error
}
