package main

import (
	// "fmt"
	// "log"
	// "math"

	"github.com/notnil/chess"
)

/*

Version 1 piece values.

*/

// piece weights
const pawn_v1 int = 100
const knight_v1 int = 320 
const bishop_v1 int = 330
const rook_v1 int = 500
const queen_v1 int = 900
const king_v1 int = 20000

// piece map
var piece_map_v1 map[chess.PieceType]int = map[chess.PieceType]int{
	1: king_v1,
	2: queen_v1,
	3: rook_v1,
	4: bishop_v1,
	5: knight_v1,
	6: pawn_v1,
}

/*

Version 1 of the evaluation function.
Sums simple material values and detects checkmate and stalemate.

*/

func evaluate_position_v1(pos *chess.Position) int {
	squares := pos.Board().SquareMap()
	var material int = 0
	for _, piece := range squares {
		var sign int = 1
		if piece.Color() == chess.Black {
			sign = -1
		}
		material += piece_map_v1[piece.Type()] * sign
	}

	// faster than doing two comparisons
	if pos.Status() != chess.NoMethod { 
		if pos.Status() == chess.Stalemate {
			return 0
		}
		if pos.Status() == chess.Checkmate {
			if pos.Turn() == chess.White {
				return -CHECKMATE_VALUE
			} else {
				return CHECKMATE_VALUE
			}
		}
	}

	return material
}

/*

Version 2 of the evaluation function.
Sums simple material values and detects checkmate and stalemate.
Penalizes checkmates based on ply.

*/

func evaluate_position_v2(pos *chess.Position, engine_ply int, current_ply int, flip int) int {
	squares := pos.Board().SquareMap()
	var material int = 0
	for _, piece := range squares {
		var sign int = 1
		if piece.Color() == chess.Black {
			sign = -1
		}
		material += piece_map_v1[piece.Type()] * sign
	}

	// faster than doing two comparisons
	if pos.Status() != chess.NoMethod { 
		if pos.Status() == chess.Stalemate {
			return 0
		}
		if pos.Status() == chess.Checkmate {
			// calculate ply penalty
			ply_penalty := (engine_ply - current_ply) * 1
			if pos.Turn() == chess.White {
				return -CHECKMATE_VALUE * flip + ply_penalty
			} else {
				return CHECKMATE_VALUE * flip + ply_penalty
			}
		}
	}

	return material * flip
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
	if DO_Q_MOVE_PRUNING {
		if evaluate_q_move_v2(moves[start_index], board) < 0 {
			return nil
		}
	}
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
		return quicksort_prune(q_moves, board, evaluate_q_move_v2)
	}
	return quicksort(q_moves, board, evaluate_q_move_v2)
}

/*

Version 1 of the quiescence move ordering function.
Sorts moves by how likely they are to be good.

*/

func evaluate_q_move_v1(move *chess.Move, board *chess.Board) int {
	// if the move is a capture, return the difference in material
	// if move.HasTag(chess.Capture) {
		return piece_map_v1[board.Piece(move.S2()).Type()] - piece_map_v1[board.Piece(move.S1()).Type()]
	// }

	// if the move is a promotion, return promotion value
	// if move.Promo() != chess.NoPieceType {
	// 	return piece_map_v1[move.Promo()]
	// }
	
	// return 0
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

Version 2 of the quiescence move ordering function.
Sorts moves by how likely they are to be good.

*/

func evaluate_q_move_v2(move *chess.Move, board *chess.Board) int {
	return mvv_lva[board.Piece(move.S2()).Type()][board.Piece(move.S1()).Type()]
}

/*

Version 1 Piece-Square tables.
Used in evaluation version 2.
Stored as from black's perspective.

*/

// Flip the board and return the piece square table for the other side's perspective.
var FLIP = []int{
    56, 57, 58, 59, 60, 61, 62, 63,
    48, 49, 50, 51, 52, 53, 54, 55,
    40, 41, 42, 43, 44, 45, 46, 47,
    32, 33, 34, 35, 36, 37, 38, 39,
    24, 25, 26, 27, 28, 29, 30, 31,
    16, 17, 18, 19, 20, 21, 22, 23,
     8,  9, 10, 11, 12, 13, 14, 15,
     0,  1,  2,  3,  4,  5,  6,  7,
};

var KING_MG_V1 = []int{
    0,    0,     0,     0,    0,    0,    0,    0,
    0,    0,     0,     0,    0,    0,    0,    0,
    0,    0,     0,     0,    0,    0,    0,    0,
    0,    0,     0,    20,   20,    0,    0,    0,
    0,    0,     0,    20,   20,    0,    0,    0,
    0,    0,     0,     0,    0,    0,    0,    0,
    0,    0,     0,   -10,  -10,    0,    0,    0,
    0,    0,    20,   -10,  -10,    0,   20,    0,
};

var QUEEN_MG_V1 = []int{
    -30,  -20,  -10,  -10,  -10,  -10,  -20,  -30,
    -20,  -10,   -5,   -5,   -5,   -5,  -10,  -20,
    -10,   -5,   10,   10,   10,   10,   -5,  -10,
    -10,   -5,   10,   20,   20,   10,   -5,  -10,
    -10,   -5,   10,   20,   20,   10,   -5,  -10,
    -10,   -5,   -5,   -5,   -5,   -5,   -5,  -10,
    -20,  -10,   -5,   -5,   -5,   -5,  -10,  -20,
    -30,  -20,  -10,  -10,  -10,  -10,  -20,  -30, 
};

var ROOK_MG_V1 = []int{
    0,   0,   0,   0,   0,   0,   0,   0,
   15,  15,  15,  20,  20,  15,  15,  15,
    0,   0,   0,   0,   0,   0,   0,   0,
    0,   0,   0,   0,   0,   0,   0,   0,
    0,   0,   0,   0,   0,   0,   0,   0,
    0,   0,   0,   0,   0,   0,   0,   0,
    0,   0,   0,   0,   0,   0,   0,   0,
    0,   0,   0,  10,  10,  10,   0,   0,
};

var BISHOP_MG_V1 = []int{
    -20,    0,    0,    0,    0,    0,    0,  -20,
    -15,    0,    0,    0,    0,    0,    0,  -15,
    -10,    0,    0,    5,    5,    0,    0,  -10,
    -10,   10,   10,   30,   30,   10,   10,  -10,
      5,    5,   10,   25,   25,   10,    5,    5,
      5,    5,    5,   10,   10,    5,    5,    5,
    -10,    5,    5,   10,   10,    5,    5,  -10,
    -20,  -10,  -10,  -10,  -10,  -10,  -10,  -20,
};

var KNIGHT_MG_V1 = []int{
    -20, -10,  -10,  -10,  -10,  -10,  -10,  -20,
    -10,  -5,   -5,   -5,   -5,   -5,   -5,  -10,
    -10,  -5,   15,   15,   15,   15,   -5,  -10,
    -10,  -5,   15,   15,   15,   15,   -5,  -10,
    -10,  -5,   15,   15,   15,   15,   -5,  -10,
    -10,  -5,   10,   15,   15,   15,   -5,  -10,
    -10,  -5,   -5,   -5,   -5,   -5,   -5,  -10,
    -20,   0,  -10,  -10,  -10,  -10,    0,  -20,
};

var PAWN_MG_V1 = []int{
     0,   0,   0,   0,   0,   0,   0,   0,
    60,  60,  60,  60,  70,  60,  60,  60,
    40,  40,  40,  50,  60,  40,  40,  40,
    20,  20,  20,  40,  50,  20,  20,  20,
     5,   5,  15,  30,  40,  10,   5,   5,
     5,   5,  10,  20,  30,   5,   5,   5,
     5,   5,   5, -30, -30,   5,   5,   5,
     0,   0,   0,   0,   0,   0,   0,   0,
};


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
		return 1 + mvv_lva[board.Piece(move.S2()).Type()][board.Piece(move.S1()).Type()]
	}

	return mvv_lva[board.Piece(move.S2()).Type()][board.Piece(move.S1()).Type()]
}

/*

Version 3 of the evaluation function.
Sums simple material values and detects checkmate and stalemate.
Penalizes checkmates based on ply.
Uses piece-square tables.

*/

var piece_square_map_v1 = map[chess.PieceType][]int{
	chess.Pawn:   PAWN_MG_V1,
	chess.Knight: KNIGHT_MG_V1,
	chess.Bishop: BISHOP_MG_V1,
	chess.Rook:   ROOK_MG_V1,
	chess.Queen:  QUEEN_MG_V1,
	chess.King:   KING_MG_V1,
}

func evaluate_position_v3(pos *chess.Position, engine_ply int, current_ply int, flip int) int {
	squares := pos.Board().SquareMap()
	var eval int = 0
	for square, piece := range squares {
		if piece.Color() == chess.Black {
			eval += piece_square_map_v1[piece.Type()][square] * -1
			eval += piece_map_v1[piece.Type()] * -1  
		} else {
			eval += piece_square_map_v1[piece.Type()][FLIP[square]]
			eval += piece_map_v1[piece.Type()]
		}
	}

	// faster than doing two comparisons
	if pos.Status() != chess.NoMethod { 
		if pos.Status() == chess.Stalemate {
			return 0
		}
		if pos.Status() == chess.Checkmate {
			// calculate ply penalty
			ply_penalty := (engine_ply - current_ply) * 1
			if pos.Turn() == chess.White {
				return -CHECKMATE_VALUE * flip + ply_penalty
			} else {
				return CHECKMATE_VALUE * flip + ply_penalty
			}
		}
	}

	return eval * flip
}