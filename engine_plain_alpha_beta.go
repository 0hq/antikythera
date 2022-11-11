package main

import (
	"log"
	"math"

	"github.com/notnil/chess"
)

type t_engine_p_ab struct {
	EngineClass
}

// define new engine
var engine_minimax_plain_ab = t_engine_p_ab{
	EngineClass{
		name: "Minimax Plain Alpha Beta",
		features: EngineFeatures{
			plain: true,
			parallel: false,
			alphabeta: true,
			iterative_deepening: false,
			mtdf: false,
		},
	},
}

func (e *t_engine_p_ab) Run_Engine(pos *chess.Position) (best *chess.Move, eval int) {
	reset_counters()
	best, eval = e.minimax_plain_ab_starter(pos, e.engine_config.ply, pos.Turn() == chess.White)
	log.Println("Plain minimax results", best, eval)
	return
}

func (e *t_engine_p_ab) minimax_plain_ab_starter(position *chess.Position, ply int, max bool) (best *chess.Move, eval int) {
	moves := sort_moves_v1(position.ValidMoves(), position.Board())
	eval = -1 * math.MaxInt
	for _, move := range moves {
		score := -1 * e.minimax_plain_ab_searcher(position.Update(move), ply-1, !max, -1 * math.MaxInt, math.MaxInt)
		log.Println("Top Level Move:", move, "Eval:", score)
		if score > eval {
			log.Println("New best move:", move)
			eval = score
			best = move
		}
	}
	return best, eval
}

func (e *t_engine_p_ab) minimax_plain_ab_searcher(position *chess.Position, ply int, max bool, alpha int, beta int) (eval int) {
	explored++
	if ply == 0 {
		return evaluate_position_v1(position) * bool_to_int(max)
	}

	moves := sort_moves_v1(position.ValidMoves(), position.Board())
    for _, move := range moves {
        score := -1 * e.minimax_plain_ab_searcher(position.Update(move), ply - 1, !max, -beta, -alpha)
		if score >= beta {
			return beta
		}
        if score > alpha {
            alpha = score
        }
    }

	return alpha
}
