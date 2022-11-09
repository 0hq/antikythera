package main

import (
	"fmt"
	"runtime"
	"time"
	"log"
	"os"
)


func init() {
	fmt.Println("Version", runtime.Version())
    fmt.Println("NumCPU", runtime.NumCPU())
    fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))
	runtime.GOMAXPROCS(runtime.NumCPU())

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

}

func main() {
	fmt.Println("Version", runtime.Version())
    fmt.Println("NumCPU", runtime.NumCPU())
    fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))
	// defer profile.Start().Stop()

	test_m2(engine_minimax_plain)

	// for i := 0; i < 6; i++ {
	// 	elapsed :=  benchmark(i + 1)
	// 	fmt.Println("ply:", i+1)
	// 	fmt.Println("Benchmark:", explored, elapsed)
	// 	fmt.Println("Nodes per second:", float64(explored)/elapsed, "\n")
	// }
}