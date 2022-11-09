package main

import (
	"log"
	"github.com/notnil/chess"

)

func test_m2(engine Engine) {
	log.Println("Running mate in two test on plain...")
	log.Println("FEN:", "3qr2k/pbpp2pp/1p5N/3Q2b1/2P1P3/P7/1PP2PPP/R4RK1 w - - 0 1")

	fen, _ := chess.FEN("3qr2k/pbpp2pp/1p5N/3Q2b1/2P1P3/P7/1PP2PPP/R4RK1 w - - 0 1")
	game := chess.NewGame(fen)
	move, _ := engine.Run(game.Position())

	if move.String() != "d5g8" {
		panic("TEST FAILED")
	} else {
		log.Println("TEST PASSED")
	}
}