package main

import (
	"time"

	"github.com/notnil/chess"
)

// function that turns boolean into either 1 or -1
func bool_to_int(b bool) int {
	if b {
		return 1
	} else {
		return -1
	}
}

func game_from_fen(pos string) *chess.Game {
	fen, err := chess.FEN(pos)
	if err != nil {
		panic(err)
	}
	return chess.NewGame(fen)
}

func duration_from_sec(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}