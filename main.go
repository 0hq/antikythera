package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
	// "github.com/notnil/chess"
)

/*

// Replace position with board.
// Evaluation function.
UCI compatibility. Ugh, this sucks. I might give up on this and do a web server.
Add in auto-testing using EPD files.
Test different sorting algorithms.
Why isn't the parallel version faster for perft?

*/

func init() {
	fmt.Println("Initializing engine...")
	// create new log file that doesn't exist
	for i := 0; ; i++ {
		// create file name from timestamp date and hour
		date := time.Now().Format("2006-01-02")
		filename := fmt.Sprintf("logs/%s-%d.log", date, i)
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			f, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			log.SetOutput(f)
			break
		}
	}
    log.Println("File initialization.")

	log.Println("Version", runtime.Version())
    log.Println("NumCPU", runtime.NumCPU())
    log.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	defer exit()
	fmt.Println("Running engine...")
	// test_m2(engine_minimax_plain_ab)
	// benchmark_engines(plain_engines, newGame)
	// benchmark_pll(4)
	// move_sort_test(game_from_fen("1k1r3r/pp1bbp1p/5p2/1B2n3/5B2/3P1N2/PP3PPP/R4RK1 w - - 3 16").Position())
	
	// benchmark_range(2, 6, engine_minimax_parallel_plain, game_from_fen("1k1r4/pp1b1R2/6pp/4p3/2B5/4Q3/PPP2B2/3K4 b - - 0 2").Position())
	// fmt.Println(perft(5, newGame.Clone().Position()))
	// benchmark_range(1, 7, engine_perft_pll, newGame.Clone().Position())
	

	// fmt.Println(game_from_fen("rn1qkb1r/pp2pppp/5n2/3p1b2/3P4/2N1P3/PP3PPP/R1BQKBNR w KQkq - 0 1").Position().Board().Draw())
	// fmt.Println(parse_epd_record("rnbq1rk1/pppp1ppp/4pn2/8/1bPP4/P1N5/1PQ1PPPP/R1B1KBNR b KQ - 1 5 id \"CCR06\"; bm Bcx3+;"))
	// fmt.Println(parse_epd_file("tests/test1.txt"))
	run_tests(engine_minimax_plain_ab_q, EngineConfig{ply: 5}, "tests/test1.txt")
}

func exit()	{
	fmt.Println("Exiting engine.")
	log.Println("Exiting engine.")
}

