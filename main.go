package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

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

type particle struct {
	x, y      float64
	vx, vy    float64
	life      int
	maxLife   int
	r, g, b   uint8
}

// Game represents the full game with graphics
type Game struct {
	snake       *Snake
	quit        bool
	frames      int
	cellSize    int
	gridOffsetX int
	gridOffsetY int
	particles   []particle
	scoreFlash  int
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
			g.snake.QueueDirection(Up)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
			g.snake.QueueDirection(Down)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
			g.snake.QueueDirection(Left)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.snake.QueueDirection(Right)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.quit = true
			return nil
		}

		wasEaten := g.snake.WasEaten()
		g.snake.Update()

		if wasEaten {
			g.scoreFlash = 15
			food := g.snake.GetFood()
			cs := g.cellSize
			gx := g.gridOffsetX
			gy := g.gridOffsetY
			for i := 0; i < 20; i++ {
				angle := float64(i) / 20.0 * 2 * math.Pi
				speed := 1.5 + math.Abs(math.Sin(float64(i)*1.3))*2
				g.particles = append(g.particles, particle{
					x:       float64(gx + food.X*cs + cs/2),
					y:       float64(gy + food.Y*cs + cs/2),
					vx:      math.Cos(angle) * speed,
					vy:      math.Sin(angle) * speed,
					life:    30 + int(math.Abs(math.Sin(float64(i))*20)),
					maxLife: 50,
					r:       255,
					g:       uint8(100 + uint(i)*5),
					b:       50,
				})
			}
		}

		if g.scoreFlash > 0 {
			g.scoreFlash--
		}
	}

	// Update particles
	for i := len(g.particles) - 1; i >= 0; i-- {
		p := &g.particles[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += 0.05
		p.life--
		if p.life <= 0 {
			g.particles = append(g.particles[:i], g.particles[i+1:]...)
		}
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
				ebitenutil.DrawRect(screen, float64(posX), float64(posY), float64(cs), float64(cs), color.RGBA{30, 35, 45, 255})
			} else {
				ebitenutil.DrawRect(screen, float64(posX), float64(posY), float64(cs), float64(cs), color.RGBA{35, 40, 52, 255})
			}
		}
	}

	// Glowing blue border
	borderColor := color.RGBA{80, 140, 220, 255}
	ebitenutil.DrawRect(screen, float64(gx), float64(gy), float64(gridWidth*cs), 2, borderColor)
	ebitenutil.DrawRect(screen, float64(gx), float64(gy+gridHeight*cs-2), float64(gridWidth*cs), 2, borderColor)
	ebitenutil.DrawRect(screen, float64(gx), float64(gy), 2, float64(gridHeight*cs), borderColor)
	ebitenutil.DrawRect(screen, float64(gx+gridWidth*cs-2), float64(gy), 2, float64(gridHeight*cs), borderColor)
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
			// Head - brighter green circle
			cx := float64(posX + cs/2)
			cy := float64(posY + cs/2)
			radius := float64(cs) / 2
			g.drawCircle(screen, cx, cy, radius, color.RGBA{20, 240, 60, 255})

			// Eyes with pupils
			eyeR := math.Max(2, float64(cs)/7)
			eyeOff := float64(cs) / 4
			dir := g.snake.GetDirection()
			var e1x, e1y, e2x, e2y float64
			switch dir {
			case Up:
				e1x, e1y = cx-eyeOff, cy-eyeOff
				e2x, e2y = cx+eyeOff, cy-eyeOff
			case Down:
				e1x, e1y = cx-eyeOff, cy+eyeOff
				e2x, e2y = cx+eyeOff, cy+eyeOff
			case Left:
				e1x, e1y = cx-eyeOff, cy-eyeOff
				e2x, e2y = cx-eyeOff, cy+eyeOff
			case Right:
				e1x, e1y = cx+eyeOff, cy-eyeOff
				e2x, e2y = cx+eyeOff, cy+eyeOff
			}
			g.drawCircle(screen, e1x, e1y, eyeR, color.RGBA{255, 255, 255, 255})
			g.drawCircle(screen, e2x, e2y, eyeR, color.RGBA{255, 255, 255, 255})

			// Pupils
			pupilR := math.Max(1, eyeR*0.5)
			pOff := eyeR * 0.3
			switch dir {
			case Up:
				g.drawCircle(screen, e1x, e1y-pOff, pupilR, color.RGBA{30, 30, 30, 255})
				g.drawCircle(screen, e2x, e2y-pOff, pupilR, color.RGBA{30, 30, 30, 255})
			case Down:
				g.drawCircle(screen, e1x, e1y+pOff, pupilR, color.RGBA{30, 30, 30, 255})
				g.drawCircle(screen, e2x, e2y+pOff, pupilR, color.RGBA{30, 30, 30, 255})
			case Left:
				g.drawCircle(screen, e1x-pOff, e1y, pupilR, color.RGBA{30, 30, 30, 255})
				g.drawCircle(screen, e2x-pOff, e2y, pupilR, color.RGBA{30, 30, 30, 255})
			case Right:
				g.drawCircle(screen, e1x+pOff, e1y, pupilR, color.RGBA{30, 30, 30, 255})
				g.drawCircle(screen, e2x+pOff, e2y, pupilR, color.RGBA{30, 30, 30, 255})
			}
		} else {
			// Body - gradient from bright to dark green with highlights
			t := float64(i) / float64(max(1, len(body)-1))
			gc := uint8(int(255) - int(t)*120)
			cx := float64(posX + cs/2)
			cy := float64(posY + cs/2)
			radius := float64(cs) / 2 * 0.95
			g.drawCircle(screen, cx, cy, radius, color.RGBA{0, gc, 0, 255})

			// Highlight on even segments
			if i%2 == 0 {
				hl := color.RGBA{30, gc+40, 20, 80}
				g.drawCircle(screen, cx-1, cy-1, radius*0.6, hl)
			}
		}
	}
}

// drawCircle draws a filled circle
func (g *Game) drawCircle(screen *ebiten.Image, cx, cy, radius float64, c color.Color) {
	r := int(radius)
	if r < 1 {
		r = 1
	}
	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			if dx*dx+dy*dy <= r*r {
				screen.Set(int(cx)+dx, int(cy)+dy, c)
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

	cx := float64(gx + food.X*cs + cs/2)
	cy := float64(gy + food.Y*cs + cs/2)
	baseR := float64(cs) / 2 * 0.75

	// Pulsating glow effect
	pulse := math.Sin(float64(g.frames)*0.08)*0.3 + 1.0
	glowR := baseR * 1.6 * pulse
	g.drawCircle(screen, cx, cy, glowR, color.RGBA{255, 60, 30, uint8(40*pulse)})

	// Apple body
	appleR := baseR * pulse
	g.drawCircle(screen, cx, cy, appleR, color.RGBA{230, 40, 30, 255})

	// Highlight
	hlR := appleR * 0.35
	g.drawCircle(screen, cx-appleR*0.25, cy-appleR*0.3, hlR, color.RGBA{255, 180, 170, 200})

	// Stem
	stemH := int(float64(cs) * 0.25)
	ebitenutil.DrawRect(screen, cx-0.5, cy-appleR, 2, float64(stemH), color.RGBA{80, 50, 20, 255})

	// Leaf
	leafR := float64(cs) * 0.12
	g.drawCircle(screen, cx+2, cy-appleR+float64(stemH)*0.5, leafR, color.RGBA{60, 160, 40, 255})
}

// drawHUD draws the heads-up display
func (g *Game) drawHUD(screen *ebiten.Image) {
	w, _ := screen.Size()
	// HUD background panel
	ebitenutil.DrawRect(screen, 0, 0, float64(w), hudHeight, color.RGBA{15, 18, 28, 255})
	ebitenutil.DrawRect(screen, 0, float64(hudHeight)-2, float64(w), 2, color.RGBA{80, 140, 220, 255})

	scoreText := fmt.Sprintf("Score: %d", g.snake.GetScore())
	c := color.RGBA{255, 255, 255, 255}
	if g.scoreFlash > 0 {
		c = color.RGBA{100, 255, 100, 255}
	}
	text.Draw(screen, scoreText, menuFont, 10, 35, c)

	length := len(g.snake.GetBody())
	lengthText := fmt.Sprintf("Length: %d", length)
	text.Draw(screen, lengthText, menuFont, 150, 35, color.RGBA{200, 200, 200, 255})

	speed := (200 - int(g.snake.GetSpeed().Milliseconds())) / 5
	speedText := fmt.Sprintf("Speed: %d", speed)
	text.Draw(screen, speedText, menuFont, 300, 35, color.RGBA{180, 180, 220, 255})
}

// drawMenu draws the main menu
func (g *Game) drawMenu(screen *ebiten.Image) {
	w, h := screen.Size()
	ebitenutil.DrawRect(screen, 0, 0, float64(w), float64(h), color.RGBA{15, 18, 28, 255})

	title := "SNAKE GAME"
	titleX := (w - len(title)*14) / 2
	text.Draw(screen, title, menuFont, titleX, 250, color.RGBA{60, 220, 80, 255})

	subtitle := "A Go + Ebiten Game"
	subX := (w - len(subtitle)*7) / 2
	text.Draw(screen, subtitle, menuFont, subX, 310, color.RGBA{200, 200, 200, 255})

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

	y := 400
	for _, line := range instructions {
		lineX := (w - len(line)*7) / 2
		if line == "Press SPACE to start" {
			if g.frames/30%2 == 0 {
				text.Draw(screen, line, menuFont, lineX, y, color.RGBA{255, 255, 100, 255})
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
	w, h := screen.Size()
	ebitenutil.DrawRect(screen, 0, 0, float64(w), float64(h), color.RGBA{25, 15, 15, 255})

	gameOver := "GAME OVER"
	goX := (w - len(gameOver)*14) / 2
	text.Draw(screen, gameOver, menuFont, goX, 250, color.RGBA{255, 80, 80, 255})

	scoreText := fmt.Sprintf("Final Score: %d", g.snake.GetScore())
	scoreX := (w - len(scoreText)*7) / 2
	text.Draw(screen, scoreText, menuFont, scoreX, 350, color.White)

	length := len(g.snake.GetBody())
	lengthText := fmt.Sprintf("Snake Length: %d", length)
	lenX := (w - len(lengthText)*7) / 2
	text.Draw(screen, lengthText, menuFont, lenX, 400, color.RGBA{200, 200, 200, 255})

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
	text.Draw(screen, rating, menuFont, ratingX, 500, color.RGBA{255, 215, 0, 255})

	if g.frames/30%2 == 0 {
		restart := "Press SPACE to play again"
		restX := (w - len(restart)*7) / 2
		text.Draw(screen, restart, menuFont, restX, 650, color.RGBA{255, 255, 100, 255})
	}

	quit := "Press ESC to quit"
	quitX := (w - len(quit)*7) / 2
	text.Draw(screen, quit, menuFont, quitX, 700, color.RGBA{180, 180, 180, 255})
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	screenWidth, screenHeight := screen.Size()
	g.calculateLayout(screenWidth, screenHeight)

	switch g.snake.GetGameMode() {
	case ModeMenu:
		g.drawMenu(screen)
	case ModePlaying:
		// Dark background
		ebitenutil.DrawRect(screen, 0, 0, float64(screenWidth), float64(screenHeight), color.RGBA{20, 24, 35, 255})
		g.drawGrid(screen)
		g.drawFood(screen)
		g.drawSnake(screen)
		g.drawParticles(screen)
		g.drawHUD(screen)
	case ModeGameOver:
		g.drawGameOver(screen)
	}
}

// Layout returns the game size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// drawParticles draws particle effects
func (g *Game) drawParticles(screen *ebiten.Image) {
	for _, p := range g.particles {
		alpha := uint8(float64(p.life) / float64(p.maxLife) * 255)
		size := float64(3) * (float64(p.life) / float64(p.maxLife))
		ebitenutil.DrawRect(screen, p.x-size/2, p.y-size/2, size, size, color.RGBA{p.r, p.g, p.b, alpha})
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
