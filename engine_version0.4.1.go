package main

import (
	"math"
	"time"

	"github.com/0hq/chess"
)

/*

Experimental Engine
Adds MTD(bi)

*/

type t_engine_0dot4dot1 struct {
	EngineClass
	killer_moves [MAX_DEPTH][2]*chess.Move
	current_depth int
	tt TransTable[SearchEntry]
	age uint8 // this is used to age off entries in the transposition table, in the form of a half move clock
	zobristHistory [1024]uint64 // draw detection history
	zobristHistoryPly uint16 // draw detection ply
	best_guess int // score guess for mtdf
	last_depth int // for score guess replacement
}

func new_engine_0dot4dot1() t_engine_0dot4dot1 {
	return t_engine_0dot4dot1{
		EngineClass{
			name: "Engine 0.4.0 (experimental: mtdbi)",
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
		[1024]uint64{},
		0,
		0,
		0,
	} 
}

func (e *t_engine_0dot4dot1) Run_Engine(pos *chess.Position) (best *chess.Move, eval int) {
	e.Reset_Time()
	e.Print_Start()
	
	e.age ^= 1

	depth := 1
	for {

		Reset_Global_Counters()
		Reset_Hash_Counters()

		// for killer moves and depth counting
		e.current_depth = depth 
		
		// temporary storage, in case time runs out
		t_best, t_eval := e.mtdbi(pos, depth, pos.Turn() == chess.White)

		if e.Check_Time_Up() {
			break
		}

		// update best move and eval
		best = t_best
		eval = t_eval

		e.Print_Iterative_Deepening(depth, best, eval)

		if t_best == nil {
			panic("Nil move.")
		}
	
		// break only on checkmate win, not on checkmate loss
		if eval >= CHECKMATE_VALUE { 
			break
		}

		if depth >= e.last_depth - 1 {
			e.last_depth = depth
			e.best_guess = eval
		}

		depth++
	}

	e.last_depth = depth
	e.best_guess = eval
	e.Print_End(best, eval)
	return best, eval
}

func (e *t_engine_0dot4dot1) mtdbi(position *chess.Position, depth int, max bool) (*chess.Move, int) {
	var eval int
	upper_bound := CHECKMATE_VALUE
	lower_bound := -CHECKMATE_VALUE
	
	for lower_bound < upper_bound - EVAL_ROUGHNESS {
		var gamma int = (lower_bound + upper_bound + 1) / 2
		_, eval = e.minimax_id_ab_q_starter(position, depth, max, gamma, gamma + 1)

		if eval > gamma {
			lower_bound = eval
		} else {
			upper_bound = eval
		}
	}
	
	// if eval == nil {
	// 	panic()
	// }

	// get best move from transposition table
	best_move := e.tt.Probe(Zobrist.GenHash(position)).Best
	
	return &best_move, eval
}


func (e *t_engine_0dot4dot1) minimax_id_ab_q_starter(position *chess.Position, depth int, max bool, alpha int, beta int) (*chess.Move, int) {

	explored++

	if e.Check_Time_Up() {
		return nil, 0
	}

	var hash uint64 = Zobrist.GenHash(position)
	var entry *SearchEntry = e.tt.Probe(hash)
	var tt_score, should_use, tt_move = entry.Get(hash, 0, depth, -math.MaxInt, math.MaxInt)

	if should_use {
		hash_hits++
		return tt_move, tt_score
	}
	
	var moves []scored_move = score_moves_v3(position.ValidMoves(), position.Board(), e.killer_moves[e.current_depth - depth], tt_move)

	var tt_flag = AlphaFlag
	var best_move *chess.Move = nil
	var best_score = alpha

	for i := 0; i < len(moves); i++ {

		if DO_DEPTH_COUNT {
			depth_count[e.current_depth - depth]++
		}

		var move *chess.Move = pick_move_v2(moves, position.Board(), i)
		var updated_position = position.Update(move)

		var updated_hash = Zobrist.GenHash(updated_position)
		e.Add_Zobrist_History(updated_hash)

		var score int = -1 * e.minimax_id_ab_q_searcher(position.Update(move), 1, depth-1, !max, -beta, -alpha)

		e.Remove_Zobrist_History()

		if PRINT_TOP_MOVES {
			out("Top Level Move:", move, "Eval:", score,)
		}

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

			if PRINT_TOP_MOVES {
				out("New best move:", move, "Eval:", alpha)
			}

		}

	}

	if !e.Check_Time_Up() {

		var entry *SearchEntry = e.tt.Store(hash, depth, e.age)
		entry.Set(hash, best_score, best_move, 0, depth, tt_flag, e.age)

		hash_writes++

	}
	
	return best_move, best_score
}

func (e *t_engine_0dot4dot1) minimax_id_ab_q_searcher(position *chess.Position, ply int, depth int, max bool, alpha int, beta int) (eval int) {
	
	explored++

	if e.Check_Time_Up() {
		return 0
	}

	var hash uint64 = Zobrist.GenHash(position)
	var entry *SearchEntry = e.tt.Probe(hash)
	var tt_score, should_use, tt_move = entry.Get(hash, ply, depth, alpha, beta)

	if should_use {
		hash_hits++
		return tt_score
	}

	if depth == 0 {
		return e.quiescence_minimax_id_ab_q(position, 0, max, alpha, beta)
	}

	if e.Is_Draw_By_Repition(hash) {
		return 0
	}
	
	var moves []scored_move = score_moves_v3(position.ValidMoves(), position.Board(), e.killer_moves[e.current_depth - depth], tt_move)

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
		var updated_position = position.Update(move)
		
		var updated_hash = Zobrist.GenHash(updated_position)
		e.Add_Zobrist_History(updated_hash)

		var score int = -1 * e.minimax_id_ab_q_searcher(updated_position, ply + 1, depth - 1, !max, -beta, -alpha)

		e.Remove_Zobrist_History()

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

func (e *t_engine_0dot4dot1) quiescence_minimax_id_ab_q(position *chess.Position, depthcount int, max bool, alpha int, beta int) (eval int) {
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

// adds to zobrist history, which is used for draw detection
func (e *t_engine_0dot4dot1) Add_Zobrist_History(hash uint64) {
	e.zobristHistoryPly++
	e.zobristHistory[e.zobristHistoryPly] = hash
}

// decrements ply counter, which means history will be overwritten
func (e *t_engine_0dot4dot1) Remove_Zobrist_History() {
	e.zobristHistoryPly--
}

func (e *t_engine_0dot4dot1) Is_Draw_By_Repition(hash uint64) bool {
	for i := uint16(0); i < e.zobristHistoryPly; i++ {
		if e.zobristHistory[i] == hash {
			return true
		}
	}
	return false
}

func (e *t_engine_0dot4dot1) Reset(position *chess.Position) {
	e.tt.Clear()
	e.tt.Resize(64, 16)
	e.time_up = false
	e.killer_moves = [MAX_DEPTH][2]*chess.Move{}
	e.current_depth = 0
	e.best_guess = 0
	e.last_depth = 0
	e.zobristHistoryPly = 0
	e.zobristHistory[e.zobristHistoryPly] = Zobrist.GenHash(position)
}


func (e *t_engine_0dot4dot1) Print_Iterative_Deepening(depth int, best *chess.Move, eval int) {
	if QUIET_MODE {
		return
	}
	out("Depth:", depth, "Nodes:", explored, "Best move:", best, "Eval:", eval, "Time:", time.Since(e.start_time), "Hash hits", hash_hits, "writes", hash_writes, "reads", hash_reads, "collisions", hash_collisions)	
}

func (e *t_engine_0dot4dot1) Print_Start() {
	if QUIET_MODE {
		return
	}
	out("Starting", e.name)
	// out("Killer moves", e.killer_moves)
	out("Duration:", e.time_duration)
	out("Guess", e.best_guess, "from", e.last_depth)
}

func (e *t_engine_0dot4dot1) Print_End(best *chess.Move, eval int) {
	if QUIET_MODE {
		return
	}
	out("Engine results", best, eval)
	out("Nodes searched", explored, "Quiescence search explored", q_explored, "nodes")
	// out("Depth count", depth_count)
	out("Time", time.Since(e.start_time))
	// out("Killer moves", e.killer_moves)
	out()
}
