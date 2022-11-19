package main

import (
	"github.com/0hq/chess"
)

// transposition.go contains an implementation of a transposition table (TT) to use
// in searching and perft.

const (
	// Default size of the transposition table, in MB.
	DefaultTTSize = 64

	// The number of buckets to have per transposition table index.
	NumTTBuckets = 2

	// Constant for the size of a transposition table search and perft entry, in bytes,
	// considering memory alignment.
	SearchEntrySize uint64 = 16
	PerftEntrySize  uint64 = 24

	// Constants representing the different flags for a transposition table entry,
	// which determine what kind of entry it is. If the entry has a score from
	// a fail-low node (alpha wasn't raised), it's an alpha entry. If the entry has
	// a score from a fail-high node (a beta cutoff occured), it's a beta entry. And
	// if the entry has an exact score (alpha was raised), it's an exact entry.
	AlphaFlag uint8 = 1
	BetaFlag  uint8 = 2
	ExactFlag uint8 = 3
)

// A struct for a transposition table entry used in the search.
type SearchEntry struct {
	Hash       uint64
	Depth      int
	Score      int
	Best       chess.Move
	FlagAndAge uint8
}

// A struct for a transposition table entry used in perft.
type PerftEntry struct {
	Hash  uint64
	Nodes uint64
	Depth uint8
}

func (entry SearchEntry) GetHash() uint64 {
	return entry.Hash
}

func (entry SearchEntry) GetDepth() int {
	return entry.Depth
}

func (entry SearchEntry) GetFlag() uint8 {
	return (entry.FlagAndAge & 0xc0) >> 6

}
func (entry *SearchEntry) SetFlag(flag uint8) {
	entry.FlagAndAge &= 0x3f
	entry.FlagAndAge |= flag << 6
}

func (entry SearchEntry) GetAge() uint8 {
	return (entry.FlagAndAge & 0x30) >> 4
}

func (entry *SearchEntry) SetAge(age uint8) {
	entry.FlagAndAge &= 0xcf
	entry.FlagAndAge |= age << 4
}

func (entry *SearchEntry) Get(hash uint64, ply int, depth int, alpha int, beta int, best *chess.Move) (int, bool) {
	var adjustedScore int = 0
	shouldUse := false

	hash_reads++

	// Since index collisions can occur, test if the hash of the entry at this index
	// actually matches the hash for the current position.
	if entry.Hash == hash {

		// Even if we don't get a score we can use from the table, we can still
		// use the best move in this entry and put it first in our move ordering
		// scheme.
		if best != nil {
			*best = entry.Best
		}

		// Return the score of the position to use as an estimate for various
		// pruning and extension techniques in the search.
		adjustedScore = entry.Score

		// To be able to get an accurate value from this entry, make sure the results of
		// this entry are from a search that is equal or greater than the current
		// depth of our search.
		if entry.Depth >= depth {
			score := entry.Score

			// If the score we get from the transposition table is a checkmate score, we need
			// to do a little extra work. This is because we store checkmates in the table using
			// their distance from the node they're found in, not their distance from the root.
			// So if we found a checkmate-in-8 in a node that was 5 plies from the root, we need
			// to store the score as a checkmate-in-3. Then, if we read the checkmate-in-3 from
			// the table in a node that's 4 plies from the root, we need to return the score as
			// checkmate-in-7.
			if score > 30000 {
				score -= ply
			}

			if score < -30000 {
				score += ply
			}

			if entry.GetFlag() == ExactFlag {
				// If we have an exact entry, we can use the saved score.
				adjustedScore = score
				shouldUse = true
			}

			if entry.GetFlag() == AlphaFlag  { // && score >= alpha
				// If we have an alpha entry, and the entry's score is less than our
				// current alpha, then we know that our current alpha is the best score
				// we can get in this node, so we can stop searching and use alpha.
				adjustedScore = score
				shouldUse = true
			}

			if entry.GetFlag() == BetaFlag { //   && score <= beta
				// If we have a beta entry, and the entry's score is greater than our
				// current beta, then we have a beta-cutoff, since while
				// searching this node previously, we found a value greater than the current
				// beta. so we can stop searching and use beta.
				adjustedScore = score
				shouldUse = true
			}
		}

		// if ply == 1 {
		// 	out("Hash hit:", entry.Score, entry.Depth, entry.FlagAndAge, shouldUse)
		// }

		// if ply == 1 && !shouldUse {
		// 	out("Right depth?", entry.Depth >= depth)
		// 	out("What flag? (exact, alpha, beta)", entry.GetFlag() == ExactFlag, entry.GetFlag() == AlphaFlag, entry.GetFlag() == BetaFlag)
		// 	out("Right less than alpha?", entry.Score <= alpha, entry.Score, alpha)
		// 	out("Right greater than beta?", entry.Score >= beta, entry.Score, beta)
		// 	out()
		// }
	} else {
		if entry.Hash != 0 {
			// out(entry.Hash, hash)
			// out(entry)
			// panic("hash mismatch")
			hash_collisions++
		}
	}

	// Return the score
	return adjustedScore, shouldUse
}

func (entry *SearchEntry) Set(hash uint64, score int, best chess.Move, ply, depth int, flag, age uint8) {
	entry.Hash = hash
	entry.Depth = depth
	entry.Best = best
	entry.SetFlag(flag)
	entry.SetAge(age)

	// If the score we get from the transposition table is a checkmate score, we need
	// to do a little extra work. This is because we store checkmates in the table using
	// their distance from the node they're found in, not their distance from the root.
	// So if we found a checkmate-in-8 in a node that was 5 plies from the root, we need
	// to store the score as a checkmate-in-3. Then, if we read the checkmate-in-3 from
	// the table in a node that's 4 plies from the root, we need to return the score as
	// checkmate-in-7.
	if score > 30000 {
		score += ply
	}

	if score < -30000 {
		score -= ply
	}

	entry.Score = score
}

func (entry PerftEntry) GetHash() uint64 {
	return entry.Hash
}

func (entry PerftEntry) GetAge() uint8 {
	// Age doesn't matter when considering perft transposition table
	// entries, so 1 will always be returned, so that a current age of
	// 0 will be given when probing the table, meaning the ages
	// will never match and a replace-always scheme will still be used.
	return 1
}

func (entry PerftEntry) GetDepth() uint8 {
	return entry.Depth
}

func (entry *PerftEntry) Get(hash uint64, depth uint8) (nodeCount uint64, ok bool) {
	if entry.Hash == hash && entry.Depth == depth {
		return entry.Nodes, true
	}
	return 0, false
}

func (entry *PerftEntry) Set(hash uint64, depth uint8, nodes uint64) {
	entry.Hash = hash
	entry.Depth = depth
	entry.Nodes = nodes
}

// A struct for a transposition table.
type TransTable[Entry interface {
	SearchEntry | PerftEntry
	GetHash() uint64
	GetAge() uint8
	GetDepth() int
}] struct {
	entries []Entry
	size    uint64
}

// Resize the transposition table given what the size should be in MB.
func (tt *TransTable[Entry]) Resize(sizeInMB uint64, entrySize uint64) {
	size := (sizeInMB * 1024 * 1024) / entrySize
	tt.entries = make([]Entry, size)
	tt.size = size
}

// Get an entry from the table to use it.
func (tt *TransTable[Entry]) Probe(hash uint64) *Entry {
	// Get the entry from the table, calculating an index by modulo-ing the hash of
	// the position by the size of the table. A two-bucket system is used to
	// more efficently make use of the table.

	index := hash % tt.size
	// if index+1 == tt.size {
	// 	return &tt.entries[index]
	// }

	first := tt.entries[index]
	return &first
	// if first.GetHash() == hash {
	// 	return &tt.entries[index]
	// }

	// return &tt.entries[index+1]
}

// Get an entry from the table to store in it.
func (tt *TransTable[Entry])  Store(hash uint64, depth uint8, currAge uint8) *Entry {
	index := hash % tt.size
	// if index+1 == tt.size {
	// 	return &tt.entries[index]
	// }

	// first := tt.entries[index]
	// if first.GetDepth() <= depth { // || first.GetAge() != currAge
	// 	// Note that returning &first caused a bug where the transposition
	// 	// table entry was never modifed, and so the table was always empty.
	// 	// Have to figure out why that is, but for now return &tt.entries[index]
	// 	// directly.
	// 	return &tt.entries[index]
	// }

	return &tt.entries[index]
}

// Unitialize the memory used by the transposition table
func (tt *TransTable[Entry]) Unitialize() {
	tt.entries = nil
	tt.size = 0
}

// Clear the transposition table
func (tt *TransTable[Entry]) Clear() {
	for idx := uint64(0); idx < tt.size; idx++ {
		tt.entries[idx] = *new(Entry)
	}
}