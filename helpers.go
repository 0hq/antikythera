package main

import (
	"fmt"
	"log"
	"time"

	"github.com/0hq/chess"
)

// wraps an engine with opening book
func wrap_engine(engine Engine, time float64, game *chess.Game) Engine {
	subengine := engine
	subengine.Reset(game.Position())
	subengine.Set_Time(time)
	return NewOpeningWrapper(subengine, game)
}

// returns engine from engine and time
func new_engine(engine Engine, time float64, game *chess.Game) Engine {
	if game != nil {
		engine.Reset(game.Position())
	}
	engine.Set_Time(time)
	return engine
}


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

func Reset_Hash_Counters() {
	hash_hits = 0
	hash_reads = 0
	hash_writes = 0
	hash_collisions = 0
}


func Reset_Global_Counters() {
	// initialize depth_count to all 0s
	for i := range depth_count {
		depth_count[i] = 0
	}
	explored = 0
	q_explored = 0
}

func reset_test_counters() {
	tests_run = 0
	tests_passed = 0
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