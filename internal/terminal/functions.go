package terminal

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// Commandes ANSI
func moveCursor(x, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

func HideCursor() {
	fmt.Print("\033[?25l")
}

func ShowCursor() {
	fmt.Print("\033[?25h")
}

func SetupTerminal() *term.State {
	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	return oldState
}

func ResetTerminal(state *term.State) {
	term.Restore(int(os.Stdin.Fd()), state)
}

func ReadInput(inputChan chan string) {
	buffer := make([]byte, 3)

	for {
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			inputChan <- "ERROR"
			continue
		}

		// Detect arrow keys
		if n == 3 && buffer[0] == 0x1b && buffer[1] == 0x5b {
			switch buffer[2] {
			case 0x41:
				inputChan <- "UP"
			case 0x42:
				inputChan <- "DOWN"
			case 0x43:
				inputChan <- "RIGHT"
			case 0x44:
				inputChan <- "LEFT"
			}
		} else if n == 1 { // Detect other keys
			switch buffer[0] {
			case 0x20:
				inputChan <- "PAUSE"
			case 0x03:
				inputChan <- "QUIT"
				return
			default:
				inputChan <- "UNKNOWN"
			}
		}
	}
}

func getScreensize() (int, int) {
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	return width, height
}

func getColor(color uint8) uint8 {
	// Custom colors code to ANSI
	colors := map[uint8]uint8{
		BLACK:   30,
		RED:     31,
		GREEN:   32,
		YELLOW:  33,
		BLUE:    34,
		MAGENTA: 35,
		CYAN:    36,
		WHITE:   37,
	}

	if c, ok := colors[color]; ok {
		return c
	}
	return 37
}

func printText(text string, color uint8) {
	c := getColor(color)
	fmt.Printf("\x1b[%dm%s\x1b[0m", int(c), text)
}
