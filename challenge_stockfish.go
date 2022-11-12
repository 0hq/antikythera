package main

import (
	"log"
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

func challenge_stockfish(engine Engine, play_as_color chess.Color, starting_position string) {
	max_moves := 20
	game := game_from_fen(starting_position)

	eng, err := uci.New("stockfish")
	if err != nil {
		panic(err)
	}
	defer eng.Close()
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		panic(err)
	}

	log.Println("Stockfish challenge started.")
	log.Println("Engine:", engine.Name())
	log.Println("Playing as", play_as_color)
	log.Println()
	
	for game.Outcome() == chess.NoOutcome && len(game.Moves()) < max_moves {
		var move *chess.Move
		var eval int
		if game.Position().Turn() == play_as_color {
			move, eval = engine.Run_Engine(game.Position())
		} else {
			move, eval = stockfish(game, eng)
		}

		if move == nil {
			panic("NO MOVE")
		}

		game.Move(move)
		log.Println(game.Position().Turn(), move, eval)
		log.Println(game.Position().Board().Draw())
		log.Println()
	}

}

func stockfish(game *chess.Game, eng *uci.Engine) (*chess.Move, int) {
	cmdPos := uci.CmdPosition{Position: game.Position()}
	cmdGo := uci.CmdGo{MoveTime: time.Second / 100}
	if err := eng.Run(cmdPos, cmdGo); err != nil {
		panic(err)
	}
	return eng.SearchResults().BestMove, eng.SearchResults().Info.Score.CP / 100
}