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

func test_exchange_7move(engine Engine, cfg EngineConfig) {
	test(engine, cfg, "6r1/pppk4/3p4/8/2PnPp1Q/7P/PP4r1/R5RK b - - 1 24", "g2g1")
}

func test_exchange_5move(engine Engine, cfg EngineConfig) {
	test(engine, cfg, "6r1/pppk4/3p4/8/2PnPp1Q/7P/PP6/6RK b - - 0 25", "g8g1")
}

func test_exchange_3move(engine Engine, cfg EngineConfig) {
	test(engine, cfg, "8/pppk4/3p4/8/2PnPp1Q/7P/PP6/6K1 b - - 0 26", "d4f3")
}


func test(engine Engine, cfg EngineConfig, pos string, expected string) {
	log.Println("Running test on engine:", engine.Name())
	log.Println("FEN:", pos)

	fen, err := chess.FEN(pos)
	if err != nil {
		log.Fatal(err, pos)
		panic(err)
	}
	game := chess.NewGame(fen)
	move, _ := engine.Run(game.Position(), cfg)

	if move.String() != expected {
		fmt.Println("TEST FAILED", move.String(), expected)
		log.Println("TEST FAILED", move.String(), expected)
	} else {
		fmt.Println("TEST PASSED")
		log.Println("TEST PASSED")
	}
}

type test_record struct {
	pos string
	expected string
}

func run_tests(engine Engine, cfg EngineConfig, records []test_record) {
	for _, record := range records {
		test(engine, cfg, record.pos, record.expected)
	}
	return
}

// parse EPD-file and return positions and expected moves
func parse_test_file(filename string, method func (string) (string, string)) ([]test_record) {
	// read file
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
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
		pos, move := method(record)		
		records = append(records, test_record{pos: pos, expected: move})
	}

	return records
}

// parse EPD-record and return position and expected move
func parse_epd_record(record string) (string, string) {
	record = strings.TrimSpace(record)
	parts := strings.Split(record, " ")
	expected_move := strings.TrimSuffix(parts[5], ";")
	turn := parts[1]
	move_clock := "0 1" 
	if turn == "b" {
		move_clock = "1 2"
	}
	position := strings.Join(parts[:4], " ") + " " + move_clock
	return position, expected_move
}

// parse EPD-record and return position and expected move
func parse_epd_record_off(record string) (string, string) {
	record = strings.TrimSpace(record)
	parts := strings.Split(record, " ")
	parts[9] = strings.TrimSuffix(parts[9], ";")
	return strings.Join(parts[:6], " "), parts[9]
}

func parse_fen_record(record string) (string, string) {
	record = strings.TrimSpace(record)
	parts := strings.Split(record, " ")
	return strings.Join(parts[:6], " "), parts[6]
}
