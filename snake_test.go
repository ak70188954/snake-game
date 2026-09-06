package main

import (
	"testing"
	"time"
)

func TestNewSnake(t *testing.T) {
	snake := NewSnake(40, 30, 20)
	
	if snake.GetGridWidth() != 40 {
		t.Errorf("expected grid width 40, got %d", snake.GetGridWidth())
	}
	if snake.GetGridHeight() != 30 {
		t.Errorf("expected grid height 30, got %d", snake.GetGridHeight())
	}
	if snake.GetCellSize() != 20 {
		t.Errorf("expected cell size 20, got %d", snake.GetCellSize())
	}
	if snake.GetGameMode() != ModeMenu {
		t.Errorf("expected game mode Menu, got %v", snake.GetGameMode())
	}
}

func TestStartGame(t *testing.T) {
	snake := NewSnake(40, 30, 20)
	snake.StartGame()
	
	if snake.GetGameMode() != ModePlaying {
		t.Errorf("expected game mode Playing, got %v", snake.GetGameMode())
	}
	if snake.GetScore() != 0 {
		t.Errorf("expected score 0, got %d", snake.GetScore())
	}
	
	body := snake.GetBody()
	if len(body) != 3 {
		t.Errorf("expected snake length 3, got %d", len(body))
	}
	
	if snake.GetDirection() != Right {
		t.Errorf("expected direction Right, got %v", snake.GetDirection())
	}
}

func TestMoveRight(t *testing.T) {
	snake := NewSnake(40, 30, 20)
	snake.StartGame()
	
	initialBody := snake.GetBody()
	initialHead := initialBody[0]
	
	// Move right (same direction as initial)
	snake.SetDirection(Right)
	
	// Simulate time passage
	snake.lastMove = time.Now().Add(-300 * time.Millisecond)
	snake.Update()
	
	newBody := snake.GetBody()
	newHead := newBody[0]
	
	// Head should move right (X+1)
	if newHead.X != initialHead.X+1 {
		t.Errorf("expected head X to increase by 1, got %d (expected %d)", newHead.X, initialHead.X+1)
	}
}

func TestMoveLeft(t *testing.T) {
	snake := NewSnake(40, 30, 20)
	snake.StartGame()
	
	initialBody := snake.GetBody()
	initialHead := initialBody[0]
	
	// Move left (opposite direction - should be blocked)
	snake.SetDirection(Left)
	
	// Simulate time passage
	snake.lastMove = time.Now().Add(-300 * time.Millisecond)
	snake.Update()
	
	newBody := snake.GetBody()
	newHead := newBody[0]
	
	// Head should NOT move left (180 turn blocked)
	if newHead.X != initialHead.X {
		t.Errorf("expected head to stay at X=%d (180 turn blocked), got %d", initialHead.X, newHead.X)
	}
}

func TestMoveUp(t *testing.T) {
	snake := NewSnake(40, 30, 20)
	snake.StartGame()
	
	initialBody := snake.GetBody()
	initialHead := initialBody[0]
	
	// Move up
	snake.SetDirection(Up)
	
	// Simulate time passage
	snake.lastMove = time.Now().Add(-300 * time.Millisecond)
	snake.Update()
	
	newBody := snake.GetBody()
	newHead := newBody[0]
	
	// Head should move up (Y-1)
	if newHead.Y != initialHead.Y-1 {
		t.Errorf("expected head Y to decrease by 1, got %d (expected %d)", newHead.Y, initialHead.Y-1)
	}
}

func TestMoveDown(t *testing.T) {
	snake := NewSnake(40, 30, 20)
	snake.StartGame()
	
	initialBody := snake.GetBody()
	initialHead := initialBody[0]
	
	// Move down
	snake.SetDirection(Down)
	
	// Simulate time passage
	snake.lastMove = time.Now().Add(-300 * time.Millisecond)
	snake.Update()
	
	newBody := snake.GetBody()
	newHead := newBody[0]
	
	// Head should move down (Y+1)
	if newHead.Y != initialHead.Y+1 {
		t.Errorf("expected head Y to increase by 1, got %d (expected %d)", newHead.Y, initialHead.Y+1)
	}
}

func TestPrevent180Turn(t *testing.T) {
	snake := NewSnake(40, 30, 20)
	snake.StartGame()
	
	// Snake is moving Right, try to go Left (180 turn)
	snake.SetDirection(Left)
	
	// nextDir should be Left but direction should stay Right
	if snake.nextDir != Left {
		t.Errorf("expected nextDir to be Left, got %v", snake.nextDir)
	}
	if snake.GetDirection() != Right {
		t.Errorf("expected direction to stay Right, got %v", snake.GetDirection())
	}
}

func TestWallCollision(t *testing.T) {
	snake := NewSnake(10, 10, 20)
	snake.StartGame()
	
	// Move snake all the way right until it hits the wall
	for i := 0; i < 15; i++ {
		snake.SetDirection(Right)
		snake.lastMove = time.Now().Add(-300 * time.Millisecond)
		snake.Update()
	}
	
	if snake.GetGameMode() != ModeGameOver {
		t.Errorf("expected game over after wall collision, got %v", snake.GetGameMode())
	}
}

func TestSelfCollision(t *testing.T) {
	snake := NewSnake(10, 10, 20)
	snake.StartGame()
	
	// Make snake long enough to collide with itself
	for i := 0; i < 8; i++ {
		snake.SetDirection(Right)
		snake.lastMove = time.Now().Add(-300 * time.Millisecond)
		snake.Update()
	}
	
	// Create a loop: up, left, down
	snake.SetDirection(Up)
	snake.lastMove = time.Now().Add(-300 * time.Millisecond)
	snake.Update()
	
	snake.SetDirection(Left)
	snake.lastMove = time.Now().Add(-300 * time.Millisecond)
	snake.Update()
	
	snake.SetDirection(Down)
	snake.lastMove = time.Now().Add(-300 * time.Millisecond)
	snake.Update()
	
	// Should collide with itself
	if snake.GetGameMode() != ModeGameOver {
		t.Errorf("expected game over after self collision, got %v", snake.GetGameMode())
	}
}

func TestEatFood(t *testing.T) {
	snake := NewSnake(40, 30, 20)
	snake.StartGame()
	
	initialScore := snake.GetScore()
	initialLength := len(snake.GetBody())
	
	// Manually move food to head position to simulate eating
	head := snake.GetBody()[0]
	snake.food = Point{head.X + 1, head.Y}
	
	// Move right to eat food
	snake.SetDirection(Right)
	snake.lastMove = time.Now().Add(-300 * time.Millisecond)
	snake.Update()
	
	newScore := snake.GetScore()
	newLength := len(snake.GetBody())
	
	if newScore != initialScore+1 {
		t.Errorf("expected score to increase by 1, got %d", newScore-initialScore)
	}
	if newLength != initialLength+1 {
		t.Errorf("expected length to increase by 1, got %d", newLength-initialLength)
	}
}

func TestScoreIncreases(t *testing.T) {
	snake := NewSnake(40, 30, 20)
	snake.StartGame()
	
	initialScore := snake.GetScore()
	
	// Eat 5 food items by placing food in front of snake
	for i := 0; i < 5; i++ {
		head := snake.GetBody()[0]
		// Place food to the right of head
		snake.food = Point{head.X + 1, head.Y}
		
		// Move right to eat
		snake.SetDirection(Right)
		snake.lastMove = time.Now().Add(-300 * time.Millisecond)
		snake.Update()
		
		if snake.GetGameMode() == ModeGameOver {
			t.Fatal("game over before eating all food")
		}
	}
	
	finalScore := snake.GetScore()
	if finalScore != initialScore+5 {
		t.Errorf("expected score to be %d, got %d", initialScore+5, finalScore)
	}
}

func TestSpeedIncreases(t *testing.T) {
	snake := NewSnake(40, 30, 20)
	snake.StartGame()
	
	initialSpeed := snake.GetSpeed()
	
	// Eat food to increase speed
	head := snake.GetBody()[0]
	snake.food = Point{head.X + 1, head.Y}
	
	snake.SetDirection(Right)
	snake.lastMove = time.Now().Add(-300 * time.Millisecond)
	snake.Update()
	
	newSpeed := snake.GetSpeed()
	
	// Speed should decrease (game gets faster)
	if newSpeed >= initialSpeed {
		t.Errorf("expected speed to decrease, got %v (initial: %v)", newSpeed, initialSpeed)
	}
}

func TestChangeDirection(t *testing.T) {
	snake := NewSnake(40, 30, 20)
	
	// Test arrow keys - direction is stored in nextDir until Update() is called
	snake.ChangeDirection("up")
	if snake.nextDir != Up {
		t.Errorf("expected nextDir Up, got %v", snake.nextDir)
	}
	
	snake.ChangeDirection("down")
	if snake.nextDir != Down {
		t.Errorf("expected nextDir Down, got %v", snake.nextDir)
	}
	
	snake.ChangeDirection("left")
	if snake.nextDir != Left {
		t.Errorf("expected nextDir Left, got %v", snake.nextDir)
	}
	
	snake.ChangeDirection("right")
	if snake.nextDir != Right {
		t.Errorf("expected nextDir Right, got %v", snake.nextDir)
	}
}

func TestWASDControls(t *testing.T) {
	snake := NewSnake(40, 30, 20)
	
	snake.ChangeDirection("w")
	if snake.nextDir != Up {
		t.Errorf("expected nextDir Up with 'w', got %v", snake.nextDir)
	}
	
	snake.ChangeDirection("s")
	if snake.nextDir != Down {
		t.Errorf("expected nextDir Down with 's', got %v", snake.nextDir)
	}
	
	snake.ChangeDirection("a")
	if snake.nextDir != Left {
		t.Errorf("expected nextDir Left with 'a', got %v", snake.nextDir)
	}
	
	snake.ChangeDirection("d")
	if snake.nextDir != Right {
		t.Errorf("expected nextDir Right with 'd', got %v", snake.nextDir)
	}
}

func TestRestartGame(t *testing.T) {
	snake := NewSnake(40, 30, 20)
	snake.StartGame()
	
	// Game over
	snake.gameOver()
	if snake.GetGameMode() != ModeGameOver {
		t.Errorf("expected ModeGameOver, got %v", snake.GetGameMode())
	}
	
	// Restart
	snake.ChangeDirection(" ")
	if snake.GetGameMode() != ModePlaying {
		t.Errorf("expected ModePlaying after restart, got %v", snake.GetGameMode())
	}
	if snake.GetScore() != 0 {
		t.Errorf("expected score to reset to 0, got %d", snake.GetScore())
	}
}

func TestFoodNotOnSnake(t *testing.T) {
	snake := NewSnake(40, 30, 20)
	snake.StartGame()
	
	body := snake.GetBody()
	food := snake.GetFood()
	
	// Food should never be on the snake
	for _, segment := range body {
		if food.X == segment.X && food.Y == segment.Y {
			t.Errorf("food spawned on snake body at (%d, %d)", food.X, food.Y)
		}
	}
}

func TestGetBodyReturnsCopy(t *testing.T) {
	snake := NewSnake(40, 30, 20)
	snake.StartGame()
	
	body1 := snake.GetBody()
	body1[0] = Point{999, 999}
	
	body2 := snake.GetBody()
	if body2[0].X == 999 {
		t.Error("GetBody should return a copy, not the original slice")
	}
}
