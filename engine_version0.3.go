package main

import (
	"math"
	"time"

	"github.com/notnil/chess"
)

/*

Improvements over engine_iterative_alpha_beta_q.
Starter now maintains an alpha value (massive improvement).
Picks moves and sorts as you go, instead of sorting all moves at the start.

*/

type t_engine_0dot3 struct {
	EngineClass
	tt TransTable[SearchEntry]
}

// define new engine
var engine_0dot3 = t_engine_0dot3{
	EngineClass{
		name: "Engine 0.3",
		features: EngineFeatures{
			plain: true,
			parallel: false,
			alphabeta: true,
			iterative_deepening: true,
			mtdf: false,
		},
		engine_config: EngineConfig{
			ply: 0,
		},
		time_up: false,
	},
	TransTable[SearchEntry]{},
	// engine_func: minimax_id_ab_q_engine_func,
}

func (e *t_engine_0dot3) Reset() {
	e.tt.Clear()
	e.tt.Resize(64, 16)
}

func (e *t_engine_0dot3) Run_Engine(pos *chess.Position) (best *chess.Move, eval int) {
	reset_counters()
	out("Running", e.name, "as player", pos.Turn())
	e.time_up = false
	e.start_time = time.Now()
	// out("Duration:", e.time_duration)
	depth := 1
	for {
		// out()
		// out("Iterative deepening depth", depth)
		t_best, t_eval := e.minimax_id_ab_q_starter(pos, depth, pos.Turn() == chess.White)
		if e.Check_Time_Up() {
			// out("Time up, returning best move so far.")
			break
		} else {
			best = t_best
			eval = t_eval
		}
		out("Depth:", depth, "Nodes:", explored, "Best move:", best, "Eval:", eval, "Time:", time.Since(e.start_time))

		// out("Time since start_time:", time.Since(e.start_time))
		if eval >= 30000 { // break on checkmate win
			break
		}
		depth++
	}
	out()
	out("Engine results", best, eval)
	out("Total nodes", explored, "Quiescence search explored", q_explored, "nodes")
	return
}

func (e *t_engine_0dot3) minimax_id_ab_q_starter(position *chess.Position, ply int, max bool) (best *chess.Move, eval int) {
	moves := position.ValidMoves()
	eval = -1 * math.MaxInt // functions as alpha
	for i := 0; i < len(moves); i++ {
		if e.Check_Time_Up() {
			break
		}
		move := pick_move_v1(moves, position.Board(), i) // mutates move list, moves best move to front
		score := -1 * e.minimax_id_ab_q_searcher(position.Update(move), ply-1, !max, -math.MaxInt, -eval)
		if PRINT_TOP_MOVES {
			out("Top Level Move:", move, "Eval:", score,)
		}
		if score > eval {
			if PRINT_TOP_MOVES {
				out("New best move:", move, "Eval:", score)
			}
			eval = score
			best = move
		}
	}
	return best, eval
}

func (e *t_engine_0dot3) minimax_id_ab_q_searcher(position *chess.Position, ply int, max bool, alpha int, beta int) (eval int) {
	explored++
	if ply == 0 {
		return e.quiescence_minimax_id_ab_q(position, 0, max, alpha, beta)
	}
	if e.Check_Time_Up() {
		return 0
	}

	hash := Zobrist.GenHash(position)
	entry := e.tt.Probe(hash)
	ttScore, shouldUse := entry.Get(hash, uint8(ply), uint8(ply), int16(alpha), int16(beta), nil)
	if shouldUse {
		return int(ttScore)
	}

	moves := position.ValidMoves()
	if len(moves) == 0 {
		return evaluate_position_v3(position, e.engine_config.ply, ply, bool_to_int(max))
	}
	var tt_flag = AlphaFlag
	var best_move *chess.Move = nil
	var best_score = alpha
    for i := 0; i < len(moves); i++ {
		move := pick_move_v1(moves, position.Board(), i) // mutates move list, moves best move to front
        score := -1 * e.minimax_id_ab_q_searcher(position.Update(move), ply - 1, !max, -beta, -alpha)
		if score >= beta {
			tt_flag = BetaFlag
			best_score = beta // fail hard beta-cutoff
			best_move = move
			break
		}
        if score > alpha {
			tt_flag = ExactFlag
			alpha = score
			best_score = score
        }
    }

	// // If we're not out of time, store the result of the search for this position.
	if !e.Check_Time_Up() && best_move != nil {
		entry := e.tt.Store(hash, uint8(ply), 0)
		entry.Set(
			hash, int16(best_score), *best_move, uint8(ply), uint8(ply), tt_flag, 0,
		)
	}

	return best_score
}

func (e *t_engine_0dot3) quiescence_minimax_id_ab_q(position *chess.Position, plycount int, max bool, alpha int, beta int) (eval int) {
	explored++
	q_explored++

	stand_pat := evaluate_position_v3(position, e.engine_config.ply, -plycount, bool_to_int(max))
	if stand_pat >= beta {
        return beta;
	}
    if alpha < stand_pat {
        alpha = stand_pat;
	}

	moves := quiescence_moves_v2(position.ValidMoves())

	if len(moves) == 0 || plycount > MAX_DEPTH {
		return stand_pat
	}

    for i := 0; i < len(moves); i++ {
		move := pick_qmove_v1(moves, position.Board(), i) // mutates move list, moves best move to front
		if move == nil { // other moves are pruned
			break
		}
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