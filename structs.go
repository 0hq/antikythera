package main

import (
	"time"

	"github.com/notnil/chess"
)

// engine
// owns config, features, name, and start engine_func
// define new engine that houses the methods


type Engine interface {
	Run_Engine(*chess.Position) (*chess.Move, int)
	Run_Engine_Game(*chess.Game) (*chess.Move, int)
	Name() string
	Set_Config(EngineConfig)
	Check_Time_Up() bool
}

type EngineClass struct {
    name string
	features EngineFeatures
	engine_config EngineConfig
	start_time time.Time
	time_duration time.Duration
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

func (e *EngineClass) Run_Engine(pos *chess.Position) (*chess.Move, int) {
	panic("Run_Engine not implemented on null engine.")
}

func (e *EngineClass) Run_Engine_Game(game *chess.Game) (*chess.Move, int) {
	return e.Run_Engine(game.Position())
}

func (e *EngineClass) Set_Time(seconds int) {
	e.time_duration = time.Duration(seconds) * time.Second
	e.time_up = false
}

func (e *EngineClass) Set_Config(cfg EngineConfig) {
	e.engine_config = cfg
}

func (e EngineClass) Name() string {
	return e.name
}

// check time up
func (e EngineClass) Check_Time_Up() bool {
	if time.Since(e.start_time) > e.time_duration {
		e.time_up = true
	}
	return e.time_up
}