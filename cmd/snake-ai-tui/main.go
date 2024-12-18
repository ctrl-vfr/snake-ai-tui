package main

import (
	"flag"
	"fmt"
	"os"
	"snake/internal/game"
	"snake/internal/terminal"
	"time"
)

const (
	USER = iota
	MACHINE
)

var (
	mode   = USER
	height = 0
	width  = 0
	speed  int64
)

func printHelp() {
	fmt.Println("Snake game")
	fmt.Println("Usage: snake-ai-tui [options]")
	fmt.Println("Options:")
	flag.PrintDefaults()

	fmt.Println("In game as a user:")
	fmt.Println("  - Use arrow keys to move the snake")
	fmt.Println("  - CTRL+C to quit")
	fmt.Println("  - SPACE to pause")

	fmt.Println("In game as a machine:")
	fmt.Println("  - CTRL+C to quit")
	fmt.Println("  - SPACE to pause")

	os.Exit(0)
}

func init() {
	// def flags
	flag.IntVar(&mode, "m", MACHINE, "Game mode (0: User, 1: Machine)")
	flag.IntVar(&height, "h", 0, "Game's frame height")
	flag.IntVar(&width, "w", 0, "Game's frame width")
	flag.Int64Var(&speed, "s", 0, "Game speed (refresh rate in ms)")
	flag.Bool("help", false, "Display help")
	flag.Parse()
}

func main() {
	if flag.Lookup("help").Value.String() == "true" {
		printHelp()
	}

	fmt.Print("\033[2J")
	previousTerminalState := terminal.SetupTerminal()
	defer terminal.ResetTerminal(previousTerminalState)

	terminal.HideCursor()
	defer terminal.ShowCursor()

	screen := terminal.Screen{}
	err := screen.New(width, height)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	screen.DrawFrame()

	snake := game.Snake{}
	snake.New(screen.GameWidth, screen.GameHeight)

	snake.NewFood()
	screen.DrawRune(snake.Food, '◈', terminal.RED)

	screen.DrawRune(snake.GetHead(), snake.GetHeadRune(), terminal.GREEN)
	screen.DrawScore(snake.Score)

	rate := time.Duration(30)
	if speed <= 0 && mode == USER {
		rate = time.Duration(200)
	} else if speed != 0 {
		rate = time.Duration(speed)
	}

	ticker := time.NewTicker(rate * time.Millisecond)
	defer ticker.Stop()

	end := false
	pause := false
	// Gestion des entrées utilisateur
	inputChan := make(chan string)
	go terminal.ReadInput(inputChan)

	direction := game.Position{X: 1, Y: 0}

	go func() {
		for {
			input := <-inputChan
			switch input {
			case "QUIT":
				end = true
			case "PAUSE":
				if pause {
					screen.ErasePause()
				} else {
					screen.DrawPause()
				}
				pause = !pause
			case "UP":
				if mode == USER {
					direction = game.Position{X: 0, Y: -1}
				}
			case "DOWN":
				if mode == USER {
					direction = game.Position{X: 0, Y: 1}
				}
			case "LEFT":
				if mode == USER {
					direction = game.Position{X: -1, Y: 0}
				}
			case "RIGHT":
				if mode == USER {
					direction = game.Position{X: 1, Y: 0}
				}
			}
		}
	}()

	for !end {
		<-ticker.C

		for pause {
			time.Sleep(400 * time.Millisecond)
		}

		tail := snake.GetTail()

		if mode == MACHINE {
			direction = snake.GetNextDirection()
			if direction == (game.Position{}) {
				direction = snake.GetPreviousDirection()
			}
		}

		snake.Move(direction)

		if snake.IsEating() {
			snake.Eat()
			snake.NewFood()
			screen.DrawRune(snake.Food, '◈', terminal.RED)
			screen.DrawScore(snake.Score)
		}
		screen.EraseRune(game.Position{X: tail.X + screen.OffSetX, Y: tail.Y + screen.OffSetY})

		if len(snake.Body) > 1 {
			screen.DrawRune(snake.Body[1], snake.GetBodyRune(), terminal.GREEN)
		}
		screen.DrawRune(snake.GetHead(), snake.GetHeadRune(), terminal.GREEN)

		if snake.HaveWon() {
			screen.DrawEndScreen("You won!", terminal.GREEN)
			time.Sleep(1 * time.Second)
			end = true
		}

		if snake.IsDead() {
			screen.DrawEndScreen("You lose!", terminal.RED)
			time.Sleep(1 * time.Second)
			end = true
		}

		if snake.TurnsWithoutEating > snake.Width*snake.Height*30 {
			screen.DrawEndScreen("You're starving!", terminal.YELLOW)
			time.Sleep(1 * time.Second)
			end = true
		}

	}
	screen.DrawEndScreen("Goodbye! Thanks for playing! 🐍", terminal.BLUE)
	time.Sleep(1 * time.Second)
	screen.ClearScreen()

}
