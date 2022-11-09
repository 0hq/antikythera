package main

import (
	"log"
	"github.com/notnil/chess"
)
	
type Engine struct {
    name string
	features EngineFeatures
	engine_func func(*chess.Position, EngineConfig) (*chess.Move, int)
}

type EngineFeatures struct {
	plain bool
	parallel bool
	alphabeta bool
	iterative_deepening bool
	mtdf bool
}

type EngineConfig struct {
	ply int
}

func (e *Engine) Run(pos *chess.Position, cfg EngineConfig) (best *chess.Move, eval int) {
	log.Println("Starting engine", e.name)
	log.Println("Features:", e.features)
	return e.engine_func(pos, cfg)
}

func (e *Engine) Name() string {
	return e.name
}