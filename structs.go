package main

import (
	"time"

	"github.com/0hq/chess"
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
	Set_Time(float64)
	Reset()
	Reset_Time()
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

func (e *EngineClass) Set_Time(ms float64) {
	e.time_duration = time.Duration(ms * 1000) * time.Millisecond
	e.time_up = false
	// out("Setting time duration:", e.time_duration, "seconds:", seconds)
	// out("Setting time up:", e.time_up)
}

func (e *EngineClass) Reset() {
	e.time_up = false
}

func (e *EngineClass) Reset_Time() {
	e.time_up = false
	e.start_time = time.Now()
}


func (e *EngineClass) Set_Config(cfg EngineConfig) {
	e.engine_config = cfg
}

func (e EngineClass) Name() string {
	return e.name
}

// check time up
func (e *EngineClass) Check_Time_Up() bool {
	// out("Checking time up")
	// out("Time since start_time:", time.Since(e.start_time))
	// out("Time duration:", e.time_duration)
	if e.time_duration < time.Since(e.start_time) {
		e.time_up = true
	}
	return e.time_up
}