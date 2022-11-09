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

func benchmark_plain_ab_move_ordering() {
	fen, _ := chess.FEN("3qr2k/pbpp2pp/1p5N/3Q2b1/2P1P3/P7/1PP2PPP/R4RK1 w - - 0 1")
	game := chess.NewGame(fen)
	benchmark_range(4, 4, engine_minimax_plain, game.Clone().Position())
	DO_MOVE_SORTING = false
	benchmark_range(4, 4, engine_minimax_plain_ab, game.Clone().Position())
	DO_MOVE_SORTING = true
	benchmark_range(4, 4, engine_minimax_plain_ab, game.Clone().Position())
}