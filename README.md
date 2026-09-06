# Snake Game

Классическая игра "Змейка", написанная на Go с графикой через [Ebitengine](https://github.com/hajimehoshi/ebiten).

## Особенности

- 🐍 Адаптивное игровое поле (подстраивается под размер окна)
- 🎨 Графика с градиентами, глазами змейки и анимированным меню
- 🎮 Управление стрелками или WASD
- 📈 Увеличение скорости с каждой съеденной едой
- 🏆 Система рейтинга (Beginner → Legend)
- ✅ 20 unit-тестов

## Установка

### Требования

- Go 1.25+

### Сборка

```bash
git clone https://github.com/ak70188954/snake-game.git
cd snake-game
go build -o snake-game .
./snake-game
```

### Запуск

```bash
go run .
```

## Управление

| Клавиша        | Действие           |
|----------------|--------------------|
| ↑ / W          | Движение вверх     |
| ↓ / S          | Движение вниз      |
| ← / A          | Движение влево     |
| → / D          | Движение вправо    |
| Space / Enter  | Начать игру        |
| Escape         | Выйти              |

## Структура проекта

```
snake-game/
├── main.go          # Графика, рендеринг, игровой цикл
├── snake.go         # Логика игры (движение, столкновения, еда)
├── snake_test.go    # Unit-тесты (20 тестов)
├── go.mod           # Зависимости
└── .gitignore
```

## Архитектура

### Логика (snake.go)

- `Snake` — основная структура с телом, направлением, скоростью
- `Direction` — тип направления (Up, Down, Left, Right)
- `GameMode` — состояние игры (Menu, Playing, GameOver)
- `SetDirection()` — блокирует разворот на 180°
- `Update()` — игровой цикл логики (движение, коллизии)

### Графика (main.go)

- `Game` — структура с рендерингом
- `calculateLayout()` — адаптивный размер клетки
- `drawGrid()` — шахматное поле с рамкой
- `drawSnake()` — змейка с градиентом и глазами
- `drawFood()` — еда с эффектом блеска
- `drawMenu()` / `drawGameOver()` — экраны меню

### Тесты (snake_test.go)

| Тест                  | Что проверяет                          |
|-----------------------|----------------------------------------|
| TestNewSnake          | Инициализация змейки                   |
| TestStartGame         | Старт игры, начальная позиция          |
| TestMoveRight/Left/Up/Down | Движение в каждом направлении   |
| TestPrevent180Turn    | Блокировка разворота на 180°          |
| TestWallCollision     | Столкновение со стеной                 |
| TestSelfCollision     | Столкновение с собой                   |
| TestEatFood           | Поедание еды, увеличение счёта         |
| TestScoreIncreases    | Рост счёта при поедании                |
| TestSpeedIncreases    | Увеличение скорости                    |
| TestChangeDirection   | Все направления клавиш                 |
| TestWASDControls      | WASD управление                        |
| TestTurnFromUp/Down/Left | Повороты из каждого направления    |
| TestRestartGame       | Перезапуск игры                        |
| TestFoodNotOnSnake    | Еда не появляется на змейке            |
| TestGetBodyReturnsCopy | GetBody возвращает копию, не оригинал |

## Запуск тестов

```bash
go test -v ./...
```

## Лицензия

MIT
