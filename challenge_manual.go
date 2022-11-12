package main

import (
	"fmt"
	"log"

	"github.com/notnil/chess"
)

func challenge_manual(engine Engine, play_as_color chess.Color, starting_position string) {
	max_moves := 100
	game := game_from_fen(starting_position)

	log.Println("Manual challenge started.")
	log.Println("Engine:", engine.Name())
	log.Println("Playing as", play_as_color)
	log.Println()

	var test chess.UCINotation 
	
	for game.Outcome() == chess.NoOutcome && len(game.Moves()) < max_moves {
		var move *chess.Move
		var eval int
		if game.Position().Turn() == play_as_color {
			move, eval = engine.Run_Engine(game.Position())
		} else {
			var input string
			fmt.Scanln(&input)
			move, _ = test.Decode(game.Position(), input) 
		}

		if move == nil {
			panic("NO MOVE")
		}

		err := game.Move(move)
		if err != nil {
			panic(err)
		}
		fmt.Println(game.Position().Turn(), move, eval)
		fmt.Println(game.Position().Board().Draw())
		fmt.Println()
		log.Println(game.Position().Turn(), move, eval)
		log.Println(game.Position().Board().Draw())
		log.Println()
	}

}