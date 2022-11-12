package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/notnil/chess"
)

/*

Move Sorting Test
Tests the move ordering function on a given position.

*/

func move_sort_test(position *chess.Position) {
	out(position.Board().Draw())
	moves := position.ValidMoves()
	for _, move := range moves {
		out("Top Level Move:", move, "Move order score:", evaluate_move_v1(move, position.Board()))
	}
}

/*

Manual test positions.
Simple checkmate or exchange positions.

*/

func test_m2(engine Engine) {
	test(engine, "3qr2k/pbpp2pp/1p5N/3Q2b1/2P1P3/P7/1PP2PPP/R4RK1 w - - 0 1", "d5g8")
}

func test_exchange_7move(engine Engine) {
	test(engine, "6r1/pppk4/3p4/8/2PnPp1Q/7P/PP4r1/R5RK b - - 1 24", "g2g1")
}

func test_exchange_5move(engine Engine) {
	test(engine, "6r1/pppk4/3p4/8/2PnPp1Q/7P/PP6/6RK b - - 0 25", "g8g1")
}

func test_exchange_3move(engine Engine) {
	test(engine, "8/pppk4/3p4/8/2PnPp1Q/7P/PP6/6K1 b - - 0 26", "d4f3")
}

/*

Test position and expected move.

*/

type test_record struct {
	pos string
	expected string
}


func test(engine Engine, pos string, expected string) {
	out("Running test on engine:", engine.Name())
	out("FEN:", pos)

	fen, err := chess.FEN(pos)
	if err != nil {
		log.Fatal(err, pos)
		panic(err)
	}
	game := chess.NewGame(fen)
	move, _ := engine.Run_Engine(game.Position())

	// this is to format the move in a way that is compatible with the test file
	moves := game.Position().ValidMoves()
	possible_moves := make([]string, 0)
	for _, move := range moves {
		// if last two characters are the same, save the move
		if move.String()[2:4] == expected[len(expected)-2:] {
			possible_moves = append(possible_moves, move.String())
		}
	}

	if move.String() != expected {
		if (move.String()[2:4] == expected[len(expected)-2:]) {
			if len(possible_moves) > 1 {
				out("!!! POSSIBLE PASS:", move, expected)
				out("Move is ambiguous, possible moves are:", possible_moves)
			} else {
				out("TEST PASSED")
				return 
			}
		}
		fmt.Println("TEST FAILED", move.String(), expected)
		// fmt.Println("Possible moves:", possible_moves)
		out("TEST FAILED", move.String(), expected)
	} else {
		fmt.Println("TEST PASSED")
		out("TEST PASSED")
	}
}


func run_tests(engine Engine,records []test_record) {
	for _, record := range records {
		test(engine, record.pos, record.expected)
	}
}

/*

Load and test positions from test banks.
Stored in .txt files, loaded and parsed.

*/

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

// parse normal EPD-record 
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

// parse special EPD-record
func parse_epd_record_off(record string) (string, string) {
	record = strings.TrimSpace(record)
	parts := strings.Split(record, " ")
	parts[9] = strings.TrimSuffix(parts[9], ";")
	return strings.Join(parts[:6], " "), parts[9]
}

// parse FEN-record
func parse_fen_record(record string) (string, string) {
	record = strings.TrimSpace(record)
	parts := strings.Split(record, " ")
	return strings.Join(parts[:6], " "), parts[6]
}
