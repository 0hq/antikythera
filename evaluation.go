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

func evaluate_position_v1(board *chess.Board, max bool) int {
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