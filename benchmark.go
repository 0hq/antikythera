package main

import (
	"time"
	"log"
	"github.com/notnil/chess"
	"fmt"
)

// measure how long minimax_plain takes run
// returns time in seconds
func benchmark(ply int, engine Engine, pos *chess.Position) float64 {
	explored = 0
	log.Println("BEGIN BENCHMARKING -", engine.Name())
	log.Println("Ply:", ply)
	log.Println("Starting at time", time.Now())

	start := time.Now()
	move, eval := engine.Run(pos, EngineConfig{ply: ply})
	elapsed := time.Since(start)

	log.Println("Complete at time", time.Now())
	log.Println("Best move:", move)
	log.Println("Evaluation:", eval)
	log.Println("Elapsed time:", elapsed.Seconds(), "seconds")
	log.Println("Nodes explored:", explored)
	log.Println("END BENCHMARKING -")

	return elapsed.Seconds()
}

func benchmark_range(plymin int, plymax int, engine Engine, pos *chess.Position) {
	for i := plymin; i <= plymax; i++ {
		elapsed := benchmark(i, engine, pos)
		fmt.Println("Benchmark Ply:", i)
		fmt.Println("Benchmark:", explored, elapsed)
		fmt.Println("Nodes per second:", float64(explored)/elapsed)
		fmt.Println()
	}
}

func benchmark_engines(engines []Engine, pos *chess.Position) {
	for _, engine := range engines {
		benchmark_range(2, 4, engine, pos)
	}
}