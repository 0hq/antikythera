package main

import (
	"log"
	"github.com/notnil/chess"
	"math"
)

// define new engine
var engine_minimax_plain_ab Engine = Engine{
	name: "Minimax Plain Alpha Beta",
	features: EngineFeatures{
		plain: true,
		parallel: false,
		alphabeta: true,
		iterative_deepening: false,
		mtdf: false,
	},
	engine_func: minimax_plain_ab_engine_func,
}

func minimax_plain_ab_engine_func(pos *chess.Position, cfg EngineConfig) (best *chess.Move, eval int) {
	best, eval = minimax_plain_ab_starter(pos, cfg.ply, pos.Turn() == chess.White)
	log.Println("Plain minimax results", best, eval)
	return
}

func minimax_plain_ab_starter(position *chess.Position, ply int, max bool) (best *chess.Move, eval int) {
	moves := sort_moves_v0(position.ValidMoves(), position.Board())
	eval = -1 * math.MaxInt
	for _, move := range moves {
		score := -1 * minimax_plain_ab_searcher(position.Update(move), ply-1, !max, -1 * math.MaxInt, math.MaxInt)
		log.Println("Top Level Move:", move, "Eval:", score)
		if score > eval {
			log.Println("New best move:", move)
			eval = score
			best = move
		}
	}
	return best, eval
}

func minimax_plain_ab_searcher(position *chess.Position, ply int, max bool, alpha int, beta int) (eval int) {
	explored++
	if ply == 0 {
		return evaluate_position_v1(position.Board()) * bool_to_int(max)
	}

	moves := sort_moves_v0(position.ValidMoves(), position.Board())
    for _, move := range moves {
        score := -1 * minimax_plain_ab_searcher(position.Update(move), ply - 1, !max, -beta, -alpha)
		if score >= beta {
			return beta
		}
        if score > alpha {
            alpha = score
        }
    }

	return alpha
}
