package main

import (
	"math"
	"time"

	"github.com/notnil/chess"
)

type t_engine_p_ab_q_id struct {
	EngineClass
}

// define new engine
var engine_minimax_id_ab_q = t_engine_p_ab_q_id{
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
			ply: 0,
		},
		time_up: false,
	},
	// engine_func: minimax_id_ab_q_engine_func,
} 

func (e *t_engine_p_ab_q_id) Run_Engine(pos *chess.Position) (best *chess.Move, eval int) {
	reset_counters()
	out("Running", e.name)
	e.time_up = false
	e.start_time = time.Now()
	depth := 1
	for {
		out()
		out("Iterative deepening depth", depth)
		t_best, t_eval := e.minimax_id_ab_q_starter(pos, depth, pos.Turn() == chess.White)
		if e.Check_Time_Up() {
			out("Time up, returning best move so far.")
			break
		} else {
			best = t_best
			eval = t_eval
		}
		out("Best move:", best, "Eval:", eval)
		out("Depth:", depth, "Nodes:", explored)
		
		out("Time since start_time:", time.Since(e.start_time))
		if eval >= 30000 { // break on checkmate win
			break
		}
		depth++
	}
	out()
	out("Engine results", best, eval)
	out("Quiescence search explored", q_explored, "nodes")
	return
}

func (e *t_engine_p_ab_q_id) minimax_id_ab_q_starter(position *chess.Position, ply int, max bool) (best *chess.Move, eval int) {
	moves := sort_moves_v1(position.ValidMoves(), position.Board())
	eval = -1 * math.MaxInt
	for _, move := range moves {
		e.Check_Time_Up()
		score := -1 * e.minimax_id_ab_q_searcher(position.Update(move), ply-1, !max, -1 * math.MaxInt, math.MaxInt)
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

func (e *t_engine_p_ab_q_id) minimax_id_ab_q_searcher(position *chess.Position, ply int, max bool, alpha int, beta int) (eval int) {
	explored++
	if ply == 0 {
		return e.quiescence_minimax_id_ab_q(position, 0, max, alpha, beta)
	}
	if e.time_up {
		return 0
	}

	moves := sort_moves_v1(position.ValidMoves(), position.Board())
	if len(moves) == 0 {
		return evaluate_position_v3(position, e.engine_config.ply, ply, bool_to_int(max))
	}
    for _, move := range moves {
        score := -1 * e.minimax_id_ab_q_searcher(position.Update(move), ply - 1, !max, -beta, -alpha)
		if score >= beta {
			return beta
		}
        if score > alpha {
            alpha = score
        }
    }

	return alpha
}

func (e *t_engine_p_ab_q_id) quiescence_minimax_id_ab_q(position *chess.Position, plycount int, max bool, alpha int, beta int) (eval int) {
	explored++
	q_explored++

	stand_pat := evaluate_position_v3(position, e.engine_config.ply, -plycount, bool_to_int(max))
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
        score := -1 * e.quiescence_minimax_id_ab_q(position.Update(move), plycount + 1, !max, -beta, -alpha)
		if score >= beta {
			return beta
		}
        if score > alpha {
            alpha = score
        }
    }

	return alpha
}