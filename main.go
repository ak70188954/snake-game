package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	gridWidth  = 40
	gridHeight = 30
	hudHeight  = 50
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// menuFont is a large font for menus (20pt, about 2x the basic 7x13 font)
var menuFont font.Face

func init() {
	ttf, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Println("failed to parse font:", err)
		return
	}
	menuFont = truetype.NewFace(ttf, &truetype.Options{
		Size:    40,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

// Game represents the full game with graphics
type Game struct {
	snake       *Snake
	quit        bool
	frames      int
	cellSize    int
	gridOffsetX int
	gridOffsetY int
}

// NewGame creates a new game instance
func NewGame() *Game {
	snake := NewSnake(gridWidth, gridHeight, 15)
	return &Game{
		snake: snake,
	}
}

// calculateLayout computes cell size and offsets to fit the grid on screen
func (g *Game) calculateLayout(screenWidth, screenHeight int) {
	availableW := screenWidth
	availableH := screenHeight - hudHeight

	maxCellW := availableW / gridWidth
	maxCellH := availableH / gridHeight

	g.cellSize = maxCellW
	if maxCellH < maxCellW {
		g.cellSize = maxCellH
	}
	if g.cellSize < 5 {
		g.cellSize = 5
	}

	g.gridOffsetX = (screenWidth - gridWidth*g.cellSize) / 2
	g.gridOffsetY = hudHeight + (availableH-gridHeight*g.cellSize)/2
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
	cs := g.cellSize
	gx := g.gridOffsetX
	gy := g.gridOffsetY

	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			posX := gx + x*cs
			posY := gy + y*cs

			if (x+y)%2 == 0 {
				ebitenutil.DrawRect(screen, float64(posX), float64(posY), float64(cs), float64(cs), color.RGBA{34, 34, 34, 255})
			} else {
				ebitenutil.DrawRect(screen, float64(posX), float64(posY), float64(cs), float64(cs), color.RGBA{40, 40, 40, 255})
			}
		}
	}

	// Border
	ebitenutil.DrawRect(screen, float64(gx), float64(gy), float64(gridWidth*cs), 2, color.RGBA{100, 100, 100, 255})
	ebitenutil.DrawRect(screen, float64(gx), float64(gy+gridHeight*cs-2), float64(gridWidth*cs), 2, color.RGBA{100, 100, 100, 255})
	ebitenutil.DrawRect(screen, float64(gx), float64(gy), 2, float64(gridHeight*cs), color.RGBA{100, 100, 100, 255})
	ebitenutil.DrawRect(screen, float64(gx+gridWidth*cs-2), float64(gy), 2, float64(gridHeight*cs), color.RGBA{100, 100, 100, 255})
}

// drawSnake draws the snake on screen
func (g *Game) drawSnake(screen *ebiten.Image) {
	body := g.snake.GetBody()
	cs := g.cellSize
	gx := g.gridOffsetX
	gy := g.gridOffsetY

	for i, segment := range body {
		posX := gx + segment.X*cs
		posY := gy + segment.Y*cs

		if i == 0 {
			ebitenutil.DrawRect(screen, float64(posX), float64(posY), float64(cs), float64(cs), color.RGBA{0, 255, 0, 255})

			// Eyes
			eyeSize := max(2, cs/6)
			eyeOffset := cs/4
			dir := g.snake.GetDirection()

			var eye1X, eye1Y, eye2X, eye2Y int
			switch dir {
			case Up:
				eye1X, eye1Y = posX+eyeOffset, posY+eyeOffset
				eye2X, eye2Y = posX+cs-eyeOffset-eyeSize, posY+eyeOffset
			case Down:
				eye1X, eye1Y = posX+eyeOffset, posY+cs-eyeOffset-eyeSize
				eye2X, eye2Y = posX+cs-eyeOffset-eyeSize, posY+cs-eyeOffset-eyeSize
			case Left:
				eye1X, eye1Y = posX+eyeOffset, posY+eyeOffset
				eye2X, eye2Y = posX+eyeOffset, posY+cs-eyeOffset-eyeSize
			case Right:
				eye1X, eye1Y = posX+cs-eyeOffset-eyeSize, posY+eyeOffset
				eye2X, eye2Y = posX+cs-eyeOffset-eyeSize, posY+cs-eyeOffset-eyeSize
			}

			ebitenutil.DrawRect(screen, float64(eye1X), float64(eye1Y), float64(eyeSize), float64(eyeSize), color.RGBA{255, 255, 255, 255})
			ebitenutil.DrawRect(screen, float64(eye2X), float64(eye2Y), float64(eyeSize), float64(eyeSize), color.RGBA{255, 255, 255, 255})
		} else {
			t := float64(i) / float64(len(body))
			green := uint8(float64(255) * (1.0 - t*0.5))
			ebitenutil.DrawRect(screen, float64(posX), float64(posY), float64(cs), float64(cs), color.RGBA{0, green, 0, 255})

			if i%2 == 0 {
				ebitenutil.DrawRect(screen, float64(posX)+1, float64(posY)+1, float64(cs-2), float64(cs-2), color.RGBA{0, green-30, 0, 100})
			}
		}
	}
}

// drawFood draws the food on screen
func (g *Game) drawFood(screen *ebiten.Image) {
	food := g.snake.GetFood()
	cs := g.cellSize
	gx := g.gridOffsetX
	gy := g.gridOffsetY

	posX := gx + food.X*cs + cs/5
	posY := gy + food.Y*cs + cs/5
	size := cs - cs/5*2

	ebitenutil.DrawRect(screen, float64(posX), float64(posY), float64(size), float64(size), color.RGBA{255, 0, 0, 255})
	ebitenutil.DrawRect(screen, float64(posX)+1, float64(posY)+1, float64(size-2), float64(size-2), color.RGBA{255, 50, 50, 255})
	ebitenutil.DrawRect(screen, float64(posX)+1, float64(posY)+1, float64(max(1, size/4)), float64(max(1, size/4)), color.RGBA{255, 255, 255, 200})
}

// drawHUD draws the heads-up display
func (g *Game) drawHUD(screen *ebiten.Image) {
	scoreText := fmt.Sprintf("Score: %d", g.snake.GetScore())
	text.Draw(screen, scoreText, menuFont, 10, 50, color.White)

	length := len(g.snake.GetBody())
	lengthText := fmt.Sprintf("Length: %d", length)
	text.Draw(screen, lengthText, menuFont, 150, 50, color.White)

	speed := (200 - int(g.snake.GetSpeed().Milliseconds())) / 5
	speedText := fmt.Sprintf("Speed: %d", speed)
	text.Draw(screen, speedText, menuFont, 300, 50, color.White)
}

// drawMenu draws the main menu
func (g *Game) drawMenu(screen *ebiten.Image) {
	utilFill(screen, color.RGBA{20, 20, 30, 255})

	w, _ := screen.Size()

	title := "SNAKE GAME"
	titleX := (w - len(title)*14) / 2
	text.Draw(screen, title, menuFont, titleX, 300, color.RGBA{0, 255, 0, 255})

	subtitle := "A Go + Ebiten Game"
	subX := (w - len(subtitle)*7) / 2
	text.Draw(screen, subtitle, menuFont, subX, 360, color.RGBA{200, 200, 200, 255})

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

	y := 500
	for _, line := range instructions {
		lineX := (w - len(line)*7) / 2
		if line == "Press SPACE to start" {
			if g.frames/30%2 == 0 {
				text.Draw(screen, line, menuFont, lineX, y, color.RGBA{255, 255, 0, 255})
			} else {
				text.Draw(screen, line, menuFont, lineX, y, color.RGBA{100, 100, 100, 255})
			}
		} else if line == "Controls:" {
			text.Draw(screen, line, menuFont, lineX, y, color.RGBA{100, 200, 255, 255})
		} else {
			text.Draw(screen, line, menuFont, lineX, y, color.White)
		}
		y += 50
	}
}

// drawGameOver draws the game over screen
func (g *Game) drawGameOver(screen *ebiten.Image) {
	utilFill(screen, color.RGBA{30, 10, 10, 255})

	w, _ := screen.Size()

	gameOver := "GAME OVER"
	goX := (w - len(gameOver)*14) / 2
	text.Draw(screen, gameOver, menuFont, goX, 300, color.RGBA{255, 0, 0, 255})

	scoreText := fmt.Sprintf("Final Score: %d", g.snake.GetScore())
	scoreX := (w - len(scoreText)*7) / 2
	text.Draw(screen, scoreText, menuFont, scoreX, 440, color.White)

	length := len(g.snake.GetBody())
	lengthText := fmt.Sprintf("Snake Length: %d", length)
	lenX := (w - len(lengthText)*7) / 2
	text.Draw(screen, lengthText, menuFont, lenX, 500, color.White)

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
	ratingX := (w - len(rating)*7) / 2
	text.Draw(screen, rating, menuFont, ratingX, 600, color.RGBA{255, 215, 0, 255})

	if g.frames/30%2 == 0 {
		restart := "Press SPACE to play again"
		restX := (w - len(restart)*7) / 2
		text.Draw(screen, restart, menuFont, restX, 800, color.RGBA{255, 255, 0, 255})
	}

	quit := "Press ESC to quit"
	quitX := (w - len(quit)*7) / 2
	text.Draw(screen, quit, menuFont, quitX, 860, color.RGBA{200, 200, 200, 255})
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	screenWidth, screenHeight := screen.Size()
	g.calculateLayout(screenWidth, screenHeight)

	switch g.snake.GetGameMode() {
	case ModeMenu:
		g.drawMenu(screen)
	case ModePlaying:
		utilFill(screen, color.RGBA{20, 20, 30, 255})
		g.drawGrid(screen)
		g.drawFood(screen)
		g.drawSnake(screen)
		g.drawHUD(screen)
	case ModeGameOver:
		g.drawGameOver(screen)
	}
}

// Layout returns the game size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// utilFill fills the screen with a color (replacement for ebiten/util)
func utilFill(screen *ebiten.Image, c color.Color) {
	w, h := screen.Size()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			screen.Set(x, y, c)
		}
	}
}

func main() {
	ebiten.SetWindowTitle("Snake Game - Go + Ebiten")
	ebiten.SetWindowResizable(true)

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
