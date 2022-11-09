package main

import (
	"fmt"
	"runtime"
	"time"
	"log"
	"os"
	// "github.com/notnil/chess"

)

/*

// Replace position with board.
// Evaluation function.
UCI compatibility. Ugh, this sucks. I might give up on this and do a web server.

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
	test_m2(engine_minimax_parallel_plain)
	// benchmark_range(2, 6, engine_minimax_parallel_plain, chess.NewGame().Position())
	// benchmark_pll(4)
}

func exit()	{
	fmt.Println("Exiting engine.")
	log.Println("Exiting engine.")
}