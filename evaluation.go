package main

import (
	// "fmt"
	// "log"
	// "math"

	"log"

	"github.com/notnil/chess"
)

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

func move_sort_test(position *chess.Position) {
	log.Println(position.Board().Draw())
	moves := position.ValidMoves()
	for _, move := range moves {
		log.Println("Top Level Move:", move, "Move order score:", evaluate_move_v1(move, position.Board()))
	}
}

func sort_moves_v1(moves []*chess.Move, board *chess.Board) []*chess.Move {
	if !DO_MOVE_SORTING {
		return moves
	}
	return quicksort(moves, board, evaluate_move_v1)
}

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
		return quicksort_prune(q_moves, board, evaluate_move_v1)
	}
	return quicksort(q_moves, board, evaluate_q_move_v1)
}


func evaluate_q_move_v1(move *chess.Move, board *chess.Board) int {
	// if the move is a capture, return the difference in material
	if move.HasTag(chess.Capture) {
		return piece_map_v1[board.Piece(move.S2()).Type()] - piece_map_v1[board.Piece(move.S1()).Type()]
	}

	// if the move is a promotion, return promotion value
	// if move.Promo() != chess.NoPieceType {
	// 	return piece_map_v1[move.Promo()]
	// }
	
	return 0
}