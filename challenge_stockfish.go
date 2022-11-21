package main

import (
	"math"
	"time"

	"github.com/0hq/chess"
	"github.com/0hq/chess/uci"
)

func challenge_stockfish(engine Engine, play_as_color chess.Color, game *chess.Game) {
	max_moves := math.MaxInt
	// game.UseNotation(global_AlgebraicNotation)

	eng, err := uci.New("stockfish")
	if err != nil {
		panic(err)
	}
	defer eng.Close()
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		panic(err)
	}

	out("Stockfish challenge started.")
	out("Engine 1:", engine.Name(), "playing as", play_as_color)
	out("Game:", game.String())
	out(game.Position().Board().Draw())
	out()
	
	for game.Outcome() == chess.NoOutcome && len(game.Moves()) < max_moves {
		var move *chess.Move
		var eval int
		if game.Position().Turn() == play_as_color {
			move, eval = engine.Run_Engine_Game(game)
		} else {
			move, eval = stockfish(game, eng)
		}

		if move == nil {
			panic("NO MOVE")
		}

		err := game.Move(move)
		if err != nil {
			out("ERROR:", err)
			out("GAME:", game.String())
			out("MOVE:", move)
			out("FEN:", game.FEN())
			panic(err)
		}
		out(game.Position().Turn(), move, eval)
		out(game.FEN())
		out(game.String())
		out(game.Position().Board().Draw())
		out()
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