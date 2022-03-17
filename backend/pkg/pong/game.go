package pong

type Game struct {
	State    GameState
	Ball     *Ball
	Player1  *Paddle
	Player2  *Paddle
	Rally    int
	Level    int
	MaxScore int
	Post     chan bool
}
type GameSend struct {
	State    GameState
	Ball     *Ball
	Player1  *Paddle
	Player2  *Paddle
	Rally    int
	Level    int
	MaxScore int
}
type GameWindow struct {
	Width  int
	Height int
}

const (
	initBallVelocity = 5.0
	initPaddleSpeed  = 10.0
	speedUpdateCount = 6
	speedIncrement   = 0.5
)

const (
	windowWidth  = 400
	windowHeight = 400
)

func (s *GameWindow) Size() (int, int) {
	return s.Width, s.Height
}

var screen = GameWindow{
	Width:  windowWidth,
	Height: windowHeight,
}

func (g *Game) Init() {
	g.State = StartState
	g.MaxScore = 11

	g.Player1 = &Paddle{
		Position: Position{
			X: InitPaddleShift,
			Y: float32(windowHeight / 2)},
		Score:  0,
		Speed:  initPaddleSpeed,
		Width:  InitPaddleWidth,
		Height: InitPaddleHeight,
		Up:     KeyUp,
		Down:   KeyDown,
	}
	g.Player2 = &Paddle{
		Position: Position{
			X: windowWidth - InitPaddleShift - InitPaddleWidth,
			Y: float32(windowHeight / 2)},
		Score:  0,
		Speed:  initPaddleSpeed,
		Width:  InitPaddleWidth,
		Height: InitPaddleHeight,
		Up:     KeyW,
		Down:   KeyS,
	}
	g.Ball = &Ball{
		Position: Position{
			X: float32(windowWidth / 2),
			Y: float32(windowHeight / 2)},
		Radius:    InitBallRadius,
		XVelocity: initBallVelocity,
		YVelocity: initBallVelocity,
	}
	g.Level = 0
}

func (g *Game) ResponseToKeyPress(key string) {
	switch key {
	case KeyUp:
		g.Player1.Update(&screen, "UP")
	case KeyDown:
		g.Player1.Update(&screen, "DOWN")

	case KeyW:
		g.Player2.Update(&screen, "UP")

	case KeyS:
		g.Player2.Update(&screen, "DOWN")

	}
}

func (g *Game) Update(key string) error {
	switch g.State {
	case StartState:
		if key == KeySpace {
			g.State = PlayState
		}
	case PlayState:
		g.ResponseToKeyPress(key)

		if g.Player1.Score == g.MaxScore || g.Player2.Score == g.MaxScore {
			g.State = GameOverState
		}

	case GameOverState:
		if key == KeySpace {
			g.reset(&screen, StartState)
		}
	}
	// g.Draw()
	return nil
}

func (g *Game) Draw() {
	w, _ := screen.Size()

	xV := g.Ball.XVelocity
	// calculate the position of the ball and
	g.Ball.Update(g.Player1, g.Player2, &screen)
	if xV*g.Ball.XVelocity < 0 {
		// score up when ball touches human player's paddle
		if g.Ball.X < float32(w/2) {
			g.Player1.Score++
		}

		g.Rally++

		// spice things up
		if (g.Rally)%speedUpdateCount == 0 {
			g.Level++
			g.Ball.XVelocity += speedIncrement
			g.Ball.YVelocity += speedIncrement
			g.Player1.Speed += speedIncrement
			g.Player2.Speed += speedIncrement
		}
	}
	if g.Ball.X < 0 {
		g.Player2.Score++
		g.reset(&screen, StartState)
	} else if g.Ball.X > float32(w) {
		g.Player1.Score++
		g.reset(&screen, StartState)
	}

	// score depend on ball position
	g.Post <- true
	// send state
}

func (g *Game) reset(screen Screen, state GameState) {
	w, _ := screen.Size()
	g.State = state
	g.Rally = 0
	g.Level = 0
	if state == StartState {
		g.Player1.Score = 0
		g.Player2.Score = 0
	}
	g.Player1.Position = Position{
		X: InitPaddleShift, Y: GetCenter(screen).Y}
	g.Player2.Position = Position{
		X: float32(w - InitPaddleShift - InitPaddleWidth), Y: GetCenter(screen).Y}
	g.Ball.Position = GetCenter(screen)
	g.Ball.XVelocity = initBallVelocity
	g.Ball.YVelocity = initBallVelocity
}
