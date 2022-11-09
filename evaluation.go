package main

import (
	// "fmt"
	// "log"
	// "math"
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

func evaluate_position_v1(board *chess.Board) int {
	squares := board.SquareMap()
	var material int = 0
	for _, piece := range squares {
		var sign int = 1
		if piece.Color() == chess.Black {
			sign = -1
		}
		material += piece_map_v1[piece.Type()] * sign
	}

	return material
}

func evaluate_move_v1(move *chess.Move, board *chess.Board) int {
	aggr := 0
	// if the move is a capture, return the difference in material
	if move.HasTag(chess.Capture) {
		aggr += piece_map_v1[board.Piece(move.S2()).Type()] - piece_map_v1[board.Piece(move.S1()).Type()]
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

// sucks
func sort_moves_v0(moves []*chess.Move, board *chess.Board) []*chess.Move {
	if !DO_MOVE_SORTING {
		return moves
	}
	return quicksort(moves, board)
}

func quicksort(moves []*chess.Move, board *chess.Board) []*chess.Move {
	if len(moves) < 2 {
		return moves
	}
	
	pivot := moves[0]
	var left, right []*chess.Move
	for _, move := range moves[1:] {
		if evaluate_move_v1(move, board) < evaluate_move_v1(pivot, board) {
			left = append(left, move)
		} else {
			right = append(right, move)
		}
	}

	return append(quicksort(left, board), append([]*chess.Move{pivot}, quicksort(right, board)...)...)
}