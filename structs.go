package main

import "github.com/notnil/chess"

// engine
// owns config, features, name, and start engine_func

// define new engine that houses the methods


type Engine interface {
	Run_Engine(*chess.Position) (*chess.Move, int)
	Name() string
	Set_Config(EngineConfig)
}

type EngineClass struct {
    name string
	features EngineFeatures
	engine_config EngineConfig
	time_up bool
}

type EngineFeatures struct {
	plain bool
	parallel bool
	alphabeta bool
	iterative_deepening bool
	mtdf bool
	quiescence bool
}

type EngineConfig struct {
	ply int
}

type MetaEngineConfig struct {
	max_time int
	max_depth int
}

func (e *EngineClass) Set_Config(cfg EngineConfig) {
	e.engine_config = cfg
}

func (e EngineClass) Name() string {
	return e.name
}