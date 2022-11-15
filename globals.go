package main

import (
	"github.com/notnil/chess"
	"github.com/notnil/chess/opening"
)

/*

Counters.

*/

var tests_run int = 0
var tests_passed int = 0
var explored int = 0
var q_explored int = 0
// var depth_count []int = make([]int, 100)


/*

Game constants.

*/

var global_UCINotation chess.UCINotation 
var global_AlgebraicNotation chess.AlgebraicNotation
var global_Opening_Book = opening.NewBookECO()
const CHESS_START_POSITION = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
var newGame = chess.NewGame()

/*

Engine constants.

*/

const CHECKMATE_VALUE int = 30000

/*

Configuration.

*/

var production_mode = false
const DO_MOVE_SORTING bool = true
const DO_Q_MOVE_SORTING bool = true
const MAX_DEPTH int = 16
const DO_Q_MOVE_PRUNING bool = true
const PRINT_TOP_MOVES bool = false
// const DO_Q_MOVE_CHECKS bool = false
// const DO_Q_MOVE_PROMOS bool = false // disabled manually

// var all_engines = []Engine{
// 	&engine_minimax_plain_ab,
// 	&engine_minimax_parallel_plain,
// 	&engine_minimax_plain,
// }
// var plain_engines = []Engine{
// 	&engine_minimax_plain_ab,
// 	&engine_minimax_plain,
// }

