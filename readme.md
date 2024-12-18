# 🐍 Snake TUI

A simple and interactive **Terminal User Interface (TUI)** Snake game developed in **Go**. The game offers two modes:
- **Player**: Manual control of the snake using arrow keys.
- **Machine**: The snake plays automatically, aiming to maximize its score.

## 🚀 Features

1. **Two game modes**:
   - **Player Mode**: Use arrow keys to control the snake.
   - **Machine Mode**: An algorithm guides the snake to avoid obstacles and reach the food.

## Demonstration 🎮
Here's a demonstration of the game in action:

![Snake TUI Demo](demo.gif)

## 🛠️ Installation

### Prerequisites
- **Go** (version 1.21 or higher) installed on your machine.

### Clone the project
```bash
# Clone the repository
git clone https://github.com/ctrl-vfr/snake-ai-tui.git
cd snake-ai-tui
# Install dependencies
go mod tidy
# Build the project
go build ./cmd/snake-ai-tui
```

### Alternative installation

Alternatively, you can use the `go install` command to install the project directly from GitHub.

```bash
go install github.com/ctrl-vfr/snake-ai-tui/cmd/snake-ai-tui
```

## 🎮 Usage

### Launch the game
#### Player Mode:
```bash
./snake-tui -m 0
```
#### Machine Mode:
```bash
./snake-tui -m 1
```

### Available options:
| **Option**       | **Description**                              | **Default Value**      |
|------------------|----------------------------------------------|------------------------|
| `-m`            | Game mode: 0 = Player, 1 = Machine            | `1` (Machine)         |
| `-h`            | Grid height                                   | Dynamic size           |
| `-w`            | Grid width                                    | Dynamic size           |
| `-s`            | Game speed (in ms, screen refresh rate)       | `30` (Machine), `200` (Player) |
| `-h`            | Display help                                  |                        |

### In-game controls

#### Player Mode:

| **Key**         | **Action**                |
|-----------------|--------------------------|
| Arrow keys      | Move the snake            |
| `SPACE`         | Pause the game            |
| `CTRL+C`        | Exit the game             |

#### Machine Mode:

| **Key**         | **Action**                |
|-----------------|--------------------------|
| `SPACE`         | Pause the game            |
| `CTRL+C`        | Exit the game             |


## 📜 License

This project is licensed under the **MIT** license.