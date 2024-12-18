package game

import (
	"math/rand"
	"sort"
	"sync"
)

func (s *Snake) New(width, height int) {
	s.Body = append(s.Body, Position{X: width / 2, Y: height/2 + 1}, Position{X: width/2 + 1, Y: height/2 + 1})
	s.Score = 0
	s.Width = width
	s.Height = height
	s.IgnoreTail = false
	s.Directions = []Position{{X: 1, Y: 0}}
	s.InitialLength = len(s.Body)
	s.TurnsWithoutEating = 0

	s.NewFood()
}

func (s *Snake) GetHead() Position {
	return s.Body[0]
}

func (s *Snake) GetTail() Position {
	return s.Body[len(s.Body)-1]
}

func (s *Snake) Move(direction Position) {
	newHead := Position{
		X: s.GetHead().X + direction.X,
		Y: s.GetHead().Y + direction.Y,
	}
	s.Directions = append(s.Directions, direction)

	s.Body = append([]Position{newHead}, s.Body[:len(s.Body)-1]...)
}

func (s *Snake) grow() {
	s.Body = append(s.Body, s.GetTail())
}

func (s *Snake) Copy() Snake {
	body := make([]Position, len(s.Body))
	copy(body, s.Body)
	return Snake{
		Width:      s.Width,
		Height:     s.Height,
		IgnoreTail: s.IgnoreTail,
		Food:       s.Food,
		Body:       body,
		Score:      s.Score,
		Directions: s.Directions,
	}
}

func (s *Snake) isCollidesWithSelf() bool {
	for i, p := range s.Body[1:] {
		if i == 0 {
			continue
		}
		if p == s.GetHead() {
			return true
		}
	}
	return false
}

func (s *Snake) isCollidesWithWall() bool {
	head := s.GetHead()
	return head.X < 0 || head.X >= s.Width || head.Y < 0 || head.Y >= s.Height
}

func (s *Snake) IsDead() bool {
	if s.isCollidesWithSelf() {
		return true
	}
	if s.isCollidesWithWall() {
		return true
	}
	return false
}

func (s *Snake) HaveWon() bool {
	return s.Score == s.Width*s.Height-s.InitialLength
}

func (s *Snake) IsEating() bool {
	return s.GetHead() == s.Food
}

// Eat increases the score of the snake and makes it grow
func (s *Snake) Eat() {
	s.grow()
	s.TurnsWithoutEating = 0
	s.Score++
}

func (s *Snake) getPathTo(pos Position) []Position {
	grid := initGrid(s.Width, s.Height, s.Body)
	if s.IgnoreTail {
		grid[s.GetTail().Y][s.GetTail().X] = true
	}
	return bfs(grid, s.GetHead(), pos, s.Width, s.Height)
}

func (s *Snake) NewFood() {
	grid := initGrid(s.Width, s.Height, s.Body)
	possible := make([]Position, 0)
	for i := 0; i < s.Width; i++ {
		for j := 0; j < s.Height; j++ {
			if grid[j][i] {
				possible = append(possible, Position{X: i, Y: j})
			}
		}
	}
	if len(possible) > 1 {
		s.Food = possible[rand.Intn(len(possible))]
	}
}

func (s *Snake) getAccessibleDirection() []Position {
	accessibles := make([]Position, 0)
	for _, d := range directions {
		snakeCopy := s.Copy()
		snakeCopy.Move(d)
		if !snakeCopy.IsDead() {
			accessibles = append(accessibles, d)
		}
	}
	return accessibles
}

// getLongestPathToTail returns the longest path to the tail of the snake
func (s *Snake) getLongestPathToTail() []Position {
	var wg sync.WaitGroup

	s.TurnsWithoutEating++
	resultCh := make(chan []Position)

	// Start a goroutine for each accessible direction
	for _, d := range s.getAccessibleDirection() {
		wg.Add(1)
		go func(direction Position) {
			defer wg.Done()

			snakeCopy := s.Copy()
			snakeCopy.Move(direction)
			if snakeCopy.IsEating() {
				snakeCopy.Eat()
			}
			if !snakeCopy.IsDead() {
				snakeCopy.IgnoreTail = true
				defer func() { snakeCopy.IgnoreTail = false }()
				pathToTail := snakeCopy.getPathTo(snakeCopy.GetTail())

				if len(pathToTail) > 0 && len(s.Body) > 2 {
					headPos := snakeCopy.GetHead()
					nextMove := Position{
						X: pathToTail[0].X - headPos.X,
						Y: pathToTail[0].Y - headPos.Y,
					}
					snakeCopy.Move(nextMove)
					if snakeCopy.IsEating() {
						snakeCopy.Eat()
					}
					if !snakeCopy.IsDead() {
						resultCh <- append([]Position{headPos}, pathToTail...)
					}
				}
			}
		}(d)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	tailPaths := [][]Position{}
	for path := range resultCh {
		tailPaths = append(tailPaths, path)
	}

	// If no path found, return empty path
	if len(tailPaths) == 0 {
		return []Position{}
	}

	sort.Slice(tailPaths, func(i, j int) bool {
		return len(tailPaths[i]) > len(tailPaths[j])
	})

	// If snake is stuck in loop, return a random direction to tail
	if s.TurnsWithoutEating > s.Width*s.Height*4 {
		return tailPaths[rand.Intn(len(tailPaths))]
	}
	// Return the longest direction to tail
	return tailPaths[0]
}

func (s *Snake) GetNextDirection() Position {
	snakeCopy := s.Copy()
	// Quiker way to get the path to the food
	path := snakeCopy.getPathTo(s.Food)

	var pathToTail []Position
	// For each path to food, try to get the path to tail
	if len(path) > 0 && s.TurnsWithoutEating != 0 {
		for p := range path {
			d := Position{X: path[p].X - snakeCopy.GetHead().X, Y: path[p].Y - snakeCopy.GetHead().Y}
			snakeCopy.Move(d)
		}
		snakeCopy.Eat()
		snakeCopy.IgnoreTail = true
		pathToTail = snakeCopy.getPathTo(snakeCopy.GetTail())
		defer func() {
			snakeCopy.IgnoreTail = false
		}()

		if len(pathToTail) > 0 && len(s.Body) > 2 {
			d := Position{X: pathToTail[0].X - snakeCopy.GetHead().X, Y: pathToTail[0].Y - snakeCopy.GetHead().Y}
			snakeCopy.Move(d)
			if snakeCopy.IsDead() {
				// If next move leads to win even if it leads to death, do it
				if s.Score == s.Width*s.Height-s.InitialLength-1 {
					return Position{X: path[0].X - s.GetHead().X, Y: path[0].Y - s.GetHead().Y}
				}
				pathToTail = []Position{}
			}
		}
	}

	// Dead end
	if len(path) > 0 && len(pathToTail) > 0 {
		d := Position{X: path[0].X - s.GetHead().X, Y: path[0].Y - s.GetHead().Y}
		return d
	}

	// If no path to tail, try to get the longest path to tail
	path = s.getLongestPathToTail()
	if len(path) > 0 {
		d := Position{X: path[0].X - s.GetHead().X, Y: path[0].Y - s.GetHead().Y}
		return d
	}

	// If no path to tail, try to get the longest path to tail
	path = s.getAccessibleDirection()
	if len(path) > 0 {
		return path[rand.Intn(len(path))]
	}

	// Probably dead
	return Position{}
}

func (s *Snake) GetCurrentDirection() Position {
	head := s.GetHead()
	next := s.Body[1]
	return Position{X: head.X - next.X, Y: head.Y - next.Y}
}

func (s *Snake) GetHeadRune() rune {
	return '▇'
}

// GetPreviousDirection returns the previous direction of the snake
func (s *Snake) GetPreviousDirection() Position {
	if len(s.Directions) < 2 {
		return Position{
			X: 1,
			Y: 0,
		}
	}
	prev := s.Directions[len(s.Directions)-2]
	return prev
}

// GetBodyRune returns the rune to draw the body of the snake
func (s *Snake) GetBodyRune() rune {
	if len(s.Body) < 2 {
		return '?'
	}

	prevDir := directionOf(s.GetPreviousDirection())
	curDir := directionOf(s.GetCurrentDirection())

	if r, ok := bodyRunes[[2]int{prevDir, curDir}]; ok {
		return r
	}
	return '?'
}
