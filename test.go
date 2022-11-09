package main

import (
	"log"
	"github.com/notnil/chess"
	"fmt"
)

func test_m2(engine Engine) {
	test(engine, EngineConfig{ply: 4}, "3qr2k/pbpp2pp/1p5N/3Q2b1/2P1P3/P7/1PP2PPP/R4RK1 w - - 0 1", "d5g8")
}

func test(engine Engine, cfg EngineConfig, pos string, expected string) {
	log.Println("Running test on plain...")
	log.Println("FEN:", pos)

	fen, _ := chess.FEN(pos)
	game := chess.NewGame(fen)
	move, _ := engine.Run(game.Position(), cfg)

	if move.String() != expected {
		panic("TEST FAILED")
	} else {
		log.Println("TEST PASSED")
		fmt.Println("Test passed.")
	}
}