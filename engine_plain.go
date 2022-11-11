package main

import (
	"log"
	"math"

	"github.com/notnil/chess"
)

// define new engine
var engine_minimax_plain Engine = Engine{
	name: "Minimax Plain",
	features: EngineFeatures{
		plain: true,
		parallel: false,
		alphabeta: false,
		iterative_deepening: false,
		mtdf: false,
	},
	engine_func: minimax_plain_engine_func,
}

func minimax_plain_engine_func(pos *chess.Position, cfg EngineConfig) (best *chess.Move, eval int) {
	best, eval = minimax_plain_starter(pos, cfg.ply, pos.Turn() == chess.White)
	log.Println("Plain minimax results", best, eval)
	return
}

func minimax_plain_starter(position *chess.Position, ply int, max bool) (best *chess.Move, eval int) {
	moves := position.ValidMoves()
	eval = -1 * math.MaxInt
	for _, move := range moves {
		tempeval := -1 * minimax_plain_searcher(position.Update(move), ply-1, !max)
		log.Println("Top Level Move:", move, "Eval:", tempeval)
		if tempeval > eval {
			log.Println("New best move:", move)
			eval = tempeval
			best = move
		}
	}
	return best, eval
}

func minimax_plain_searcher(position *chess.Position, ply int, max bool) (eval int) {
	explored++
	if ply == 0 {
		return evaluate_position_v1(position) * bool_to_int(max)
	}

	moves := position.ValidMoves()
	// if (len(moves) == 0) {
	// 	log.Println("No moves left")
	// 	log.Println(position, position.Board().Draw(), moves, position.Status(), evaluate_position_v1(position) * bool_to_int(max))
	// 	return evaluate_position_v1(position) * bool_to_int(max) // checkmate or stalemate, fix this later
	// }
    eval = -1 * math.MaxInt
    for _, move := range moves {
        tempeval := -1 * minimax_plain_searcher(position.Update(move), ply - 1, !max)
        if tempeval > eval {
            eval = tempeval
        }
    }

	return eval
}
