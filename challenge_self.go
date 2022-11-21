package main

import (
	"math"

	"github.com/0hq/chess"
)

func challenge_self(white Engine, black Engine, game *chess.Game) {
	max_moves := math.MaxInt
	// game.UseNotation(global_AlgebraicNotation)

	out("Self challenge started.")
	out("Engine 1 (white):", white.Name())
	out("Engine 2 (black):", black.Name())
	out("Game:", game.String())
	out(game.Position().Board().Draw())
	out()
	
	for game.Outcome() == chess.NoOutcome && len(game.Moves()) < max_moves {
		var move *chess.Move
		var eval int
		if game.Position().Turn() == chess.White {
			move, eval = white.Run_Engine_Game(game)
		} else {
			move, eval = black.Run_Engine_Game(game)
		}

		if move == nil {
			panic("NO MOVE")
		}

		err := game.Move(move)
		if err != nil {
			panic(err)
		}
		out(game.Position().Turn(), move, eval)
		out(game.FEN())
		out(game.String())
		out(game.Position().Board().Draw())
		out()
	}

}