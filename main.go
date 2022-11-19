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
	out()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	defer exit()
	out("Running engine...")
	engine := new_engine(&engine_0dot2, 15)

	simple_tests(engine)
}

func exit()	{
	out("Exiting engine.")
}

// current will vs antikythera 2kr1b1r/ppq2ppp/2n1bn2/1R2p1B1/8/2NP1N2/P1P2PPP/3QR1K1 w - - 10 15

func mini_self_challenge() {
	game := game_from_fen(CHESS_START_POSITION)
	engine1 := wrap_engine(&engine_0dot2, 15, game)
	engine2 := wrap_engine(&engine_0dot2, 15, game)
	challenge_self(engine1, engine2, game)
}

func mini_challenge_manual() {
	game := game_from_fen(CHESS_START_POSITION)
	engine := wrap_engine(&engine_0dot2, 15, game)
	challenge_manual(engine, chess.Black, game)
}

func mini_challenge_stockfish() {
	game := game_from_fen(CHESS_START_POSITION)
	engine := wrap_engine(&engine_0dot2, 15, game)
	challenge_stockfish(engine, chess.Black, game)
}

func mini_simple() {
	engine := new_engine(&engine_0dot2, 15)
	simple_tests(engine)
}