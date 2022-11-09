package main

import (
	"time"
	"log"
	"github.com/notnil/chess"
)

// measure how long minimax_plain takes run
// returns time in seconds
func benchmark_plain(ply int) float64 {
	log.Println("Benchmarking ply", ply)
	log.Println("Starting benchmark at time", time.Now())

	start := time.Now()
	move, eval := minimax_plain_starter(chess.NewGame().Position(), ply, true)
	elapsed := time.Since(start)

	log.Println("Best move:", move)
	log.Println("Evaluation:", eval)
	log.Println("Benchmark complete at time", time.Now())
	log.Println("Elapsed time:", elapsed.Seconds(), "seconds")

	return elapsed.Seconds()
}

func benchmark_pll(ply int) float64 {
	explored = 0
	log.Println("Benchmarking ply", ply)
	log.Println("Starting benchmark at time", time.Now())

	start := time.Now()

	move_channel := make(chan *chess.Move)
	eval_channel := make(chan int)
	go minimax_pll(chess.NewGame().Position(), ply, true, nil, move_channel, eval_channel, true)
	move := <-move_channel
	eval := <- eval_channel
	elapsed := time.Since(start)

	log.Println("Best move:", move)
	log.Println("Evaluation:", eval)
	log.Println("Benchmark complete at time", time.Now())
	log.Println("Elapsed time:", elapsed.Seconds(), "seconds")
	log.Println("Nodes per second:", float64(explored)/elapsed.Seconds())

	return elapsed.Seconds()
}