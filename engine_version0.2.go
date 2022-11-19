package main

import (
	"math"
	"time"

	"github.com/0hq/chess"
)

/*

Improvements over 0.1.
New evaluation function that examines pawn structure, mobility, and space.
SEE?
Killer moves?
Null move pruning?

*/

type t_engine_0dot2 struct {
	EngineClass
	killer_moves [MAX_DEPTH][2]*chess.Move
	current_depth int
}

// define new engine
var engine_0dot2 = t_engine_0dot2{
	EngineClass{
		name: "Engine 0.2",
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
	[MAX_DEPTH][2]*chess.Move{},
	0,
	// engine_func: minimax_id_ab_q_engine_func,
} 

func (e *t_engine_0dot2) Reset() {
	e.time_up = false
	e.killer_moves = [MAX_DEPTH][2]*chess.Move{}
	e.current_depth = 0
}

func (e *t_engine_0dot2) Run_Engine(pos *chess.Position) (best *chess.Move, eval int) {
	reset_counters()
	e.Reset_Time()
	out("Running", e.name)
	
	// out("Killer moves", e.killer_moves)
	// out("Duration:", e.time_duration)
	depth := 1
	for {
		// out()
		// out("Iterative deepening depth", depth)
		e.current_depth = depth
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
	out("Depth count", depth_count)
	return
}

func (e *t_engine_0dot2) minimax_id_ab_q_starter(position *chess.Position, ply int, max bool) (best *chess.Move, eval int) {
	moves := e.score_moves(position.ValidMoves(), position.Board())
	eval = -1 * math.MaxInt // functions as alpha
	for i := 0; i < len(moves); i++ {
		if e.Check_Time_Up() {
			break
		}
		if DO_DEPTH_COUNT {
			depth_count[e.current_depth - ply]++
		}
		move := e.pick_move(moves, position.Board(), i) // mutates move list, moves best move to front
		// out("Move", move)
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

func store_killer_move(km *[2]*chess.Move, move *chess.Move) {
	if move != km[0] {
		km[1] = km[0]
		km[0] = move
	}
}

func (e *t_engine_0dot2) minimax_id_ab_q_searcher(position *chess.Position, ply int, max bool, alpha int, beta int) (eval int) {
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
	
	moves := e.score_moves(position.ValidMoves(), position.Board())
	if len(moves) == 0 {
		return evaluate_position_v3(position, e.engine_config.ply, ply, bool_to_int(max))
	}
    for i := 0; i < len(moves); i++ {
		if DO_DEPTH_COUNT {
			depth_count[e.current_depth - ply]++
		}
		move := e.pick_move(moves, position.Board(), i) // mutates move list, moves best move to front
		var score = -1 * e.minimax_id_ab_q_searcher(position.Update(move), ply - 1, !max, -beta, -alpha)
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

type scored_move struct {
	move *chess.Move
	score int
}

func (e *t_engine_0dot2) score_moves(moves []*chess.Move, board *chess.Board) []scored_move {
	scores := make([]scored_move, len(moves))
	for i := 0; i < len(moves); i++ {
		scores[i] = scored_move{moves[i], evaluate_move_v2(moves[i], board)}
	}
	return scores
}

func (e *t_engine_0dot2) pick_move(moves []scored_move, board *chess.Board, start_index int) *chess.Move {
	best_index := start_index
	best_score := moves[start_index].score
	for i := start_index; i < len(moves); i++ {
		if moves[i].score > best_score {
			best_score = moves[i].score
			best_index = i
		}
	}

	temp := moves[start_index]
	moves[start_index] = moves[best_index]
	moves[best_index] = temp
	return moves[start_index].move
}

func (e *t_engine_0dot2) quiescence_minimax_id_ab_q(position *chess.Position, plycount int, max bool, alpha int, beta int) (eval int) {
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