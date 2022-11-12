package main

import (
	"math"

	"github.com/notnil/chess"
)

type t_engine_p_ab_q struct {
	EngineClass
}

// define new engine
var engine_minimax_plain_ab_q = t_engine_p_ab_q{
	EngineClass{
		name: "Minimax Plain Alpha Beta + Quiescence",
		features: EngineFeatures{
			plain: true,
			parallel: false,
			alphabeta: true,
			iterative_deepening: false,
			mtdf: false,
		},
		engine_config: EngineConfig{
			ply: 3,
		},
		time_up: false,
	},
	// engine_func: minimax_plain_ab_q_engine_func,
}


func (e *t_engine_p_ab_q) Run_Engine(pos *chess.Position) (best *chess.Move, eval int) {
	reset_counters()
	out("Running minimax_plain_ab_q_engine_func")
	best, eval = e.minimax_plain_ab_q_starter(pos, e.engine_config.ply, pos.Turn() == chess.White)
	out("Plain minimax results", best, eval)
	out("Quiescence search explored", q_explored, "nodes")
	return
}

func (e *t_engine_p_ab_q) minimax_plain_ab_q_starter(position *chess.Position, ply int, max bool) (best *chess.Move, eval int) {
	moves := sort_moves_v1(position.ValidMoves(), position.Board())
	eval = -1 * math.MaxInt
	for _, move := range moves {
		score := -1 * e.minimax_plain_ab_q_searcher(position.Update(move), ply-1, !max, -1 * math.MaxInt, math.MaxInt)
		if PRINT_TOP_MOVES {
			out("Top Level Move:", move, "Eval:", score)
		}
		if score > eval {
			out("New best move:", move)
			eval = score
			best = move
		}
	}
	return best, eval
}

func (e *t_engine_p_ab_q) minimax_plain_ab_q_searcher(position *chess.Position, ply int, max bool, alpha int, beta int) (eval int) {
	explored++
	if ply == 0 {
		return e.quiescence_minimax_plain_ab_q(position, 0, max, alpha, beta)
	}

	moves := sort_moves_v1(position.ValidMoves(), position.Board())
	if len(moves) == 0 {
		return evaluate_position_v2(position, e.engine_config.ply, ply, bool_to_int(max))
	}
    for _, move := range moves {
        score := -1 * e.minimax_plain_ab_q_searcher(position.Update(move), ply - 1, !max, -beta, -alpha)
		if score >= beta {
			return beta
		}
        if score > alpha {
            alpha = score
        }
    }

	return alpha
}

func (e *t_engine_p_ab_q) quiescence_minimax_plain_ab_q(position *chess.Position, plycount int, max bool, alpha int, beta int) (eval int) {
	explored++
	q_explored++

	stand_pat := evaluate_position_v2(position, e.engine_config.ply, -plycount, bool_to_int(max))
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
        score := -1 * e.quiescence_minimax_plain_ab_q(position.Update(move), plycount + 1, !max, -beta, -alpha)
		if score >= beta {
			return beta
		}
        if score > alpha {
            alpha = score
        }
    }

	return alpha
}