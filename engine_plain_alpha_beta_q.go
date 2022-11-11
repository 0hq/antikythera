package main

import (
	"log"
	"math"

	"github.com/notnil/chess"
)

// define new engine
var engine_minimax_plain_ab_q Engine = Engine{
	name: "Minimax Plain Alpha Beta + Quiescence",
	features: EngineFeatures{
		plain: true,
		parallel: false,
		alphabeta: true,
		iterative_deepening: false,
		mtdf: false,
	},
	engine_func: minimax_plain_ab_q_engine_func,
}

func minimax_plain_ab_q_engine_func(pos *chess.Position, cfg EngineConfig) (best *chess.Move, eval int) {
	best, eval = minimax_plain_ab_q_starter(pos, cfg.ply, pos.Turn() == chess.White)
	log.Println("Plain minimax results", best, eval)
	log.Println("Quiescence search explored", q_explored, "nodes")
	return
}

func minimax_plain_ab_q_starter(position *chess.Position, ply int, max bool) (best *chess.Move, eval int) {
	moves := sort_moves_v1(position.ValidMoves(), position.Board())
	eval = -1 * math.MaxInt
	for _, move := range moves {
		score := -1 * minimax_plain_ab_q_searcher(position.Update(move), ply-1, !max, -1 * math.MaxInt, math.MaxInt)
		log.Println("Top Level Move:", move, "Eval:", score)
		if score > eval {
			log.Println("New best move:", move)
			eval = score
			best = move
		}
	}
	return best, eval
}

func minimax_plain_ab_q_searcher(position *chess.Position, ply int, max bool, alpha int, beta int) (eval int) {
	explored++
	if ply == 0 {
		return quiescence_minimax_plain_ab_q(position, 0, max, alpha, beta)
	}

	moves := sort_moves_v1(position.ValidMoves(), position.Board())
	if len(moves) == 0 {
		return evaluate_position_v1(position) * bool_to_int(max)
	}
    for _, move := range moves {
        score := -1 * minimax_plain_ab_q_searcher(position.Update(move), ply - 1, !max, -beta, -alpha)
		if score >= beta {
			return beta
		}
        if score > alpha {
            alpha = score
        }
    }

	return alpha
}

func quiescence_minimax_plain_ab_q(position *chess.Position, plycount int, max bool, alpha int, beta int) (eval int) {
	explored++
	q_explored++

	stand_pat := evaluate_position_v1(position) * bool_to_int(max)
	if stand_pat >= beta {
        return beta;
	}
    if alpha < stand_pat {
        alpha = stand_pat;
	}
	
	moves := quiescence_moves_v1(position.ValidMoves(), position.Board())

	if plycount > MAX_DEPTH || len(moves) == 0 {
		return stand_pat 
	}

    for _, move := range moves {
        score := -1 * quiescence_minimax_plain_ab_q(position.Update(move), plycount + 1, !max, -beta, -alpha)
		if score >= beta {
			return beta
		}
        if score > alpha {
            alpha = score
        }
    }

	return alpha
}