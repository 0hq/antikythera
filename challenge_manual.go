package main

import (
	"fmt"

	"github.com/0hq/chess"
)

func challenge_manual(engine Engine, play_as_color chess.Color, game *chess.Game) {
	max_moves := 100
	// game.UseNotation(global_AlgebraicNotation)

	out("Manual challenge started.")
	out("Engine:", engine.Name())
	out("Playing as", play_as_color)
	out("Game:", game.String())
	out("PGN:", game)
	out(game.Position().Board().Draw())
	out()
	
	for game.Outcome() == chess.NoOutcome && len(game.Moves()) < max_moves {
		var move *chess.Move
		var eval int
		if game.Position().Turn() == play_as_color {
			move, eval = engine.Run_Engine_Game(game)
		} else {
			move = take_input_move(game)
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

func take_input_move(game *chess.Game) *chess.Move {
	out("Enter move:")
	var input string
	fmt.Scanln(&input)
	move, err := global_AlgebraicNotation.Decode(game.Position(), input)
	if err != nil {
		out("Invalid move.")
		out("Did you mean?", valid_move_strings(game))
		return take_input_move(game)
	}
	return move
}		

func valid_move_strings(game *chess.Game) []string {
	moves := game.ValidMoves()
	move_strings := make([]string, len(moves))
	for i, move := range moves {
		move_strings[i] = global_AlgebraicNotation.Encode(game.Position(), move)
	}
	return move_strings
}