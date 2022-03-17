package pong

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type GameState byte

type Screen interface {
	Size() (int, int)
}

func GetCenter(screen Screen) Position {
	w, h := screen.Size()
	return Position{
		X: float32(w / 2),
		Y: float32(h / 2),
	}
}

const (
	StartState GameState = iota
	PlayState
	GameOverState
)

const (
	KeyW     = "KEYW"
	KeyS     = "KEYS"
	KeyUp    = "ARROWUP"
	KeyDown  = "ARROWDOWN"
	KeySpace = "SPACE"
)
