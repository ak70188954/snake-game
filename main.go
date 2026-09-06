package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 800
	screenHeight = 600
	gridWidth    = 40
	gridHeight   = 30
	cellSize     = 20
)

// fillScreen fills the screen with a color
func fillScreen(screen *ebiten.Image, c color.Color) {
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			screen.Set(x, y, c)
		}
	}
}

// Game represents the full game with graphics
type Game struct {
	snake  *Snake
	quit   bool
	frames int
}

// NewGame creates a new game instance
func NewGame() *Game {
	snake := NewSnake(gridWidth, gridHeight, cellSize)
	return &Game{
		snake: snake,
	}
}

// Update handles game logic
func (g *Game) Update() error {
	if g.quit {
		return nil
	}

	// Handle input
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.snake.ChangeDirection(" ")
	}

	if g.snake.GetGameMode() == ModePlaying {
		// Arrow keys and WASD
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
			g.snake.ChangeDirection("up")
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
			g.snake.ChangeDirection("down")
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
			g.snake.ChangeDirection("left")
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.snake.ChangeDirection("right")
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.quit = true
			return nil
		}

		g.snake.Update()
	}

	g.frames++
	return nil
}

// drawGrid draws the game grid
func (g *Game) drawGrid(screen *ebiten.Image) {
	gridX := (screenWidth - gridWidth*cellSize) / 2
	gridY := 80

	// Draw grid background
	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			posX := gridX + x*cellSize
			posY := gridY + y*cellSize

			// Checkerboard pattern
			if (x+y)%2 == 0 {
				ebitenutil.DrawRect(screen, float64(posX), float64(posY), float64(cellSize), float64(cellSize), color.RGBA{34, 34, 34, 255})
			} else {
				ebitenutil.DrawRect(screen, float64(posX), float64(posY), float64(cellSize), float64(cellSize), color.RGBA{40, 40, 40, 255})
			}
		}
	}

	// Draw grid border
	ebitenutil.DrawRect(screen, float64(gridX), float64(gridY), float64(gridWidth*cellSize), 2, color.RGBA{100, 100, 100, 255})
	ebitenutil.DrawRect(screen, float64(gridX), float64(gridY+gridHeight*cellSize-2), float64(gridWidth*cellSize), 2, color.RGBA{100, 100, 100, 255})
	ebitenutil.DrawRect(screen, float64(gridX), float64(gridY), 2, float64(gridHeight*cellSize), color.RGBA{100, 100, 100, 255})
	ebitenutil.DrawRect(screen, float64(gridX+gridWidth*cellSize-2), float64(gridY), 2, float64(gridHeight*cellSize), color.RGBA{100, 100, 100, 255})
}

// drawSnake draws the snake on screen
func (g *Game) drawSnake(screen *ebiten.Image) {
	body := g.snake.GetBody()
	gridX := (screenWidth - gridWidth*cellSize) / 2
	gridY := 80

	for i, segment := range body {
		posX := gridX + segment.X*cellSize
		posY := gridY + segment.Y*cellSize

		// Head is brighter
		if i == 0 {
			// Head gradient
			headR, headG, headB := 0, 255, 0
			ebitenutil.DrawRect(screen, float64(posX), float64(posY), cellSize, cellSize, color.RGBA{byte(headR), byte(headG), byte(headB), 255})

			// Draw eyes
			eyeSize := 3
			eyeOffset := 5
			dir := g.snake.GetDirection()

			var eye1X, eye1Y, eye2X, eye2Y int
			switch dir {
			case Up:
				eye1X, eye1Y = posX+eyeOffset, posY+eyeOffset
				eye2X, eye2Y = posX+cellSize-eyeOffset-eyeSize, posY+eyeOffset
			case Down:
				eye1X, eye1Y = posX+eyeOffset, posY+cellSize-eyeOffset-eyeSize
				eye2X, eye2Y = posX+cellSize-eyeOffset-eyeSize, posY+cellSize-eyeOffset-eyeSize
			case Left:
				eye1X, eye1Y = posX+eyeOffset, posY+eyeOffset
				eye2X, eye2Y = posX+eyeOffset, posY+cellSize-eyeOffset-eyeSize
			case Right:
				eye1X, eye1Y = posX+cellSize-eyeOffset-eyeSize, posY+eyeOffset
				eye2X, eye2Y = posX+cellSize-eyeOffset-eyeSize, posY+cellSize-eyeOffset-eyeSize
			}

			ebitenutil.DrawRect(screen, float64(eye1X), float64(eye1Y), float64(eyeSize), float64(eyeSize), color.RGBA{255, 255, 255, 255})
			ebitenutil.DrawRect(screen, float64(eye2X), float64(eye2Y), float64(eyeSize), float64(eyeSize), color.RGBA{255, 255, 255, 255})
		} else {
			// Body gradient - gets darker towards tail
			t := float64(i) / float64(len(body))
			green := uint8(float64(255) * (1.0 - t*0.5))
			ebitenutil.DrawRect(screen, float64(posX), float64(posY), float64(cellSize), float64(cellSize), color.RGBA{0, green, 0, 255})

			// Body pattern
			if i%2 == 0 {
				ebitenutil.DrawRect(screen, float64(posX)+2, float64(posY)+2, float64(cellSize-4), float64(cellSize-4), color.RGBA{0, green-30, 0, 100})
			}
		}
	}
}

// drawFood draws the food on screen
func (g *Game) drawFood(screen *ebiten.Image) {
	food := g.snake.GetFood()
	gridX := (screenWidth - gridWidth*cellSize) / 2
	gridY := 80

	posX := gridX + food.X*cellSize + 2
	posY := gridY + food.Y*cellSize + 2
	size := cellSize - 4

	// Apple shape with shine
	ebitenutil.DrawRect(screen, float64(posX), float64(posY), float64(size), float64(size), color.RGBA{255, 0, 0, 255})
	ebitenutil.DrawRect(screen, float64(posX)+2, float64(posY)+2, float64(size-4), float64(size-4), color.RGBA{255, 50, 50, 255})

	// Shine effect
	ebitenutil.DrawRect(screen, float64(posX)+3, float64(posY)+3, 4, 4, color.RGBA{255, 255, 255, 200})
}

// drawHUD draws the heads-up display
func (g *Game) drawHUD(screen *ebiten.Image) {
	// Score
	scoreText := fmt.Sprintf("Score: %d", g.snake.GetScore())
	text.Draw(screen, scoreText, basicfont.Face7x13, 20, 30, color.White)

	// Length
	length := len(g.snake.GetBody())
	lengthText := fmt.Sprintf("Length: %d", length)
	text.Draw(screen, lengthText, basicfont.Face7x13, 200, 30, color.White)

	// Speed
	speed := (200 - int(g.snake.GetSpeed().Milliseconds())) / 5
	speedText := fmt.Sprintf("Speed: %d", speed)
	text.Draw(screen, speedText, basicfont.Face7x13, 400, 30, color.White)
}

// drawMenu draws the main menu
func (g *Game) drawMenu(screen *ebiten.Image) {
	fillScreen(screen, color.RGBA{20, 20, 30, 255})

	// Title
	title := "SNAKE GAME"
	titleX := (screenWidth - len(title)*14) / 2
	text.Draw(screen, title, basicfont.Face7x13, titleX, 150, color.RGBA{0, 255, 0, 255})

	// Subtitle
	subtitle := "A Go + Ebiten Game"
	subX := (screenWidth - len(subtitle)*7) / 2
	text.Draw(screen, subtitle, basicfont.Face7x13, subX, 180, color.RGBA{200, 200, 200, 255})

	// Instructions
	instructions := []string{
		"",
		"Controls:",
		"Arrow Keys / WASD - Move",
		"Space / Enter - Start",
		"Escape - Quit",
		"",
		"Eat red apples to grow!",
		"Don't hit walls or yourself!",
		"",
		"Press SPACE to start",
	}

	y := 250
	for _, line := range instructions {
		lineX := (screenWidth - len(line)*7) / 2
		if line == "Press SPACE to start" {
			// Blinking effect
			if g.frames/30%2 == 0 {
				text.Draw(screen, line, basicfont.Face7x13, lineX, y, color.RGBA{255, 255, 0, 255})
			} else {
				text.Draw(screen, line, basicfont.Face7x13, lineX, y, color.RGBA{100, 100, 100, 255})
			}
		} else if line == "Controls:" {
			text.Draw(screen, line, basicfont.Face7x13, lineX, y, color.RGBA{100, 200, 255, 255})
		} else {
			text.Draw(screen, line, basicfont.Face7x13, lineX, y, color.White)
		}
		y += 25
	}
}

// drawGameOver draws the game over screen
func (g *Game) drawGameOver(screen *ebiten.Image) {
	fillScreen(screen, color.RGBA{30, 10, 10, 255})

	// Game Over text
	gameOver := "GAME OVER"
	goX := (screenWidth - len(gameOver)*14) / 2
	text.Draw(screen, gameOver, basicfont.Face7x13, goX, 150, color.RGBA{255, 0, 0, 255})

	// Score
	scoreText := fmt.Sprintf("Final Score: %d", g.snake.GetScore())
	scoreX := (screenWidth - len(scoreText)*7) / 2
	text.Draw(screen, scoreText, basicfont.Face7x13, scoreX, 220, color.White)

	// Length
	length := len(g.snake.GetBody())
	lengthText := fmt.Sprintf("Snake Length: %d", length)
	lenX := (screenWidth - len(lengthText)*7) / 2
	text.Draw(screen, lengthText, basicfont.Face7x13, lenX, 250, color.White)

	// Rating
	rating := ""
	score := g.snake.GetScore()
	switch {
	case score < 5:
		rating = "Rating: Beginner"
	case score < 10:
		rating = "Rating: Novice"
	case score < 20:
		rating = "Rating: Skilled"
	case score < 30:
		rating = "Rating: Expert"
	default:
		rating = "Rating: Legend!"
	}
	ratingX := (screenWidth - len(rating)*7) / 2
	text.Draw(screen, rating, basicfont.Face7x13, ratingX, 300, color.RGBA{255, 215, 0, 255})

	// Restart prompt
	if g.frames/30%2 == 0 {
		restart := "Press SPACE to play again"
		restX := (screenWidth - len(restart)*7) / 2
		text.Draw(screen, restart, basicfont.Face7x13, restX, 400, color.RGBA{255, 255, 0, 255})
	}

	// Quit prompt
	quit := "Press ESC to quit"
	quitX := (screenWidth - len(quit)*7) / 2
	text.Draw(screen, quit, basicfont.Face7x13, quitX, 430, color.RGBA{200, 200, 200, 255})
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.snake.GetGameMode() {
	case ModeMenu:
		g.drawMenu(screen)
	case ModePlaying:
		// Background
	fillScreen(screen, color.RGBA{20, 20, 30, 255})

		// Draw grid
		g.drawGrid(screen)

		// Draw food
		g.drawFood(screen)

		// Draw snake
		g.drawSnake(screen)

		// Draw HUD
		g.drawHUD(screen)

	case ModeGameOver:
		g.drawGameOver(screen)
	}
}

// Layout returns the game size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake Game - Go + Ebiten")
	ebiten.SetWindowResizable(true)

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
