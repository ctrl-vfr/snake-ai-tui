package game

const (
	UP = iota
	DOWN
	LEFT
	RIGHT
	INVALID
)

var bodyRunes = map[[2]int]rune{
	{UP, UP}:       '║',
	{DOWN, DOWN}:   '║',
	{LEFT, LEFT}:   '═',
	{RIGHT, RIGHT}: '═',

	{UP, RIGHT}:   '╔',
	{UP, LEFT}:    '╗',
	{DOWN, RIGHT}: '╚',
	{DOWN, LEFT}:  '╝',

	{LEFT, UP}:    '╚',
	{LEFT, DOWN}:  '╔',
	{RIGHT, UP}:   '╝',
	{RIGHT, DOWN}: '╗',
}

var directions = []Position{
	{X: 0, Y: -1}, // UP
	{X: 0, Y: 1},  // DOWN
	{X: -1, Y: 0}, // LEFT
	{X: 1, Y: 0},  // RIGHT
}
