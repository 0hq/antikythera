package main

import (
	"math"
	"time"

	"github.com/0hq/chess"
)

/*

Improvements over 0.3.1.
Implements TT move ordering.

*/

type t_engine_0dot3dot2 struct {
	EngineClass
	killer_moves [MAX_DEPTH][2]*chess.Move
	current_depth int
	tt TransTable[SearchEntry]
	age uint8 // this is used to age off entries in the transposition table, in the form of a half move clock
}

var engine_0dot3dot2 = t_engine_0dot3dot2{
	EngineClass{
		name: "Engine 0.3.2",
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
	TransTable[SearchEntry]{},
	0,
} 


func (e *t_engine_0dot3dot2) Run_Engine(pos *chess.Position) (best *chess.Move, eval int) {
	Reset_Global_Counters()
	e.Reset_Time()
	e.Print_Start()
	
	e.age ^= 1

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

		out("Depth:", depth, "Nodes:", explored, "Best move:", best, "Eval:", eval, "Time:", time.Since(e.start_time), "Hash hits", hash_hits, "writes", hash_writes, "reads", hash_reads, "collisions", hash_collisions)	
	

		// break only on checkmate win, not on checkmate loss
		if eval >= CHECKMATE_VALUE { 
			break
		}

		depth++
	}

	e.Print_End(best, eval)
	return best, eval
}

func (e *t_engine_0dot3dot2) minimax_id_ab_q_starter(position *chess.Position, depth int, max bool) (best *chess.Move, eval int) {
	moves := score_moves_v2(position.ValidMoves(), position.Board(), e.killer_moves[e.current_depth - depth])
	eval = -1 * math.MaxInt // functions as alpha

	for i := 0; i < len(moves); i++ {

		if e.Check_Time_Up() {
			break
		}

		move := pick_move_v2(moves, position.Board(), i) // mutates move list, moves best move to front
		score := -1 * e.minimax_id_ab_q_searcher(position.Update(move), 1, depth-1, !max, -math.MaxInt, -eval)


		if PRINT_TOP_MOVES {
			out("Top Level Move:", move, "Eval:", score,)
		}

		if DO_DEPTH_COUNT {
			depth_count[e.current_depth - depth]++
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

func (e *t_engine_0dot3dot2) minimax_id_ab_q_searcher(position *chess.Position, ply int, depth int, max bool, alpha int, beta int) (eval int) {
	
	explored++

	if e.Check_Time_Up() {
		return 0
	}

	var tt_move *chess.Move = nil
	var hash uint64 = Zobrist.GenHash(position)
	var entry *SearchEntry = e.tt.Probe(hash)
	var ttScore, shouldUse = entry.Get(hash, ply, depth, alpha, beta, tt_move)

	if shouldUse {
		hash_hits++
		return int(ttScore)
	}

	if depth == 0 {
		return e.quiescence_minimax_id_ab_q(position, 0, max, alpha, beta)
	}
	
	var moves []scored_move = score_moves_v2(position.ValidMoves(), position.Board(), e.killer_moves[e.current_depth - depth])

	// if no moves, checkmate or stalemate
	if len(moves) == 0 {
		return evaluate_position_v3(position, e.engine_config.ply, depth, bool_to_int(max))
	}

	var tt_flag = AlphaFlag
	var best_move *chess.Move = nil
	var best_score = alpha

    for i := 0; i < len(moves); i++ {

		if DO_DEPTH_COUNT {
			depth_count[e.current_depth - depth]++
		}

		var move *chess.Move = pick_move_v2(moves, position.Board(), i)
		var score int = -1 * e.minimax_id_ab_q_searcher(position.Update(move), ply + 1, depth - 1, !max, -beta, -alpha)

		if score >= beta {

			if !move.HasTag(chess.Capture) {
				store_killer_move(&e.killer_moves[e.current_depth - depth], move)
			}

			tt_flag = BetaFlag
			best_move = move

			best_score = beta
			break
		}

        if score > alpha {

            alpha = score
			best_score = score

			tt_flag = ExactFlag
			best_move = move

        }
    }

	if !e.Check_Time_Up() {

		var entry *SearchEntry = e.tt.Store(hash, depth, e.age)
		entry.Set(hash, best_score, best_move, ply, depth, tt_flag, e.age)

		hash_writes++

	}

	return best_score
}

func (e *t_engine_0dot3dot2) quiescence_minimax_id_ab_q(position *chess.Position, depthcount int, max bool, alpha int, beta int) (eval int) {
	explored++
	q_explored++

	stand_pat := evaluate_position_v3(position, e.engine_config.ply, -depthcount, bool_to_int(max))

	if stand_pat >= beta {
        return beta;
	}

    if alpha < stand_pat {
        alpha = stand_pat;
	}
	
	moves := score_q_moves_v2(quiescence_moves_v2(position.ValidMoves()), position.Board())

	if len(moves) == 0 || depthcount > MAX_DEPTH {
		return stand_pat 
	}

    for i := 0; i < len(moves); i++ {

		move := pick_move_v2(moves, position.Board(), i)

		if move == nil { // other moves are pruned
			break
		}

		if DO_DEPTH_COUNT {
			depth_count[e.current_depth + depthcount]++
		}

        score := -1 * e.quiescence_minimax_id_ab_q(position.Update(move), depthcount + 1, !max, -beta, -alpha)

		if score >= beta {
			return beta
		}

        if score > alpha {
            alpha = score
        }

    }

	return alpha
}

func (e *t_engine_0dot3dot2) Reset() {
	e.tt.Clear()
	e.tt.Resize(64, 16)
	e.time_up = false
	e.killer_moves = [MAX_DEPTH][2]*chess.Move{}
	e.current_depth = 0
}

func (e *t_engine_0dot3dot2) Print_Start() {
	out("Starting", e.name)
	// out("Killer moves", e.killer_moves)
	out("Duration:", e.time_duration)
}

func (e *t_engine_0dot3dot2) Print_End(best *chess.Move, eval int) {
	out("Engine results", best, eval)
	out("Total nodes", explored, "Quiescence search explored", q_explored, "nodes")
	out("Depth count", depth_count)
	out("Time", time.Since(e.start_time))
	out("Killer moves", e.killer_moves)
	out()
}
