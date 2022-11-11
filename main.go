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
// Add in auto-testing using EPD files.
Test different sorting algorithms.
Why isn't the parallel version faster for perft?
// Change checkmate value to not be max int.
// Engine now prioritize the shortest checkmate.
// Add in Eigenmann rapid engine test.

Turn iterative deepening into an engine.
Change engine struct to be smarter.
Make iterative deepening play in time.

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
	log.Println()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	defer exit()
	fmt.Println("Running engine...")
	// test_m2(engine_minimax_plain_ab_q)
	// benchmark_engines(plain_engines, newGame)
	// benchmark_pll(4)
	// x move_sort_test(game_from_fen("1k1r3r/pp1bbp1p/5p2/1B2n3/5B2/3P1N2/PP3PPP/R4RK1 w - - 3 16").Position())
	// fmt.Println(game_from_fen("3q2r1/4n2k/p1p1rBpp/PpPpPp2/1P3P1Q/2P3R1/7P/1R5K w - - 1 1").Position().Turn())
	// 6r1/pppk4/3p4/8/2PnPp1Q/7P/PP4r1/R5RK b - - 1 24
	// benchmark_range(4, 6, engine_minimax_plain_ab_q, game_from_fen("6r1/pppk4/3p4/8/2PnPp1Q/7P/PP4r1/R5RK b - - 1 24").Position())
	// fmt.Println(evaluate_position_v1(game_from_fen("3q2rk/pbpp1Npp/1p6/Q5b1/2P1P3/P7/1PP2PPP/R4RK1 b - - 3 2 ").Position()))
	// benchmark_range(2, 2, engine_minimax_plain_ab_q, game_from_fen("8/pppk4/3p4/8/2P1Pp1n/7P/PP3K2/8 b - - 1 28").Position())
	// benchmark_range(4, 6, engine_minimax_plain_ab_q, game_from_fen("6r1/pppk4/3p4/8/2PnPp1Q/7P/PP6/6RK b - - 0 25").Position())

	// benchmark_range(4, 4, engine_minimax_plain_ab_q, game_from_fen("3q2r1/4n3/p1p1rBpk/PpPpPp2/1P3P2/2P3R1/7P/1R5K w - - 0 2").Position())
	// test_exchange_4move(engine_minimax_plain, EngineConfig{ply: 4})
	// test_m2(&engine_minimax_plain_ab_q)
	// position := game_from_fen("r2qk2r/pb4pp/1n2Pb2/2B2Q2/p1p5/2P5/2B2PPP/RN2R1K1 w - - 1 1").Position()
	// fmt.Println(engine.Run_Engine(position))
	// engine.Set_Config(EngineConfig{ply: 4})
	// fmt.Println(engine.Run_Engine(position))
	// move, eval := engine.Run_Engine(game_from_fen("3qr2k/pbpp2pp/1p5N/3Q2b1/2P1P3/P7/1PP2PPP/R4RK1 w - - 0 1").Position())
	// fmt.Println(parse_test_file("tests/EigenmannRapidEngineTest.txt", parse_epd_record))
	go_simple_tests()

	// simple_tests()
	// benchmark_range(6, 6, engine_minimax_plain_ab_q, game_from_fen("5rk1/5Npp/8/3Q4/8/8/8/7K w - - 0 1").Position())
	// test_m2(engine_minimax_plain)

	// fmt.Println(perft(5, newGame.Clone().Position()))
	// benchmark_range(1, 7, engine_perft_pll, newGame.Clone().Position())

	// fmt.Println(game_from_fen("rn1qkb1r/pp2pppp/5n2/3p1b2/3P4/2N1P3/PP3PPP/R1BQKBNR w KQkq - 0 1").Position().Board().Draw())
	// fmt.Println(parse_epd_record("r1bqk1r1/1p1p1n2/p1n2pN1/2p1b2Q/2P1Pp2/1PN5/PB4PP/R4RK1 w - - bm Rxf4; id \"ERET 001 - Relief\";"))
	// fmt.Println(parse_epd_file("tests/EigenmannRapidEngineTest.txt"))
	// run_tests(engine_minimax_plain_ab_q, EngineConfig{ply: 3}, "tests/EigenmannRapidEngineTest.txt")
	// iterative_deepening_v0(engine_minimax_plain_ab_q, game_from_fen("5r1k/6pp/7N/3Q4/8/8/8/7K w - - 2 2").Position(), 30)
}

func go_simple_tests() {
	engine := engine_minimax_plain_ab_q
	engine.Set_Config(EngineConfig{ply: 4})
	simple_tests(&engine)
}

func exit()	{
	fmt.Println("Exiting engine.")
	log.Println("Exiting engine.")
}

func simple_tests(engine Engine) {
	// test_exchange_7move(engine)
	// test_exchange_5move(engine)
	// test_exchange_3move(engine)
	// test_m2(engine)
	// run_tests(engine, parse_test_file("tests/WillsMateInThree.txt", parse_fen_record))
	run_tests(engine, parse_test_file("tests/EigenmannRapidEngineTest.txt", parse_epd_record))
}