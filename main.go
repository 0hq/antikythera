package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/0hq/chess"
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
// Clone go chess package to customize it.
// Killer moves.
// Mobility is done, but seems to be a massive perf. hit.
// Pick and sort, changes engine structure.
// Hash MVV/LVA
// Transposition tables.
// TT Move ordering.

Work saving.
Check extentions.
Null Move Pruning.
Investigate mobility.
Passed pawns, blocked pawns, etc.
Endgames
	Draw detection
	Tapered eval
Better move ordering.
   SEE
   History heuristic
PVS or MTD(f)
UCI compatibility. Ugh, this sucks. I might give up on this and do a web server.


Magic bitboards.
Why isn't the parallel version faster for perft?

Possible:
Aspirations?

*/

func init() {
	out("Initializing engine...")
	InitZobrist()
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
	out("Initialization complete.")
	out()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	out("Running main program.", "\n")
	defer exit()

	// engine := new_engine(&engine_0dot3dot2, 1)
	// simple_tests(engine)
	// eigenmann_tests(engine)

	// mini_performance_challenge()
	// mini_test_transposition()
	// mini_self_challenge()
	mini_challenge_stockfish()
}

func exit()	{
	out("Exiting engine.")
}

func mini_test_transposition() {
	game := game_from_fen("6r1/pppk4/3p4/8/2PnPp1Q/7P/PP4r1/R5RK b - - 1 24")
	engine2 := wrap_engine(&engine_0dot3, 5, game)
	// stop_at_depth = 7
	engine2.Run_Engine(game.Position())
	// do_tt_output = true
	// stop_at_depth = 12
	engine2.Run_Engine(game.Position())
}

func mini_performance_challenge() {
	game := game_from_fen("6r1/pppk4/3p4/8/2PnPp1Q/7P/PP4r1/R5RK b - - 1 24")
	engine1 := wrap_engine(&engine_0dot3dot2, 5, game)
	engine2 := wrap_engine(&engine_0dot3dot1, 5, game)
	engine1.Run_Engine(game.Position())
	engine2.Run_Engine(game.Position())
}

func mini_self_challenge() {
	game := game_from_fen(CHESS_START_POSITION)
	engine1 := wrap_engine(&engine_0dot3dot2, 5, game)
	engine2 := wrap_engine(&engine_0dot3dot1, 5, game)
	challenge_self(engine2, engine1, game)
}

func mini_challenge_manual() {
	game := game_from_fen(CHESS_START_POSITION)
	engine := wrap_engine(&engine_0dot2, 15, game)
	challenge_manual(engine, chess.Black, game)
}

func mini_challenge_stockfish() {
	game := game_from_fen(CHESS_START_POSITION)
	engine := wrap_engine(&engine_0dot3dot2, 15, game)
	challenge_stockfish(engine, chess.White, game)
}

func mini_simple() {
	engine := new_engine(&engine_0dot2dot1, 15)
	simple_tests(engine)
}


// current will vs antikythera 2kr1b1r/ppq2ppp/2n1bn2/1R2p1B1/8/2NP1N2/P1P2PPP/3QR1K1 w - - 10 15
// out(evaluate_position_v3(game_from_fen("rnb1k2r/ppp3pp/8/5p2/3qnP2/2N5/PPPNQbPP/R1BK1B1R w kq - 18 19").Position(), 0, 0, 1))
// test(engine, "6r1/pppk4/3p4/8/2PnPp1Q/7P/PP4r1/R5RK b - - 1 24", "g2g1", false)