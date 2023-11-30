package stubs

var ProcessGameOfLife = "GameOfLife.ProcessAllTurns"
var ProcessAliveCells = "GameOfLife.ProcessAliveCells"
var PauseGame = "GameOfLife.PauseGameProcess"
var UnPauseGame = "GameOfLife.UnPauseGameProcess"
var WorkerProcess = "Worker.WorkerProcess"

type Params struct {
	Turns   int
	Threads int
	Height  int
	Width   int
}

type Request struct {
	World  [][]byte
	Params Params
	StartY int
	EndY   int
}

type Response struct {
	NextWorld     [][]byte
	CompletedTurn int
	InstantWorld  [][]byte
}
