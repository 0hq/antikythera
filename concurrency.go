package main

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/notnil/chess"
)

var explored int = 0

func main() {
	fmt.Println("Version", runtime.Version())
    fmt.Println("NumCPU", runtime.NumCPU())
    fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))
	// defer profile.Start().Stop()

	for i := 0; i < 6; i++ {
		elapsed :=  benchmark_pll(i + 1)
		fmt.Println("Depth:", i+1)
		fmt.Println("Benchmark:", explored, elapsed)
		fmt.Println("Nodes per second:", float64(explored)/elapsed, "\n")
	}
}

// measure how long minimax_plain takes run
// returns time in seconds
func benchmark(depth int) float64 {
	start := time.Now()
	minimax_plain(chess.NewGame(), depth, true)
	elapsed := time.Since(start)
	return elapsed.Seconds()
}

func benchmark_pll(depth int) float64 {
	start := time.Now()
	move_channel := make(chan *chess.Move)
	eval_channel := make(chan int)
	go minimax_pll(chess.NewGame(), depth, true, nil, move_channel, eval_channel)
	move := <-move_channel
	eval := <- eval_channel
	fmt.Println(move, eval)
	elapsed := time.Since(start)
	return elapsed.Seconds()
}


func minimax_plain(game *chess.Game, depth int, max bool) (best *chess.Move, eval int) {
	explored++
	if depth == 0 {
		return nil, evaluate_position(max)
	}

	moves := game.ValidMoves()
	
    eval = -1 * math.MaxInt
    for _, move := range moves {
        post := game.Clone()
        post.Move(move)
        _, tempeval := minimax_plain(post, depth-1, !max)
        if tempeval > eval {
            eval = tempeval
            best = move
        }
    }

	return best, eval
}

func minimax_pll(game *chess.Game, depth int, max bool, last_move *chess.Move, move_channel chan *chess.Move, eval_channel chan int) {
	explored++
	if depth == 0 {
		// fmt.Println("reached bottom")
		move_channel <- last_move
		eval_channel <- evaluate_position(max)
		return
	}
    eval := -1 * math.MaxInt
	var best *chess.Move = nil
	moves := game.ValidMoves()
	

	// create channel to pass back move and eval
	move_channel_local := make(chan *chess.Move, len(moves))
	eval_channel_local := make(chan int, len(moves))

	// create goroutines for each move
    for _, move := range moves {
        post := game.Clone()
        post.Move(move)
        go minimax_pll(post, depth-1, !max, move, move_channel_local, eval_channel_local)
    }

	// wait for all goroutines to finish
	for i := 0; i < len(moves); i++ {
		move := <-move_channel_local
		tempeval := <-eval_channel_local
		if tempeval > eval {
			eval = tempeval
			best = move
		}
	}

	// fmt.Println(best, eval, depth)
	move_channel <- best
	eval_channel <- eval
	return
}

func evaluate_position(max bool) int {
    return 0
}