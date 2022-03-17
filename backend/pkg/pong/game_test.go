package pong

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	var game_t = Game{}

	game_t.Init()
	if game_t.Ball.X != windowWidth/2 {
		t.Error("Expected x to be 200")
	}

}

func TestUpdate(t *testing.T) {
	var game_t = Game{}
	game_t.Init()
	game_t.State = PlayState
	game_t.Update("W")
	require.Equal(t, float32(190), game_t.Player2.Position.Y)
	game_t.Update("S")
	require.NotEqual(t, float32(190), game_t.Player2.Position.Y)
	game_t.Update("S")

	require.Equal(t, float32(210), game_t.Player2.Position.Y)
}
