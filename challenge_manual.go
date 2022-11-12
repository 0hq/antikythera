package main

import (
	"fmt"

	"github.com/notnil/chess"
)

func challenge_manual(engine Engine, play_as_color chess.Color, game *chess.Game) {
	max_moves := 100

	out("Manual challenge started.")
	out("Engine:", engine.Name())
	out("Playing as", play_as_color)
	out()
	
	for game.Outcome() == chess.NoOutcome && len(game.Moves()) < max_moves {
		var move *chess.Move
		var eval int
		if game.Position().Turn() == play_as_color {
			move, eval = engine.Run_Engine_Game(game)
		} else {
			var input string
			fmt.Scanln(&input)
			move, _ = global_UCINotation.Decode(game.Position(), input) 
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
		out(game.Position().Board().Draw())
		out()
	}

}