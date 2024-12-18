package game

type Snake struct {
	Body               []Position
	Score              int
	IgnoreTail         bool
	Food               Position
	Width              int
	Height             int
	Directions         []Position
	TurnsWithoutEating int
	InitialLength      int
}

type node struct {
	pos  Position
	prev *node
}

type Position struct {
	X, Y int
}
