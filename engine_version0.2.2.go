package main

import (
	"math"
	"time"

	"github.com/0hq/chess"
)

/*

Improvements over 0.2.1
New evaluation function with mobility.
Sucks.

*/

type t_engine_0dot2dot2 struct {
	EngineClass
	killer_moves [MAX_DEPTH][2]*chess.Move
	current_depth int
}

var engine_0dot2dot2 = t_engine_0dot2dot2{
	EngineClass{
		name: "Engine 0.2.2 (deprecated)",
		features: EngineFeatures{
			plain: true,
			parallel: false,
			alphabeta: true,
			iterative_deepening: true,
			mtdf: false,
		},
		engine_config: EngineConfig{0}, // redundant
		time_up: false,
	},
	[MAX_DEPTH][2]*chess.Move{},
	0,
} 

func (e *t_engine_0dot2dot2) Run_Engine(pos *chess.Position) (best *chess.Move, eval int) {
	Reset_Global_Counters()
	e.Reset_Time()
	e.Print_Start()

	depth := 1
	for {
		// for killer moves and depth counting
		e.current_depth = depth 
		
		// temporary storage, in case time runs out
		t_best, t_eval := e.minimax_id_ab_q_starter(pos, depth, pos.Turn() == chess.White)

		if e.Check_Time_Up() {
			break
		}

		// update best move and eval
		best = t_best
		eval = t_eval

		out("Depth:", depth, "Nodes:", explored, "Best move:", best, "Eval:", eval, "Time:", time.Since(e.start_time))		

		// break only on checkmate win, not on checkmate loss
		if eval >= CHECKMATE_VALUE { 
			break
		}

		depth++
	}

	e.Print_End(best, eval)
	return best, eval
}

func (e *t_engine_0dot2dot2) minimax_id_ab_q_starter(position *chess.Position, ply int, max bool) (best *chess.Move, eval int) {
	moves := score_moves_v2(position.ValidMoves(), position.Board(), e.killer_moves[e.current_depth - ply])
	eval = -1 * math.MaxInt // functions as alpha

	for i := 0; i < len(moves); i++ {

		if e.Check_Time_Up() {
			break
		}

		move := pick_move_v2(moves, position.Board(), i) // mutates move list, moves best move to front
		score := -1 * e.minimax_id_ab_q_searcher(position.Update(move), ply-1, !max, -math.MaxInt, -eval)


		if PRINT_TOP_MOVES {
			out("Top Level Move:", move, "Eval:", score,)
		}

		if DO_DEPTH_COUNT {
			depth_count[e.current_depth - ply]++
		}

		if score > eval {
			eval = score
			best = move

			if PRINT_TOP_MOVES {
				out("New best move:", move, "Eval:", score)
			}
		}
	}
	
	return best, eval
}

func (e *t_engine_0dot2dot2) minimax_id_ab_q_searcher(position *chess.Position, ply int, max bool, alpha int, beta int) (eval int) {
	explored++

	// extend search if in check, do not enter quiescence search
	// Why is this broken?
	// if position.IsInCheck() {
	// 	ply++
	// }

	if ply == 0 {
		return e.quiescence_minimax_id_ab_q(position, 0, max, alpha, beta)
	}
	if e.Check_Time_Up() {
		return 0
	}
	
	moves := score_moves_v2(position.ValidMoves(), position.Board(), e.killer_moves[e.current_depth - ply])

	// if no moves, checkmate or stalemate
	if len(moves) == 0 {
		return evaluate_position_v4(position, e.engine_config.ply, ply, bool_to_int(max))
	}

    for i := 0; i < len(moves); i++ {

		if DO_DEPTH_COUNT {
			depth_count[e.current_depth - ply]++
		}

		move := pick_move_v2(moves, position.Board(), i)
		score := -1 * e.minimax_id_ab_q_searcher(position.Update(move), ply - 1, !max, -beta, -alpha)

		if score >= beta {
			store_killer_move(&e.killer_moves[e.current_depth - ply], move)
			return beta
		}

        if score > alpha {
            alpha = score
        }
    }

	return alpha
}

func (e *t_engine_0dot2dot2) quiescence_minimax_id_ab_q(position *chess.Position, plycount int, max bool, alpha int, beta int) (eval int) {
	explored++
	q_explored++

	stand_pat := evaluate_position_v4(position, e.engine_config.ply, -plycount, bool_to_int(max))

	if stand_pat >= beta {
        return beta;
	}

    if alpha < stand_pat {
        alpha = stand_pat;
	}
	
	moves := score_q_moves_v2(quiescence_moves_v2(position.ValidMoves()), position.Board())

	if len(moves) == 0 || plycount > MAX_DEPTH {
		return stand_pat 
	}

    for i := 0; i < len(moves); i++ {

		move := pick_move_v2(moves, position.Board(), i)

		if move == nil { // other moves are pruned
			break
		}

		if DO_DEPTH_COUNT {
			depth_count[e.current_depth + plycount]++
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

func (e *t_engine_0dot2dot2) Reset(pos *chess.Position) {
	e.time_up = false
	e.killer_moves = [MAX_DEPTH][2]*chess.Move{}
	e.current_depth = 0
}

func (e *t_engine_0dot2dot2) Print_Start() {
	out("Starting", e.name)
	// out("Killer moves", e.killer_moves)
	out("Duration:", e.time_duration)
}

func (e *t_engine_0dot2dot2) Print_End(best *chess.Move, eval int) {
	out("Engine results", best, eval)
	out("Total nodes", explored, "Quiescence search explored", q_explored, "nodes")
	out("Depth count", depth_count)
	out("Time", time.Since(e.start_time))
	out("Killer moves", e.killer_moves)
	out()
}
