package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/notnil/chess"
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

Better move ordering.
   Pick and sort
   Hash MVV/LVA
// Turn iterative deepening into an engine.
// Change engine struct to be smarter.
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
	mini_challenge_manual()
	
}

func mini_challenge_manual() {
	engine := engine_minimax_id_ab_q
	engine.Set_Time(15)
	challenge_manual(&engine, chess.Black, CHESS_START_POSITION)
}

func mini_challenge_stockfish() {
	engine := engine_minimax_id_ab_q
	engine.Set_Time(15)
	challenge_stockfish(&engine, chess.White, CHESS_START_POSITION)
}

func mini_iterative_deepening_timed() {
	engine := engine_minimax_id_ab_q
	engine.Set_Time(2)
	simple_tests(&engine)
}

func mini_iterative_deepening() {
	engine := engine_minimax_plain_ab_q
	engine.Set_Config(EngineConfig{ply: 4})
	var meta_engine t_meta_engine_iterative_plain
	// meta_engine.Set_Meta_Config(MetaEngineConfig{max_depth: 8, max_time: 30})
	meta_engine.Set_Engine(&engine)
	simple_tests(&meta_engine)
}

func mini_simple_tests() {
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