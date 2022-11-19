package main

import "github.com/0hq/chess"

/*

MVV-LVA Hashtable
Hashed, used in evaluation version 2.
Indexed: Victim, Attacker

*/

// x axis is attacker: no piece, king, queen, rook, bishop, knight, pawn
var mvv_lva = [7][7]int{
	{0, 0, 0, 0, 0, 0, 0}, // victim no piece
	{0, 0, 0, 0, 0, 0, 0}, // victim king
	{0, 50, 51, 52, 53, 54, 55}, // victim queen
	{0, 40, 41, 42, 43, 44, 45}, // victim rook
	{0, 30, 31, 32, 33, 34, 35}, // victim bishop
	{0, 20, 21, 22, 23, 24, 25}, // victim knight
	{0, 10, 11, 12, 13, 14, 15}, // victim pawn
}

/*

Version 2 of the move ordering function.
Sorts moves by how likely they are to be good.

*/

func evaluate_move_v2(move *chess.Move, board *chess.Board) int {
	// add a bonus for castling
	if move.HasTag(chess.KingSideCastle) {
		return 50
	}
	// queen side castle is worth less
	if move.HasTag(chess.QueenSideCastle) {
		return 40 	
	}

	//	add a bonus for check
	if move.HasTag(chess.Check) {
		return 1 + MVV_LVA(move, board)
	}

	return MVV_LVA(move, board)
}

/*

Version 2 of the quiescence move ordering function.
Sorts moves by how likely they are to be good.

*/

func evaluate_q_move_v2(move *chess.Move, board *chess.Board) int {
	return MVV_LVA(move, board)
}

/*

Move ordering.

*/

func sort_moves_v1(moves []*chess.Move, board *chess.Board) []*chess.Move {
	if !DO_MOVE_SORTING {
		return moves
	}
	return quicksort(moves, board, evaluate_move_v2)
}


// mutates moves list
func pick_move_v1(moves []*chess.Move, board *chess.Board, start_index int) *chess.Move {
	for i := start_index; i < len(moves); i++ {
		if evaluate_move_v2(moves[i], board) > evaluate_move_v2(moves[start_index], board) {
			moves[i], moves[start_index] = moves[start_index], moves[i]
		}
	}
	return moves[start_index]
}

// mutates moves list
// quiescence search moves
func pick_qmove_v1(moves []*chess.Move, board *chess.Board, start_index int) *chess.Move {
	for i := start_index; i < len(moves); i++ {
		if evaluate_q_move_v2(moves[i], board) > evaluate_q_move_v2(moves[start_index], board) {
			moves[i], moves[start_index] = moves[start_index], moves[i]
		}
	}
	// if DO_Q_MOVE_PRUNING {
	// 	if evaluate_q_move_v2(moves[start_index], board) < 0 { // this does nothing because of evaluate_q_move_v2
	// 		return nil
	// 	}
	// }
	return moves[start_index]
}

/*

Version 1 of the move ordering function.
Sorts moves by how likely they are to be good.

*/

func evaluate_move_v1(move *chess.Move, board *chess.Board) int {
	aggr := 0
	// if the move is a capture, return the difference in material
	if move.HasTag(chess.Capture) {
		aggr += piece_map_v1[board.Piece(move.S2()).Type()] - piece_map_v1[board.Piece(move.S1()).Type()] + 1
	}

	// if the move is a promotion, return promotion value
	if move.Promo() != chess.NoPieceType {
		aggr += piece_map_v1[move.Promo()]
	}
	
	// add a bonus for castling
	if move.HasTag(chess.KingSideCastle) {
		aggr += 50
	}
	if move.HasTag(chess.QueenSideCastle) {
		aggr += 40 	// queen side castle is worth less
	}

	//	add a bonus for check
	if move.HasTag(chess.Check) {
		aggr += 10
	}

	return aggr
}


/*

Quicksort for move ordering.

*/

func quicksort(moves []*chess.Move, board *chess.Board, heuristic func(*chess.Move, *chess.Board) (int)) []*chess.Move {
	if len(moves) < 2 {
		return moves
	}
	
	pivot := moves[0]
	var left, right []*chess.Move
	for _, move := range moves[1:] {
		if heuristic(move, board) > heuristic(pivot, board) {
			left = append(left, move)
		} else {
			right = append(right, move)
		}
	}

	return append(quicksort(left, board, heuristic), append([]*chess.Move{pivot}, quicksort(right, board, heuristic)...)...)
}

/*

Quicksort function for q moves.
Removes all moves with negative scores.

*/

func quicksort_prune(moves []*chess.Move, board *chess.Board, heuristic func(*chess.Move, *chess.Board) (int)) []*chess.Move {
	if len(moves) < 2 {
		return moves
	}
	
	pivot := moves[0]
	var left, right []*chess.Move
	for _, move := range moves[1:] {
		// remove moves with negative scores
		if heuristic(move, board) < 0 {
			continue
		}
		if heuristic(move, board) > heuristic(pivot, board) {
			left = append(left, move)
		} else {
			right = append(right, move)
		}
	}

	return append(quicksort_prune(left, board, heuristic), append([]*chess.Move{pivot}, quicksort_prune(right, board, heuristic)...)...)
}

/*

Generates a list of q moves and sorts them.
Q moves may included moves that are captures, checks, or promotions.

*/

// only return moves that are captures or promotions or checks
func quiescence_moves_v1(moves []*chess.Move, board *chess.Board) []*chess.Move {
	var q_moves []*chess.Move = make([]*chess.Move, 0)
	for _, move := range moves {
		if move.HasTag(chess.Capture) { 
			q_moves = append(q_moves, move)
		}
		// if move.HasTag(chess.Check) {
		// 	q_moves = append(q_moves, move)
		// }
		// if move.Promo() != chess.NoPieceType {
		// 	q_moves = append(q_moves, move)
		// }
	}
	if !DO_Q_MOVE_SORTING {
		return q_moves
	}
	if DO_Q_MOVE_PRUNING {
		return quicksort_prune(q_moves, board, evaluate_q_move_v2) // this does nothing because evaluate_q_move_v2 is broken
	}
	return quicksort(q_moves, board, evaluate_q_move_v2)
}

/*

Version 1 of the quiescence move ordering function.
Sorts moves by how likely they are to be good.

*/

func evaluate_q_move_v1(move *chess.Move, board *chess.Board) int {
	return piece_map_v1[board.Piece(move.S2()).Type()] - piece_map_v1[board.Piece(move.S1()).Type()]
}

/*

Version 2 of the quiescence move ordering function.
Generates a list of q moves only! Doesn't sort like last version as we use pick moves func.
Q moves may included moves that are captures, checks, or promotions.

*/

// only return moves that are captures or promotions or checks
func quiescence_moves_v2(moves []*chess.Move) []*chess.Move {
	var q_moves []*chess.Move = make([]*chess.Move, 0)
	for _, move := range moves {
		if move.HasTag(chess.Capture) { 
			q_moves = append(q_moves, move)
		}
	}
	return q_moves
}


func store_killer_move(km *[2]*chess.Move, move *chess.Move) {
	if move != km[0] {
		km[1] = km[0]
		km[0] = move
	}
}

type scored_move struct {
	move *chess.Move
	score int
}

func score_moves_v1(moves []*chess.Move, board *chess.Board) []scored_move {
	scores := make([]scored_move, len(moves))
	for i := 0; i < len(moves); i++ {
		scores[i] = scored_move{moves[i], evaluate_q_move_v1(moves[i], board)}
	}
	return scores
}

func score_moves_v2(moves []*chess.Move, board *chess.Board, killer_moves [2]*chess.Move) []scored_move {
	scores := make([]scored_move, len(moves))
	for i := 0; i < len(moves); i++ {
		if moves[i].HasTag(chess.Capture) {
			scores[i] = scored_move{moves[i], MVV_LVA(moves[i], board) + MVV_LVA_OFFSET}
		} else if moves[i] == killer_moves[0] {
			scores[i] = scored_move{moves[i], KILLER_ONE_VALUE}
		} else if moves[i] == killer_moves[1] {
			scores[i] = scored_move{moves[i], KILLER_TWO_VALUE}
		} else {
			scores[i] = scored_move{moves[i], NORMAL_MOVE_VALUE}
		}
	}
	return scores
}

func score_q_moves_v2(moves []*chess.Move, board *chess.Board) []scored_move {
	scores := make([]scored_move, len(moves))
	for i := 0; i < len(moves); i++ {
		scores[i] = scored_move{moves[i], MVV_LVA(moves[i], board) + MVV_LVA_OFFSET}
	}
	return scores
}


func score_q_moves_v1(moves []*chess.Move, board *chess.Board) []scored_move {
	scores := make([]scored_move, len(moves))
	for i := 0; i < len(moves); i++ {
		scores[i] = scored_move{moves[i], evaluate_q_move_v2(moves[i], board)}
	}
	return scores
}

func MVV_LVA(move *chess.Move, board *chess.Board) int {
	return mvv_lva[board.Piece(move.S2()).Type()][board.Piece(move.S1()).Type()]
}

func pick_move_v2(moves []scored_move, board *chess.Board, start_index int) *chess.Move {
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