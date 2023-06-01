# Antikythera

> A parallelized Chess engine written in Golang.

This engine is mainly written with educational use in mind: I start with an engine that is as simple as possible and add more complex features one-by-one. You can read through each one of the engine versions in the /engines/ folder, each with different levels of sophistication.

I sided against focusing on the process of building and optimizing the board representation or move generation as it distracts from what I think should be the main aim of learning to build engines: gaining a deeper understanding of how to recursively optimize a system as well as how simple algorithms directly translate to 'increased intelligence.'

I chose to not focus on building a game representation / move generation from scratch as it distracts from the exciting part: watching simple algorithms translate to increased intelligence.

If you want to learn more about how this engine works, visit [chessengines.org](https://chessengines.org) (my website) which explains all of these concepts from scratch.

---

**Brief Explanation of each Feature.**

Evaluation Function - This evaluates whether we're winning or losing in a given situation. This is done by counting the number of pieces each player has and assessing their value.

Basic Min/Max Search - This method recursively explores all potential future moves up to a certain "depth" or number of steps. In choosing the optimal move, the perspective alternates between us and our opponent, always presuming that each of us will respond with the best move.

Alpha-Beta Pruning - With a favorable option in play, it is possible to save time by disregarding inferior scenarios. For instance, if there's a move that captures a pawn, situations leading to losing your queen can be overlooked.

Quiescence Search - Since the search is limited to a specific depth and a position is evaluated once this limit is reached, there's a risk of misestimating a position because we can't see one move ahead. To remedy this, evaluation is only carried out when no more threatening moves are apparent.

Iterative Deepening - Chess adheres to time constraints, so the chess engine should adapt by setting the depth of search at the game's outset. This strategy allows us to always return a valid move by running searches with increasing depth until time is exhausted.

Automated Testing - Regular testing of the engine is crucial! There are numerous chess puzzles designed specifically for engines, which are an excellent resource.

Move Ordering - If we prioritize our moves to examine the best ones first, alpha-beta pruning can eliminate suboptimal options more rapidly. The sequence could include captures, non-captures, and less advantageous moves.

Fast Move Ordering - Given the abundance of less-than-optimal moves in chess, considerable time can be saved by not sorting the entire list of moves. Instead, an insertion sort can be executed each time a new position is reviewed.

Hashed MVV/LVA - A key part of move ordering is analyzing moves where weaker pieces challenge stronger ones. We can hash the results of all combinations of piece attacks for rapid retrieval.

Piece-Square Tables for Evaluation - A basic solution to account for piece position and structure is to create tables that reward pieces based on their placement on the board.

Opening Book - Both humans and engines can benefit from the centuries of knowledge about the optimal initial moves and their responses. The engine is equipped with a comprehensive memory of all opening move sequences and the best counteractions, providing an instant grandmaster-level opening preparation.

Killer Moves - These are highly effective moves that yield positive outcomes irrespective of the opponent's reaction. Such moves could be memorized to speed up the search process by eliminating poor response options.

Transposition Tables - Given that multiple paths can lead to the same position, we can save time by remembering already encountered positions. This technique is particularly useful in iterative deepening and for conserving work between turns.

TT Move Ordering/PVS - Starting the search by focusing on the best moves predicted in the previous turn can be helpful for establishing a beneficial boundary for pruning.

Draw Detection - Chess features several rules regarding draws, and keeping track of them all can be challenging.

Checkmate Pruning - This feature helps the engine prioritize checkmates based on the number of moves they require. The faster the checkmate, the better.

MTD-bi/MTD-f - These techniques rely on the continuity of chess games. If an opponent has a slight edge in one turn, they're likely to maintain it in the next. But these features are experimental and are currently not enabled.

---

Want to play the engine?

_Click on releases on the repository sidebar ->_

Thanks, Will.
