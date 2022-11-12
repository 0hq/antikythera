package main

import (
	"time"

	"github.com/notnil/chess"
)

type t_meta_engine_iterative_plain struct {
	engine Engine
	config MetaEngineConfig
}

func (e *t_meta_engine_iterative_plain) Set_Engine(engine Engine) {
	e.engine = engine
}

func (e *t_meta_engine_iterative_plain) Name() string {
	return "Meta Engine Iterative Deepening v0 : " + e.engine.Name()
}

func (e *t_meta_engine_iterative_plain) Check_Time_Up() bool {
	panic("Not implemented.")
}

func (e *t_meta_engine_iterative_plain) Set_Meta_Config(c MetaEngineConfig) {
	e.config = c
}

func (e *t_meta_engine_iterative_plain) Run_Engine_Game(g *chess.Game) (*chess.Move, int) {
	return e.Run_Engine(g.Position())
}

func (e *t_meta_engine_iterative_plain) Set_Config(c EngineConfig) {
	panic("Not good practice, do not set internal engine config after meta engine is created.")
	// e.engine.Set_Config(c)
}

func (e *t_meta_engine_iterative_plain) Run_Engine(pos *chess.Position) (best *chess.Move, eval int) {
	reset_counters()
	out("Running iterative deepening v0 with engine ", e.engine.Name())
	best, eval = iterative_deepening_v0(e.engine, pos, e.config.max_time)
	out("Engine results", best, eval)
	return 
}

func iterative_deepening_v0(engine Engine, pos *chess.Position, max_time int) (best *chess.Move, eval int) {
	depth := 1
	start := time.Now()
	for {
		out()
		out("Iterative deepening depth", depth)
		engine.Set_Config(EngineConfig{ply: depth})
		best, eval = engine.Run_Engine(pos)
		elapsed := time.Since(start)
		out("Best move:", best, "Eval:", eval)
		out("Depth:", depth, "Time:", elapsed, "Nodes:", explored)
		if elapsed.Seconds() > float64(max_time) {
			break
		}
		if eval > 100000 { // break on checkmate
			break
		}
		depth++
	}
	return best, eval
}
