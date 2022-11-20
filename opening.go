package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/0hq/chess"
	"github.com/0hq/chess/opening"
)

/*

Opening Book
This is old code, so it is not documented. Sorry.

*/

type OpeningWrapper struct {
	engine Engine
	book *opening.BookECO
	game *chess.Game
	opening bool
}

func NewOpeningWrapper(engine Engine, game *chess.Game) *OpeningWrapper {
	return &OpeningWrapper{
		engine: engine,
		book: opening.NewBookECO(),
		game: game,
		opening: len(game.Moves()) == 0 && game.FEN() == chess.StartingPosition().String(),
	}
}

// func (o OpeningWrapper) Set_Time(time int) {
// 	o.engine.Set_Time(time)
// }

func (o *OpeningWrapper) Set_Time(time float64) {
	o.engine.Set_Time(time)
} 

func (o *OpeningWrapper) Reset_Time() {
	o.engine.Reset_Time()
}

func (o *OpeningWrapper) Reset(pos *chess.Position) {
	o.engine.Reset(pos)
}

func (o OpeningWrapper) Run_Engine(position *chess.Position) (*chess.Move, int) {
	return o.engine.Run_Engine(position)
}

func (o *OpeningWrapper) Run_Engine_Game(g *chess.Game) (*chess.Move, int) {
	if o.opening {
		uci := get_opening_uci(g, 0)
		if uci == "" {
			o.opening = false
			return o.engine.Run_Engine(g.Position())
		} else {
			move, _ := global_UCINotation.Decode(g.Position(), uci)
			return move, 0
		}
	}
	return o.engine.Run_Engine(g.Position())
}

func (o OpeningWrapper) Name() string {
	return o.engine.Name()
}

// Set_Config
func (o *OpeningWrapper) Set_Config(config EngineConfig) {
	o.engine.Set_Config(config)
}

// Check_Time_Up
func (o OpeningWrapper) Check_Time_Up() bool {
	return o.engine.Check_Time_Up()
}


func test_opening(){
    g := chess.NewGame(chess.UseNotation(chess.UCINotation{}))
	// g.MoveStr("e2e4")
	// g.MoveStr("e6")
	opening := true
	for opening {
		move := get_opening_uci(g, 0)
		if move == "" {
			opening = false
		}
		tmove, _ := global_UCINotation.Decode(g.Position(), move)
		g.Move(tmove)
		fmt.Println(g.Position().Board().Draw())
	}
}

func get_opening_uci(g *chess.Game, retries int) string {
	book := opening.NewBookECO()
	moves := g.Moves()
	if len(moves) == 0 {
		return "e2e4"
	}
	o := book.Find(moves) // find current opening
	if o == nil {
		return ""
	}
	p := book.Possible(g.Moves()) // all openings available
	if len(p) > 0 {
		r := p[rand.Intn(len(p))] // random opening available
		split := strings.Split(r.PGN(), " ")
		if len(split) <= len(g.Moves()) {
			if len(p) == 1 || retries > 3 {
				return ""
			}
			return get_opening_uci(g, retries + 1)
		}
		m := split[len(g.Moves())]
		return m
	}
	return ""
}