package main

import (
	"log"
	"github.com/notnil/chess"
)
	
type Engine struct {
    name string
	features EngineFeatures
	engine_func func(*chess.Position) (*chess.Move, int)
}

type EngineFeatures struct {
	plain bool
	parallel bool
	alphabeta bool
	iterative_deepening bool
	mtdf bool
}

func (e *Engine) Run(pos *chess.Position) (best *chess.Move, eval int) {
	log.Println("Starting engine", e.name)
	log.Println("Features:", e.features)
	return e.engine_func(pos)
}

// 