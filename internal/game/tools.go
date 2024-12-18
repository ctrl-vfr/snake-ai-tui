package game

func initGrid(width, height int, snake []Position) [][]bool {
	grid := make([][]bool, height)
	for y := range grid {
		grid[y] = make([]bool, width)
		for x := range grid[y] {
			grid[y][x] = true
		}
	}

	for _, seg := range snake {
		if seg.Y >= 0 && seg.Y < height && seg.X >= 0 && seg.X < width {
			grid[seg.Y][seg.X] = false
		}
	}

	return grid
}

// Find the shortest path between two points on a grid
func bfs(grid [][]bool, start, goal Position, width, height int) []Position {
	queue := []node{{pos: start}}
	visited := initVisited(width, height)

	if len(visited) <= start.Y || len(visited[0]) <= start.X {
		return []Position{}
	}

	visited[start.Y][start.X] = true

	var endNode *node

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.pos == goal {
			endNode = &current
			break
		}

		for _, d := range directions {
			nextPos := Position{X: current.pos.X + d.X, Y: current.pos.Y + d.Y}

			if !isValidCell(nextPos, width, height) {
				continue
			}
			if !grid[nextPos.Y][nextPos.X] {
				continue
			}
			if visited[nextPos.Y][nextPos.X] {
				continue
			}

			visited[nextPos.Y][nextPos.X] = true
			queue = append(queue, node{
				pos:  nextPos,
				prev: &current,
			})
		}
	}

	if endNode == nil {
		return []Position{}
	}

	path := reconstructPath(endNode)
	return path[1:]
}

// Initialize the visited grid
func initVisited(width, height int) [][]bool {
	visited := make([][]bool, height)
	for i := range visited {
		visited[i] = make([]bool, width)
	}
	return visited
}

func isValidCell(p Position, width, height int) bool {
	return p.X >= 0 && p.X < width && p.Y >= 0 && p.Y < height
}

func reconstructPath(endNode *node) []Position {
	var path []Position
	for n := endNode; n != nil; n = n.prev {
		path = append(path, n.pos)
	}
	reversePositions(path)
	return path
}

func reversePositions(path []Position) {
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
}

// Find the direction code of a position
func directionOf(p Position) int {
	switch p {
	case Position{X: 0, Y: -1}:
		return UP
	case Position{X: 0, Y: 1}:
		return DOWN
	case Position{X: -1, Y: 0}:
		return LEFT
	case Position{X: 1, Y: 0}:
		return RIGHT
	default:
		return INVALID
	}
}
