package main

import (
	"log"
	"github.com/notnil/chess"
	"math"
)

// define new engine
var engine_minimax_parallel_plain Engine = Engine{
	name: "Minimax Parallel Plain",
	features: EngineFeatures{
		plain: true,
		parallel: true,
		alphabeta: false,
		iterative_deepening: false,
		mtdf: false,
	},
	engine_func: minimax_parallel_plain_engine_func,
}

func minimax_parallel_plain_engine_func(pos *chess.Position) (best *chess.Move, eval int) {
	best, eval = minimax_pll_starter(pos, 4, true)
	log.Println("Parellel results", best, eval)
	return 
}

// starter is not a concurrent function, handles top level moves
// midstep is a concurrent function, but is passed the top level move and a move channel to send back instantly once done, so we know which process it is
// searcher is a concurrent function, handles the actual search, just returns a value

func minimax_pll_starter(position *chess.Position, ply int, max bool) (best *chess.Move, eval int) {
	// generate moves
	var moves []*chess.Move = position.ValidMoves()
	var length int = len(moves)

	// create channel to pass back move and eval
	eval_channel_local := make(chan int, length)
	move_channel_local := make(chan *chess.Move, length)

	// create goroutines for each move
    for _, move := range moves {
        go minimax_pll_midstep(position.Update(move), ply-1, !max, eval_channel_local, move_channel_local, move)
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

func minimax_pll_midstep(position *chess.Position, ply int, max bool, eval_channel chan int, move_channel chan *chess.Move, last_move *chess.Move) {
	// generate moves
	var moves []*chess.Move = position.ValidMoves()
	var length int = len(moves)

	// create channel to pass back move and eval
	eval_channel_local := make(chan int, length)

	// create goroutines for each move
    for _, move := range moves {
        go minimax_pll_plain_searcher(position.Update(move), ply-1, !max, eval_channel_local)
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

func minimax_pll_plain_searcher(position *chess.Position, ply int, max bool, eval_channel chan int) {
	explored++

	// max ply reached
	if ply == 0 {
		// evaluate position and send back to parent
		eval_channel <- evaluate_position_v1(position.Board(), max)
		return
	}

	// generate moves
	var moves []*chess.Move = position.ValidMoves()
	var length int = len(moves)

	// create channel to pass back move and eval
	eval_channel_local := make(chan int, length)

	// create goroutines for each move
    for _, move := range moves {
        go minimax_pll_plain_searcher(position.Update(move), ply-1, !max, eval_channel_local)
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

func minimax_pll(position *chess.Position, ply int, max bool, last_move *chess.Move, move_channel chan *chess.Move, eval_channel chan int, isRoot bool) {
	explored++

	// max ply reached
	if ply == 0 {
		// evaluate position and send back to parent
		move_channel <- last_move
		eval_channel <- evaluate_position_v1(position.Board(), max)
		return
	}

	// generate moves
	var moves []*chess.Move = position.ValidMoves()
	if (isRoot) {
		log.Println("Moves:", moves)
	}
	var length int = len(moves)

	// create channel to pass back move and eval
	move_channel_local := make(chan *chess.Move, length)
	eval_channel_local := make(chan int, length)

	// create goroutines for each move
    for _, move := range moves {
        go minimax_pll(position.Update(move), ply-1, !max, move, move_channel_local, eval_channel_local, false)
    }

	// wait for all goroutines to finish
	var eval int = -1 * math.MaxInt
	var best *chess.Move = nil
	for i := 0; i < length; i++ {
		move := <-move_channel_local
		tempeval := -1 * <-eval_channel_local
		if tempeval > eval {
			eval = tempeval
			best = move
		}
	}

	// pass value back to parent goroutine
	move_channel <- best
	eval_channel <- eval

	return
}