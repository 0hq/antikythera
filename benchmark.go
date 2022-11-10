package main

import (
	"fmt"
	"log"
	"time"

	"github.com/notnil/chess"
)

// measure how long minimax_plain takes run
// returns time in seconds
func benchmark(ply int, engine Engine, pos *chess.Position) float64 {
	reset_counters()
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


// define new engine
var engine_perft Engine = Engine{
	name: "Perft Test",
	features: EngineFeatures{
		plain: false,
		parallel: false,
		alphabeta: false,
		iterative_deepening: false,
		mtdf: false,
	},
	engine_func: perft_engine_func,
}

// define new engine
var engine_perft_pll Engine = Engine{
	name: "Parallel Perft Test",
	features: EngineFeatures{
		plain: false,
		parallel: true,
		alphabeta: false,
		iterative_deepening: false,
		mtdf: false,
	},
	engine_func: perft_engine_func,
}

func perft_engine_func(pos *chess.Position, cfg EngineConfig) (best *chess.Move, eval int) {
	count := perft(cfg.ply, pos)
	explored = count
	log.Println("Perft nodes searched", count)
	return nil, count
}

func perft_pll_engine_func(pos *chess.Position, cfg EngineConfig) (best *chess.Move, eval int) {
	count_channel := make(chan int, 1)
	go perft_parallel(cfg.ply, pos, count_channel)
	count := <-count_channel
	explored = count
	log.Println("Perft nodes searched", count)
	return nil, count
}

// perft test
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

// perft test for parallel
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