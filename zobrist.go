package main

import "github.com/notnil/chess"

// Borrowed from https://github.com/algerbrex/blunder/blob/main/engine/zobrist.go
// Thanks Algerbrex, I didn't want to have to deal with debugging my own implementation.

// zobrist.go contains an interface for creating and incrementally updating
// the zobrist hash value of a given position.
//
// https://www.chessprogramming.org/Zobrist_Hashing

// An implementation of a xorshift pseudo-random number
// generator for 64 bit numbers, based on the implementation
// by Stockfish.
type PseduoRandomGenerator struct {
	state uint64
}

// Seed the generator.
func (prng *PseduoRandomGenerator) Seed(seed uint64) {
	prng.state = seed
}

// Generator a random 64 bit number.
func (prng *PseduoRandomGenerator) Random64() uint64 {
	prng.state ^= prng.state >> 12
	prng.state ^= prng.state << 25
	prng.state ^= prng.state >> 27
	return prng.state * 2685821657736338717
}

// Generate a random 64 bit number with few bits. This method is
// useful in finding magic numbers faster for generating slider
// attacks.
func (prng *PseduoRandomGenerator) SparseRandom64() uint64 {
	return prng.Random64() & prng.Random64() & prng.Random64()
}

const (
	// A constant representing the value to use when seeding the random
	// numbers generated for zobrist hashing.
	ZobristSeedValue = 1

	// A constant which represents when there is no ep square. This value indexes
	// into _Zobrist.epFileRand64 to return a 0, which will not affect the zobrist
	// hash.
	NoEPFile = 8
)

// A constant which will be a singleton of the _Zobrist struct below,
// since only one instance is ever needed.
var Zobrist _Zobrist

// A struct which holds the random numbers for the zobrist hashing, and
// has methods to create and incrementally update the hashs.
type _Zobrist struct {
	// Each aspect of the board needs to be a given a unique random 64-bit number
	// that will be xor-ed together with other unique random numbers from the positions
	// other aspects to create a unique zobrist hash.
	//
	// Blunder uses 794 unique random numbers:
	// * 12x64 for each type of piece on each possible square.
	// * 8 for each possible en passant file
	// * 16 for each possible permutation of the castling right bits.
	// * 1 for when it's white to move

	pieceSqRand64        [768]uint64
	epFileRand64         [9]uint64
	castlingRightsRand64 [16]uint64
	sideToMoveRand64     uint64
}

// Populate the zobrist arrays with random 64-bit numbers.
func (zobrist *_Zobrist) init() {
	var prng PseduoRandomGenerator
	prng.Seed(ZobristSeedValue)

	for index := 0; index < 768; index++ {
		zobrist.pieceSqRand64[index] = prng.Random64()
	}

	for index := 0; index < 8; index++ {
		zobrist.epFileRand64[index] = prng.Random64()
	}

	zobrist.epFileRand64[NoEPFile] = 0

	for index := 0; index < 16; index++ {
		zobrist.castlingRightsRand64[index] = prng.Random64()
	}

	zobrist.sideToMoveRand64 = prng.Random64()
}

// Get the unique random number corresponding to the piece type, piece color, and square
// given.
func (zobrist *_Zobrist) PieceNumber(pieceType, pieceColor uint8, sq uint8) uint64 {
	return zobrist.pieceSqRand64[(uint16(pieceType - 1)*2+uint16(pieceColor - 1))*64+uint16(sq)]
}

// Get the unique random number corresponding to the en passant square
// given.
func (zobrist *_Zobrist) EPNumber(epSq uint8) uint64 {
	return zobrist.epFileRand64[fileOfEP(epSq)]
}

// Get the unique random number corresponding to castling bits permutation
// given.
func (zobrist *_Zobrist) CastlingNumber(castlingRights uint8) uint64 {
	return zobrist.castlingRightsRand64[castlingRights]
}

// Get the unique random number corresponding to the side to move given.
func (zobrist *_Zobrist) SideToMoveNumber() uint64 {
	return zobrist.sideToMoveRand64
}

// Generate a zobrist hash from scratch for the given position.
// Useful for creating hashs when loading in FEN strings and
// debugging zobrist hashing itself.
func (zobrist *_Zobrist) GenHash(pos *chess.Position) (hash uint64) {
	squares := pos.Board().SquareMap()
	for square, piece := range squares {
		hash ^= zobrist.PieceNumber(uint8(piece.Type()), uint8(piece.Color()), uint8(square))
	}

	// hash ^= zobrist.EPNumber(uint8(pos.EnPassantSquare()))
	// hash ^= zobrist.CastlingNumber(pos.CastleRights())

	if pos.Turn() == chess.White {
		hash ^= zobrist.SideToMoveNumber()
	}

	return hash
}

// Precomputing all possible en passant file numbers
// is much more efficent for Blunder than calculating
// them on the fly.
var PossibleEPFiles [65]uint8 = [65]uint8{
	8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8,
	0, 1, 2, 3, 4, 5, 6, 7,
	8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8,
	0, 1, 2, 3, 4, 5, 6, 7,
	8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8,
	8,
}

func fileOfEP(sq uint8) uint8 {
	return PossibleEPFiles[sq]
}

func InitZobrist() {
	Zobrist = _Zobrist{}
	Zobrist.init()
}