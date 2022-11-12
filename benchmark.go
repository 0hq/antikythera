package main

import (
	"time"

	"github.com/notnil/chess"
)

// measure how long minimax_plain takes run
// returns time in seconds
func benchmark(ply int, engine Engine, pos *chess.Position) float64 {
	reset_counters()
	out("BEGIN BENCHMARKING -", engine.Name())
	out("Ply:", ply)
	out("Starting at time", time.Now())

	start := time.Now()
	move, eval := engine.Run_Engine(pos)
	elapsed := time.Since(start)

	out("Complete at time", time.Now())
	out("Best move:", move)
	out("Evaluation:", eval)
	out("Elapsed time:", elapsed.Seconds(), "seconds")
	out("Nodes explored:", explored)
	out("END BENCHMARKING -")

	return elapsed.Seconds()
}

func benchmark_range(plymin int, plymax int, engine Engine, pos *chess.Position) {
	for i := plymin; i <= plymax; i++ {
		benchmark(i, engine, pos)
	}
}

func benchmark_engines(engines []Engine, pos *chess.Position) {
	for _, engine := range engines {
		benchmark_range(2, 4, engine, pos)
	}
}

/*

Move ordering and generation benchmarks.

*/

/*

Perft Test

Benchmarks move generation and board update.

*/

var engine_perft EngineClass = EngineClass{
	name: "Perft Test",
	features: EngineFeatures{
		plain: false,
		parallel: false,
		alphabeta: false,
		iterative_deepening: false,
		mtdf: false,
	},
}

type type_benchmark_perft EngineClass

func (e *type_benchmark_perft) Run_Engine(pos *chess.Position, cfg EngineConfig) (best *chess.Move, eval int) {
	count := perft(cfg.ply, pos)
	explored = count
	out("Perft nodes searched", count)
	return nil, count
}

func perft(ply int, pos *chess.Position) int {
	if ply == 0 {
		return 1
	}

	moves := pos.ValidMoves()
	count := 0
	for _, move := range moves {
		count += perft(ply-1, pos.Update(move))
	}
	return count
}

/*

Perft PLL Test
Runs in parallel.

Benchmarks move generation and board update.

*/

// define new engine
var engine_perft_pll EngineClass = EngineClass{
	name: "Parallel Perft Test",
	features: EngineFeatures{
		plain: false,
		parallel: true,
		alphabeta: false,
		iterative_deepening: false,
		mtdf: false,
	},
}

type type_benchmark_perft_pll EngineClass

func (e *type_benchmark_perft_pll) Run_Engine(pos *chess.Position, cfg EngineConfig) (best *chess.Move, eval int) {
	count_channel := make(chan int, 1)
	go perft_parallel(cfg.ply, pos, count_channel)
	count := <-count_channel
	explored = count
	out("Perft nodes searched", count)
	return nil, count
}

func perft_parallel(ply int, pos *chess.Position, count_channel chan int) int {
	if ply == 0 {
		return 1
	}

	moves := pos.ValidMoves()
	length := len(moves)
	count_channel_local := make(chan int, length)

	for _, move := range moves {
		go perft_parallel(ply-1, pos.Update(move), count_channel_local)
	}

	count := 0
	for i := 0; i < length; i++ {
		count += <-count_channel_local
	}
	return count
}