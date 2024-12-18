package terminal

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ctrl-vfr/snake-ai-tui/internal/game"
)

type Screen struct {
	Width      int
	Height     int
	GameWidth  int
	GameHeight int
	OffSetX    int
	OffSetY    int
}

func (s *Screen) New(gameWidth, gameHeight int) error {
	w, h := getScreensize()

	if w < gameWidth || h < gameHeight-4 {
		return fmt.Errorf("screen size is too small. Max are %d*%d", w, h-4)
	}

	s.Width = w
	s.Height = h

	if gameWidth == 0 || gameHeight == 0 {
		s.GameWidth = s.Width / 2
		s.GameHeight = s.Height / 2
	} else {
		s.GameWidth = gameWidth
		s.GameHeight = gameHeight
	}

	if gameWidth == 0 || gameHeight == 0 {
		s.OffSetX = (w - s.GameWidth) / 2
		s.OffSetY = (h - s.GameHeight) / 2
	} else {
		s.OffSetX = (w - gameWidth) - (w-gameWidth)/2
		s.OffSetY = (h - gameHeight) - (h-gameHeight)/2
	}

	return nil

}

func (s *Screen) DrawFrame() {
	moveCursor(s.OffSetX-1, s.OffSetY-1)
	printText("┏", BLUE)
	for i := 0; i < s.GameWidth; i++ {
		printText("━", BLUE)
	}
	printText("┓", BLUE)

	for j := 0; j < s.GameHeight; j++ {
		moveCursor(s.OffSetX-1, s.OffSetY+j)
		printText("┃", BLUE)
		moveCursor(s.OffSetX+s.GameWidth, s.OffSetY+j)
		printText("┃", BLUE)
	}

	moveCursor(s.OffSetX-1, s.OffSetY+s.GameHeight)
	printText("┗", BLUE)
	for i := 0; i < s.GameWidth; i++ {
		printText("━", BLUE)
	}
	printText("┛", BLUE)
}

func (s *Screen) DrawRune(p game.Position, r rune, color uint8) {
	moveCursor(s.OffSetX+p.X, s.OffSetY+p.Y)
	c := getColor(color)
	fmt.Printf("\x1b[%dm%c\x1b[0m", int(c), r)
}

func (s *Screen) EraseRune(p game.Position) {
	moveCursor(p.X, p.Y)
	fmt.Print(" ")
}

func (s *Screen) DrawScore(score int) {
	moveCursor(s.OffSetX, s.OffSetY+s.GameHeight+1)
	printText(fmt.Sprintf("Score: %d", score), BLUE)
}

func (s *Screen) DrawEndScreen(text string, color uint8) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	visited := make([][]bool, s.GameHeight)
	for i := range visited {
		visited[i] = make([]bool, s.GameWidth)
	}

	totalCells := s.GameWidth * s.GameHeight

	maxConcurrent := 10
	sem := make(chan struct{}, maxConcurrent)

	for i := 0; i < totalCells; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			for {
				x := rand.Intn(s.GameWidth)
				y := rand.Intn(s.GameHeight)

				mu.Lock()
				if visited[y][x] {
					mu.Unlock()
					continue
				}
				visited[y][x] = true
				mu.Unlock()

				// Dessine la rune
				moveCursor(s.OffSetY+y, s.OffSetX+x)
				s.DrawRune(game.Position{X: x, Y: y}, '▇', color)
				time.Sleep(time.Duration(rand.Intn(4)) * 10 * time.Millisecond) // Pause pour effet visuel
				break
			}
		}()
	}

	wg.Wait()

	textLen := len(text)
	moveCursor(s.OffSetX+s.GameWidth/2-textLen/2, s.OffSetY-2)
	printText(text, color)
}

func (s *Screen) DrawPause() {
	text := "PAUSE"
	moveCursor(s.OffSetX+s.GameWidth/2-len(text)/2, s.OffSetY-2)
	printText(text, YELLOW)
}

func (s *Screen) ErasePause() {
	text := "PAUSE"
	moveCursor(s.OffSetX+s.GameWidth/2-len(text)/2, s.OffSetY-2)
	fmt.Print("       ")
}

func (s *Screen) ClearScreen() {
	fmt.Print("\033[H\033[2J") // Escape code pour effacer l'écran
}
