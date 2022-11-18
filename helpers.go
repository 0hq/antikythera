package main

import (
	"fmt"
	"log"
	"time"

	"github.com/0hq/chess"
)

// Custom printLn function
// logs and prints any input
func out(a ...any) {
	if !production_mode {
		log.Println(a...)
	}
	fmt.Println(a...)
}

// function that turns boolean into either 1 or -1
func bool_to_int(b bool) int {
	if b {
		return 1
	} else {
		return -1
	}
}

func reset_counters() {
	// initialize depth_count to all 0s
	// for i := range depth_count {
	// 	depth_count[i] = 0
	// }
	// tests_run = 0
	// tests_passed = 0
	explored = 0
	q_explored = 0
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