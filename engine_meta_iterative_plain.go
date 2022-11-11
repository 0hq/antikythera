package main

// // define new engine
// var engine_meta_iterative_plain Engine = Engine{
// 	name: "Iterative Deepening Plain",
// 	features: EngineFeatures{
// 		plain: true,
// 		parallel: false,
// 		alphabeta: true,
// 		iterative_deepening: true,
// 		mtdf: false,
// 	},
// 	engine_func: engine_meta_iterative_plain_func,
// 	engine_config: EngineConfig{},
// 	time_up: false,
// }

// // make method of engine

// func engine_meta_iterative_plain_func(pos *chess.Position, cfg EngineConfig) (best *chess.Move, eval int) {
// 	// fmt.Println(time_up)
// 	// best, eval = minimax_plain_ab_q_starter(pos, cfg.ply, pos.Turn() == chess.White)
// 	log.Println("Plain minimax results", best, eval)
// 	log.Println("Quiescence search explored", q_explored, "nodes")
// 	return
// }

// func iterative_deepening_v0(engine Engine, pos *chess.Position, time_control int) (output *chess.Move) {
// 	depth := 1
// 	var best *chess.Move
// 	var eval int
// 	start := time.Now()
// 	for {
// 		log.Println()
// 		log.Println("Iterative deepening depth", depth)
// 		best, eval = engine.engine_func(pos, EngineConfig{ply: depth})
// 		elapsed := time.Since(start)
// 		log.Println("Best move:", best, "Eval:", eval)
// 		log.Println("Depth:", depth, "Time:", elapsed, "Nodes:", explored)
// 		if elapsed.Seconds() > float64(time_control) {
// 			break
// 		}
// 		if eval > 100000 { // break on checkmate
// 			break
// 		}
// 		depth++
// 	}
// 	return best
// }
