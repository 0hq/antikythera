package main

import (
	"math"

	"github.com/0hq/chess"
	"github.com/0hq/chess/opening"
)

/*

Counters.

*/

var tests_run int = 0
var tests_passed int = 0
var explored int = 0
var q_explored int = 0
var depth_count []int = make([]int, MAX_DEPTH)
var hash_writes int = 0
var hash_reads int = 0
var hash_hits int = 0
var hash_collisions int = 0


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
const MVV_LVA_OFFSET int = math.MaxInt - 256
const KILLER_VALUE int = 10;
const TT_MOVE_BONUS int = 100;
const TT_MOVE_VALUE int = MVV_LVA_OFFSET + TT_MOVE_BONUS;
const KILLER_ONE_VALUE int = MVV_LVA_OFFSET - KILLER_VALUE 
const KILLER_TWO_VALUE int = MVV_LVA_OFFSET - KILLER_VALUE * 2
const NORMAL_MOVE_VALUE int = MVV_LVA_OFFSET - KILLER_VALUE * 2 - 100

/*

Configuration.

*/

var production_mode = false
const DO_MOVE_SORTING bool = true
const DO_Q_MOVE_SORTING bool = true
const MAX_DEPTH int = 30
const DO_Q_MOVE_PRUNING bool = true 
const DO_DEPTH_COUNT bool = true // disable for performance

/*

Logging options

*/

const PRINT_TOP_MOVES bool = false 
const QUIET_MODE bool = false

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

