package pong

type Key string

type Paddle struct {
	Position
	Score        int
	Speed        float32
	Width        int
	Height       int
	Up           Key
	Down         Key
	pressed      keysPressed
	scorePrinted scorePrinted
}

type keysPressed struct {
	up   bool
	down bool
}

type scorePrinted struct {
	score   int
	printed bool
	x       int
	y       int
}

const (
	InitPaddleWidth  = 20
	InitPaddleHeight = 100
	InitPaddleShift  = 50
)

func (p *Paddle) Update(screen Screen, key string) {
	_, h := screen.Size()

	// if inpututil.IsKeyJustPressed(p.Up) {
	// 	p.pressed.down = false
	// 	p.pressed.up = true
	// } else if inpututil.IsKeyJustReleased(p.Up) || !ebiten.IsKeyPressed(p.Up) {
	// 	p.pressed.up = false
	// }
	// if inpututil.IsKeyJustPressed(p.Down) {
	// 	p.pressed.up = false
	// 	p.pressed.down = true
	// } else if inpututil.IsKeyJustReleased(p.Down) || !ebiten.IsKeyPressed(p.Down) {
	// 	p.pressed.down = false
	// }

	if key == "UP" {
		p.Y -= p.Speed
	} else if key == "DOWN" {
		p.Y += p.Speed
	}

	if p.Y-float32(p.Height/2) < 0 {
		p.Y = float32(1 + p.Height/2)
	} else if p.Y+float32(p.Height/2) > float32(h) {
		p.Y = float32(h - p.Height/2 - 1)
	}
}

func (p *Paddle) AiUpdate(b *Ball) {
	// unbeatable haha
	p.Y = b.Y
}

func (p *Paddle) Draw(screen Screen) {
	if p.scorePrinted.score != p.Score && p.scorePrinted.printed {
		p.scorePrinted.printed = false
	}
	if p.scorePrinted.score == 0 && !p.scorePrinted.printed {
		p.scorePrinted.x = int(p.X + (GetCenter(screen).X-p.X)/2)
		p.scorePrinted.y = int(2 * 30)
	}
	if (p.scorePrinted.score == 0 || p.scorePrinted.score != p.Score) && !p.scorePrinted.printed {
		p.scorePrinted.score = p.Score
		p.scorePrinted.printed = true
	}
}
