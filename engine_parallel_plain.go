package main

import (
	"math"

	"github.com/0hq/chess"
)

type t_engine_p_pll struct {
	EngineClass
}

// define new engine
var engine_minimax_parallel_plain = t_engine_p_pll{
	EngineClass{
		name: "Minimax Parallel Plain",
		features: EngineFeatures{
			plain: true,
			parallel: true,
			alphabeta: false,
			iterative_deepening: false,
			mtdf: false,
		},
	},
}

func (e *t_engine_p_pll) Run_Engine(pos *chess.Position) (best *chess.Move, eval int) {
	reset_counters()
	best, eval = e.minimax_parallel_plain_starter(pos, e.engine_config.ply, true)
	out("Parellel results", best, eval)
	return 
}

// starter is not a concurrent function, handles top level moves
// midstep is a concurrent function, but is passed the top level move and a move channel to send back instantly once done, so we know which process it is
// searcher is a concurrent function, handles the actual search, just returns a value

func (e *t_engine_p_pll) minimax_parallel_plain_starter(position *chess.Position, ply int, max bool) (best *chess.Move, eval int) {
	// generate moves
	var moves []*chess.Move = position.ValidMoves()
	var length int = len(moves)

	// create channel to pass back move and eval
	eval_channel_local := make(chan int, length)
	move_channel_local := make(chan *chess.Move, length)

	// create goroutines for each move
    for _, move := range moves {
        go e.minimax_parallel_plain_midstep(position.Update(move), ply-1, !max, eval_channel_local, move_channel_local, move)
    }

	// wait for all goroutines to finish
	eval = -1 * math.MaxInt
	best = nil
	for i := 0; i < length; i++ {
		tempeval := -1 * <-eval_channel_local
		tempmove := <-move_channel_local
		if tempeval > eval {
			eval = tempeval
			best = tempmove
		}
	}

	return best, eval
}

func (e *t_engine_p_pll) minimax_parallel_plain_midstep(position *chess.Position, ply int, max bool, eval_channel chan int, move_channel chan *chess.Move, last_move *chess.Move) {
	// generate moves
	var moves []*chess.Move = position.ValidMoves()
	var length int = len(moves)

	// create channel to pass back move and eval
	eval_channel_local := make(chan int, length)

	// create goroutines for each move
    for _, move := range moves {
        go e.minimax_parallel_plain_searcher(position.Update(move), ply-1, !max, eval_channel_local)
    }

	// wait for all goroutines to finish
	var eval int = -1 * math.MaxInt
	for i := 0; i < length; i++ {
		tempeval := -1 * <-eval_channel_local
		if tempeval > eval {
			eval = tempeval
		}
	}

	// pass value back to parent goroutine
	eval_channel <- eval
	move_channel <- last_move

	return
}

func (e *t_engine_p_pll) minimax_parallel_plain_searcher(position *chess.Position, ply int, max bool, eval_channel chan int) {
	explored++

	// max ply reached
	if ply == 0 {
		// evaluate position and send back to parent
		// multiply by -1 to flip sign if we are minimizing,
		eval_channel <- evaluate_position_v1(position) * bool_to_int(max)
		return
	}

	// generate moves
	var moves []*chess.Move = position.ValidMoves()
	var length int = len(moves)

	// create channel to pass back move and eval
	eval_channel_local := make(chan int, length)

	// create goroutines for each move
    for _, move := range moves {
        go e.minimax_parallel_plain_searcher(position.Update(move), ply-1, !max, eval_channel_local)
    }

	// wait for all goroutines to finish
	var eval int = -1 * math.MaxInt
	for i := 0; i < length; i++ {
		tempeval := -1 * <-eval_channel_local
		if tempeval > eval {
			eval = tempeval
		}
	}

	// pass value back to parent goroutine
	eval_channel <- eval

	return
}

