package main

import (
	"math/rand"
	"time"
)

// Direction represents the direction the snake is moving
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

// GameMode represents the current game state
type GameMode int

const (
	ModeMenu GameMode = iota
	ModePlaying
	ModeGameOver
)

// Point represents a coordinate on the grid
type Point struct {
	X, Y int
}

// Snake represents the snake game logic
type Snake struct {
	body       []Point
	direction  Direction
	nextDir    Direction
	food       Point
	score      int
	speed      time.Duration
	lastMove   time.Time
	gameMode   GameMode
	gridWidth  int
	gridHeight int
	cellSize   int
}

// NewSnake creates a new snake game instance
func NewSnake(gridWidth, gridHeight, cellSize int) *Snake {
	return &Snake{
		direction:  Right,
		nextDir:    Right,
		score:      0,
		speed:      200 * time.Millisecond,
		gridWidth:  gridWidth,
		gridHeight: gridHeight,
		cellSize:   cellSize,
		gameMode:   ModeMenu,
	}
}

// StartGame initializes a new game
func (s *Snake) StartGame() {
	// Start snake in the middle area, moving right with tail behind (to the left)
	centerX := s.gridWidth / 2
	centerY := s.gridHeight / 2

	s.body = []Point{
		{X: centerX, Y: centerY},
		{X: centerX - 1, Y: centerY},
		{X: centerX - 2, Y: centerY},
	}
	s.direction = Right
	s.nextDir = Right
	s.score = 0
	s.speed = 200 * time.Millisecond
	s.spawnFood()
	s.gameMode = ModePlaying
	s.lastMove = time.Now()
}

// Update handles game logic updates
func (s *Snake) Update() {
	if s.gameMode != ModePlaying {
		return
	}

	// Check if it's time to move
	if time.Since(s.lastMove) < s.speed {
		return
	}

	s.lastMove = time.Now()

	// Update direction
	s.direction = s.nextDir

	// Calculate new head position
	head := s.body[0]
	var newHead Point

	switch s.direction {
	case Up:
		newHead = Point{head.X, head.Y - 1}
	case Down:
		newHead = Point{head.X, head.Y + 1}
	case Left:
		newHead = Point{head.X - 1, head.Y}
	case Right:
		newHead = Point{head.X + 1, head.Y}
	}

	// Check wall collision
	if newHead.X < 0 || newHead.X >= s.gridWidth || newHead.Y < 0 || newHead.Y >= s.gridHeight {
		s.gameOver()
		return
	}

	// Check self collision
	for _, segment := range s.body {
		if newHead.X == segment.X && newHead.Y == segment.Y {
			s.gameOver()
			return
		}
	}

	// Add new head
	s.body = append([]Point{newHead}, s.body...)

	// Check food collision
	if newHead.X == s.food.X && newHead.Y == s.food.Y {
		s.score++
		// Increase speed slightly
		if s.speed > 80*time.Millisecond {
			s.speed -= 5 * time.Millisecond
		}
		s.spawnFood()
	} else {
		// Remove tail if no food eaten
		s.body = s.body[:len(s.body)-1]
	}
}

// SetDirection sets the snake's direction (with validation)
func (s *Snake) SetDirection(dir Direction) {
	// Prevent 180-degree turns
	if s.direction+2 == dir {
		return
	}
	s.nextDir = dir
}

// GetBody returns a copy of the snake's body
func (s *Snake) GetBody() []Point {
	body := make([]Point, len(s.body))
	copy(body, s.body)
	return body
}

// GetFood returns the food position
func (s *Snake) GetFood() Point {
	return s.food
}

// GetScore returns the current score
func (s *Snake) GetScore() int {
	return s.score
}

// GetDirection returns the current direction
func (s *Snake) GetDirection() Direction {
	return s.direction
}

// GetGameMode returns the current game mode
func (s *Snake) GetGameMode() GameMode {
	return s.gameMode
}

// GetCellSize returns the cell size
func (s *Snake) GetCellSize() int {
	return s.cellSize
}

// GetGridWidth returns the grid width
func (s *Snake) GetGridWidth() int {
	return s.gridWidth
}

// GetGridHeight returns the grid height
func (s *Snake) GetGridHeight() int {
	return s.gridHeight
}

// GetSpeed returns the current speed
func (s *Snake) GetSpeed() time.Duration {
	return s.speed
}

// spawnFood places food at a random position not occupied by the snake
func (s *Snake) spawnFood() {
	rand.Seed(time.Now().UnixNano())

	// Generate random position until it's not on the snake
	for {
		x := rand.Intn(s.gridWidth)
		y := rand.Intn(s.gridHeight)

		onSnake := false
		for _, segment := range s.body {
			if x == segment.X && y == segment.Y {
				onSnake = true
				break
			}
		}

		if !onSnake {
			s.food = Point{x, y}
			return
		}
	}
}

// gameOver handles game over state
func (s *Snake) gameOver() {
	s.gameMode = ModeGameOver
}

// ChangeDirection changes the direction based on key input
func (s *Snake) ChangeDirection(key string) {
	switch key {
	case "up", "w":
		s.SetDirection(Up)
	case "down", "s":
		s.SetDirection(Down)
	case "left", "a":
		s.SetDirection(Left)
	case "right", "d":
		s.SetDirection(Right)
	case " ":
		if s.gameMode == ModeMenu || s.gameMode == ModeGameOver {
			s.StartGame()
		}
	case "enter":
		if s.gameMode == ModeMenu || s.gameMode == ModeGameOver {
			s.StartGame()
		}
	}
}
