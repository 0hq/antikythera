# Antikythera

> A parrallelized Chess engine written in Golang.

This engine is a rebuild of the old engine I built this summer, as I was unhappy with the strength of the previous engine (considering it was built in 2.5 weeks). I've already massively increased the power of the engine through some bugfixing / rewriting.

---

Here's a quick overview if you want to read the engine code:

The files that begin with "engine\_" denote different engine versions, each with different levels of sophistication. In optimizing the engine, I've started with the simplest possible architecture to gradually improve the system. The parallelized parts are somewhat complex to wrap your head around if you don't understand Goroutines and the way Go works with concurrent operations. A good way to think about is that the engine runs all the different parts at once, then pulls them back in one by one as soon as they finish.

If you want to learn more about how this engine works, go to [chessengines.org](https://chessengines.org) (my website) which explains everything. I haven't got around to writing all too many code comments yet.

---

**Features + a tiny explainer for each.**

Evaluation function - Are we winning or losing in a certain position? Count the number of pieces each person has and how much each one is worth.  

Basic min/max search - Recursively look through all possible future moves given a certain "depth" or number of moves to look through. In order to chose the best move, we flip flop from seeing what we would do in a situation to seeing what our opponent would do, always assuming each of us will make the best response possible.   

Alpha-beta pruning - Once you have a good option on your hands, you can use it to save time by not looking at situations once you know they're worse than what you have somewhere else. For example, if you know you have a move which wins a pawn, you can ignore positions where you blunder your queen.  

Quiescence search - Because we only search at a certain depth (# of moves) and then evaluate a position once we've reached our limit, we might run into an issue where we incorrectly estimate the value of a position because we can't see one move into the future. For example, a move might be evaluated to be great since it ends with me taking my opponents pawn with my queen, but only because my search stops too early, not seeing that in the next move my queen is taken by a knight. We can solve this blindness by only evaluating positions once we know there are no more threatening moves on the horizon: once you reach the end of your normal search, extend it with a search that only considers captures and then evaluate once you know there are no captures left.

Iterative deepening - Chess is played with time constraints, so we can try to make our engine play in time by setting the # of moves we search through at the start of the game but this doesn't work very well. Some positions might be pretty simple and we finish early and some are very complicated so we might not even have time to return a move. The solution is to run searches with an increasing depth until time runs out. This way, we'll always return a valid move.

Auto-testing - It's important to test your engine! There are lists of chess puzzles out there specifically for engines which are great to use.

Move ordering - Alpha-beta pruning lets us eliminate moves that are worse than our current best option. If we could order our moves to look through the best choice first, we could eliminate alternatives way faster. We can sort captures first, moves where a weak piece takes a strong (pawn takes queen is almost always good), non-captures, and then weak moves (random queen takes pawn is usually bad).

Fast move ordering - Since there's a lot of bad moves in chess (EG bad trades: taking defended pawns with your queen), we usually are going to be eliminating a lot of moves. Sorting moves takes a lot of time, so we can speed it up by not sorting the entire move list, which we probably won't get to the end of anyways. Instead we perform a single insertion sort every time we look at a new position.

Hashed MVV/LVA - Move ordering first looks at move where weak pieces attack strong pieces. Instead of calculating this each time, we can hash the results of all piece takes piece combinations for fast access.

Piece-square tables for eval - Our engine only really cares about winning pieces at the moment which isn't very good as in Chess your position and piece structure matters a lot. A simple solution is to make tables that reward pieces depending on what square their standing on: pawns are great at the middle, knights are bad at the edge of the board, kings are good when they're in the corner. This is an easy way to hack chess engines to efficiently care about positions.

Opening book - Humans have built up hundreds of years of knowledge about the best starting moves and their responses. No human or engine can figure these out on the fly very well, so we give our chess engine a perfect memory of all opening move sequences and the best responses. Instant grandmaster opening preparation.

Killer moves - Sometimes there are very good tricky moves that work regardless of what the opponent chooses to do after it. For example, discovered attacks where if you move a piece to open up another piece to attack on the next turn, if the opponent does not respond, they lose badly. By remembering these quiet moves that end up eliminating a lot of possible moves from our opponent we can use them to save time looking through poor possible responses.

Transposition tables - You can often end up at the same position by many different manners. Sometimes you make an exchange where you can take with your knight or your bishop first and it doesn't matter either way. Instead of doing the same calculation twice, we can remember positions that we've already been to and save time. This is a big deal for iterative deepening, where we can use transposition tables to never have to do the same work twice. This also helps since you save the work you did last turn to the next, so you don't start from the beginning.

TT move ordering/PVS - Alpha-beta pruning really likes when we already know a good option that helps us eliminate all the crappy ones: chances are that the situation that we looked at last turn (I take your pawn, you take my pawn, I take back, you take back) stays the same the next turn. By always starting by looking at the best moves we predicted last turn we can easily establish a good bound to use for pruning.

Draw detection - Chess has a lot of rules on how you can create draws. Tracking them all isn't super easy!

Checkmate pruning - Make the engine weight checkmates based on how many moves they take. Faster wins.  

MTD-bi/MTD-f - Chess games usually aren't super dramatic. If your opponent is winning by +0.5 one turn, chances are they're be winning by around that much next turn. Rather than searching through all possible moves and establishing a bound, we can use our guess of the evaluation from last turn to eliminate a lot of situations. NOTE: Remember that this doesn't always work! Maybe one move you think you're winning and the next move you realize you're about to be checkmated, if something unexpected happens you'll have to search as you would normally. These are no enabled at the moment as they are experimental.

---    

Want to play the engine?

*Click on releases on the repository sidebar ->*

Thanks, Will.
