package main

import (
	"github.com/notnil/chess"
)

// constants for minimax
var newGame = chess.NewGame().Position()
var explored int = 0
var q_explored int = 0
const ENGINE_MINIMAX_PLAIN_PLY int = 4
const ENGINE_MINIMAX_PARALLEL_PLAIN_PLY int = 4
var DO_MOVE_SORTING bool = true
var DO_Q_MOVE_SORTING bool = false
const MAX_DEPTH int = 16

var all_engines = []Engine{
	engine_minimax_plain_ab,
	engine_minimax_parallel_plain,
	engine_minimax_plain,
}
var plain_engines = []Engine{
	engine_minimax_plain_ab,
	engine_minimax_plain,
}

func reset_counters() {
	explored = 0
	q_explored = 0
}