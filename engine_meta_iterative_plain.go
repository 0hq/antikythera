package main

import (
	"log"
	// "math"
	"time"

	"github.com/notnil/chess"
)


func iterative_deepening_v0(engine Engine, pos *chess.Position, time_control int) (output *chess.Move) {
	depth := 1
	var best *chess.Move
	var eval int
	start := time.Now()
	for {
		log.Println()
		log.Println("Iterative deepening depth", depth)
		best, eval = engine.engine_func(pos, EngineConfig{ply: depth})
		elapsed := time.Since(start)
		log.Println("Best move:", best, "Eval:", eval)
		log.Println("Depth:", depth, "Time:", elapsed, "Nodes:", explored)
		if elapsed.Seconds() > float64(time_control) {
			break
		}
		if eval > 100000 { // break on checkmate
			break
		}
		depth++
	}
	return best
}
