package main

import (
	"fmt"
	"math"
	"github.com/notnil/chess"
)

/*

// Replace position with board.
// Evaluation function.
UCI compatibility. Ugh, this sucks. I might give up on this and do a web server.


*/

func minimax_plain_test(pos *chess.Position) (best *chess.Move) {
	move, eval := minimax_plain_starter(pos, 4, false)
	fmt.Println("Plain minimax results", move, eval)
	return move
}

func minimax_parallel_test(pos *chess.Position) (best *chess.Move) {
	move_channel := make(chan *chess.Move)
	eval_channel := make(chan int)
	go minimax_pll(pos, 4, true, nil, move_channel, eval_channel, true)
	move := <- move_channel
	eval := <- eval_channel
	fmt.Println("Parellel results", move, eval)
	return move
}


func minimax_pll(position *chess.Position, ply int, max bool, last_move *chess.Move, move_channel chan *chess.Move, eval_channel chan int, isRoot bool) {
	explored++

	// max ply reached
	if ply == 0 {
		// evaluate position and send back to parent
		move_channel <- last_move
		eval_channel <- evaluate_position(position.Board(), max)
		return
	}

	// generate moves
	var moves []*chess.Move = position.ValidMoves()
	if (isRoot) {
		fmt.Println("Moves:", moves)
	}
	var length int = len(moves)

	// create channel to pass back move and eval
	move_channel_local := make(chan *chess.Move, length)
	eval_channel_local := make(chan int, length)

	// create goroutines for each move
    for _, move := range moves {
        go minimax_pll(position.Update(move), ply-1, !max, move, move_channel_local, eval_channel_local, false)
    }

	// wait for all goroutines to finish
	var eval int = -1 * math.MaxInt
	var best *chess.Move = nil
	for i := 0; i < length; i++ {
		move := <-move_channel_local
		tempeval := <-eval_channel_local
		if tempeval > eval {
			eval = tempeval
			best = move
		}
	}

	// pass value back to parent goroutine
	move_channel <- best
	eval_channel <- eval

	return
}

// piece weights
const pawn int = 100
const knight int = 320
const bishop int = 330
const rook int = 500
const queen int = 900
const king int = 20000

// piece map
var piece_map map[chess.PieceType]int = map[chess.PieceType]int{
	1: king,
	2: queen,
	3: rook,
	4: bishop,
	5: knight,
	6: pawn,
}

func evaluate_position(board *chess.Board, max bool) int {
	squares := board.SquareMap()
	var material int = 0
	for _, piece := range squares {
		var sign int = 1
		if piece.Color() == chess.Black {
			sign = -1
		}
		material += piece_map[piece.Type()] * sign
	}

	return material
}