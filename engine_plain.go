package main

import (
	"math"

	"github.com/0hq/chess"
)


type t_engine_p struct {
	EngineClass
}

// define new engine
var engine_minimax_plain = t_engine_p_ab_q{
	EngineClass{
		name: "Minimax Plain",
		features: EngineFeatures{
			plain: true,
			parallel: false,
			alphabeta: false,
			iterative_deepening: false,
			mtdf: false,
		},
	},
}

func (e *t_engine_p) Run_Engine(pos *chess.Position) (best *chess.Move, eval int) {
	reset_counters()
	best, eval = e.minimax_plain_starter(pos, e.engine_config.ply, pos.Turn() == chess.White)
	out("Plain minimax results", best, eval)
	return
}

func (e *t_engine_p) minimax_plain_starter(position *chess.Position, ply int, max bool) (best *chess.Move, eval int) {
	moves := position.ValidMoves()
	eval = -1 * math.MaxInt
	for _, move := range moves {
		tempeval := -1 * e.minimax_plain_searcher(position.Update(move), ply-1, !max)
		out("Top Level Move:", move, "Eval:", tempeval)
		if tempeval > eval {
			out("New best move:", move)
			eval = tempeval
			best = move
		}
	}
	return best, eval
}

func (e *t_engine_p) minimax_plain_searcher(position *chess.Position, ply int, max bool) (eval int) {
	explored++
	if ply == 0 {
		return evaluate_position_v1(position) * bool_to_int(max)
	}

	moves := position.ValidMoves()
	// if (len(moves) == 0) {
	// 	return evaluate_position_v1(position) * bool_to_int(max)
	// }
    eval = -1 * math.MaxInt
    for _, move := range moves {
        tempeval := -1 * e.minimax_plain_searcher(position.Update(move), ply - 1, !max)
        if tempeval > eval {
            eval = tempeval
        }
    }

	return eval
}
