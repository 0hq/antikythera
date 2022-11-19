package main

import "github.com/0hq/chess"

// "fmt"
// "log"
// "math"

// "github.com/0hq/chess"

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

/*

Version 4 of the evaluation function.
Considers mobility.

*/

func evaluate_position_v4(pos *chess.Position, engine_ply int, current_ply int, flip int) int {
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

	// mobility
	mobility := 0
	if pos.Turn() == chess.White {
		mobility += len(pos.ValidMoves())
		mobility -= len(pos.NullMove().ValidMoves())
	} else {
		mobility += len(pos.NullMove().ValidMoves())
		mobility -= len(pos.ValidMoves())
	}
	eval += mobility * 1
	// out("mobility:", mobility)

	return eval * flip
}