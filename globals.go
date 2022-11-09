package main

import ()

// constants for minimax

var explored int = 0
const ENGINE_MINIMAX_PLAIN_PLY int = 4
const ENGINE_MINIMAX_PARALLEL_PLAIN_PLY int = 4

var all_engines = []Engine{
	engine_minimax_plain_ab,
	engine_minimax_parallel_plain,
	engine_minimax_plain,
}
var plain_engines = []Engine{
	engine_minimax_plain_ab,
	engine_minimax_plain,
}