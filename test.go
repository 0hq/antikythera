package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/notnil/chess"
)

func test_m2(engine Engine) {
	test(engine, EngineConfig{ply: 4}, "3qr2k/pbpp2pp/1p5N/3Q2b1/2P1P3/P7/1PP2PPP/R4RK1 w - - 0 1", "d5g8")
}

func test(engine Engine, cfg EngineConfig, pos string, expected string) {
	log.Println("Running test on engine:", engine.Name())
	log.Println("FEN:", pos)

	fen, _ := chess.FEN(pos)
	game := chess.NewGame(fen)
	move, _ := engine.Run(game.Position(), cfg)

	if move.String() != expected {
		fmt.Println("TEST FAIL")
		log.Println("TEST FAILED")
	} else {
		fmt.Println("TEST PASSED")
		log.Println("TEST PASSED")
	}
}

type test_record struct {
	pos string
	expected string
}

func run_tests(engine Engine, cfg EngineConfig, filename string) {
	records, err := parse_epd_file(filename)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, record := range records {
		test(engine, cfg, record.pos, record.expected)
	}
	return
}

// parse EPD-file and return positions and expected moves
func parse_epd_file(filename string) ([]test_record, error) {
	// read file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	// read file line by line
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// parse EPD-records
	records := make([]test_record, 0)
	for scanner.Scan() {
		// parse EPD-record using parse_epd_record
		record := scanner.Text()
		pos, move := parse_epd_record(record)
		records = append(records, test_record{pos: pos, expected: move})
	}

	return records, nil
}

// parse EPD-record and return position and expected move
func parse_epd_record(record string) (string, string) {
	record = strings.TrimSpace(record)
	parts := strings.Split(record, " ")
	fmt.Println(parts, len(parts))
	parts[9] = strings.TrimSuffix(parts[9], ";")
	return strings.Join(parts[:6], " "), parts[9]
}
