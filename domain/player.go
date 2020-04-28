package domain

type Result string
type Action string

const (
	ResultMiss   = Result("miss")
	ResultHit    = Result("hit")
	ResultDead   = Result("dead")
	ResultFinish = Result("finish")

	ActionStart    = Action("start")
	ActionFinish   = Action("finish")
	ActionGetBoard = Action("getBoard")
	ActionDoStep   = Action("doStep")
)

type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type (
	CoordsMsg struct {
		Coords Coords `json:"coords"`
	}

	ResultMsg struct {
		Result Result `json:"result"`
	}

	ActionMsg struct {
		Action Action `json:"action"`
	}
)

type PlayerUsecase interface {
	DoAction(action Action) error
	DoStep(coords chan Coords, result chan Result)
	EnemyStep(coords chan Coords, result chan Result)
}
