package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/notnil/chess"
)

/*

// Replace position with board.
// Evaluation function.
// Add in auto-testing using EPD files.
// Test different sorting algorithms.
// Change checkmate value to not be max int.
// Engine now prioritize the shortest checkmate.
// Add in Eigenmann rapid engine test.
// Piece square tables.
// Turn iterative deepening into an engine.
// Change engine struct to be smarter.
// Make iterative deepening play in time.
// Check extentions.

Passed pawns, mobility, etc.
Endgames
	Draw detection
	Tapered eval
Clone go chess package to customize it.
Better move ordering.
   // Pick and sort, changes engine structure.
   // Hash MVV/LVA
   SEE
Killer moves.
PVS or MTD(f)
Transposition tables.
UCI compatibility. Ugh, this sucks. I might give up on this and do a web server.


Magic bitboards.
Why isn't the parallel version faster for perft?

Possible:
Aspirations?

*/

func init() {
	out("Initializing engine...")
	// create new log file that doesn't exist
	if !production_mode {
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
		out("File initialization.")
	}

	out("Version", runtime.Version())
    out("NumCPU", runtime.NumCPU())
    out("GOMAXPROCS", runtime.GOMAXPROCS(0))
	if production_mode {
		out("Production mode.")
	} else {
		out("WARNING: Development mode.")
	}
	out()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	defer exit()
	out("Running engine...")
	mini_self_challenge()
}

func mini_test_iterative() {
	game := game_from_fen("r1b1kb1r/pp2pppp/2n2n2/8/2q5/2NP1N2/PPP2PPP/R1BQK2R b KQkq - 0 7")
	engine := engine_minimax_id_ab_q
	engine.Set_Time(15)
	out(engine.Run_Engine(game.Position()))
}

func mini_iterative_deepening_timed() {
	engine := engine_0dot1
	engine.Set_Time(15)
	simple_tests(&engine)
	out("Tests passed:", tests_passed)
	out("Tests run:", tests_run)
}

func mini_self_challenge() {
	game := game_from_fen(CHESS_START_POSITION)
	subengine1 := engine_0dot1
	subengine1.Set_Time(0.1)
	engine1 := NewOpeningWrapper(&subengine1, game)
	subengine2 := engine_0dot1
	subengine2.Set_Time(15)
	engine2 := NewOpeningWrapper(&subengine2, game)
	challenge_self(engine1, engine2, game)
}

func mini_challenge_manual_opening_custom() {
	game := game_from_fen("2k5/pp3p1p/5p2/8/P1P5/7P/2Pr2q1/R6K w - - 0 27")
	subengine := engine_minimax_id_ab_q
	subengine.Set_Time(15)
	engine := NewOpeningWrapper(&subengine, game)
	challenge_manual(engine, chess.Black, game)
}

func mini_challenge_manual_opening() {
	game := game_from_fen(CHESS_START_POSITION)
	subengine := engine_minimax_id_ab_q
	subengine.Set_Time(2)
	engine := NewOpeningWrapper(&subengine, game)
	challenge_manual(engine, chess.Black, game)
}

func mini_challenge_manual() {
	engine := engine_minimax_id_ab_q
	engine.Set_Time(15)
	challenge_manual(&engine, chess.Black, game_from_fen(CHESS_START_POSITION))
}

func mini_challenge_stockfish() {
	engine := engine_minimax_id_ab_q
	engine.Set_Time(15)
	challenge_stockfish(&engine, chess.White, CHESS_START_POSITION)
}



func mini_simple_tests() {
	engine := engine_minimax_plain_ab_q
	engine.Set_Config(EngineConfig{ply: 4})
	simple_tests(&engine)
	out("Tests run:", tests_run)
	out("Tests passed:", tests_passed)
}

func exit()	{
	out("Exiting engine.")
}

func simple_tests(engine Engine) {
	test_exchange_7move(engine)
	test_exchange_5move(engine)
	test_exchange_3move(engine)
	test_m2(engine)
	run_tests(engine, parse_test_file("tests/WillsMateInThree.txt", parse_fen_record))
}

func eigenmann_tests(engine Engine) {
	run_tests(engine, parse_test_file("tests/EigenmannRapidEngineTest.txt", parse_epd_record))
}

// current will vs antikythera 2kr1b1r/ppq2ppp/2n1bn2/1R2p1B1/8/2NP1N2/P1P2PPP/3QR1K1 w - - 10 15