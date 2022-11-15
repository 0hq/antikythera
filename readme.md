# Antikythera

> A parrallelized Chess engine written in Golang.

This engine is a rebuild of the old engine I built this summer, as I was unhappy with the strength of the previous engine (considering it was built in 2.5 weeks). I've already massively increased the power of the engine through some bugfixing / rewriting.

Here's a quick overview if you want to read the engine code:

The files that begin with "engine\_" denote different engine versions, each with different levels of sophistication. In optimizing the engine, I've started with the simplest possible architecture to gradually improve the system. The parallelized parts are somewhat complex to wrap your head around if you don't understand Goroutines and the way Go works with concurrent operations. A good way to think about is that the engine runs all the different parts at once, then pulls them back in one by one as soon as they finish.

If you want to learn more about how this engine works, go to [chessengines.org](https://chessengines.org) (my website) which explains everything. I haven't got around to writing code comments yet.

---    

Want to play the engine?

*Click on releases on the repository sidebar ->*

![Releases](https://media.discordapp.net/attachments/494692084041908244/1041756240709165147/Screen_Shot_2022-11-14_at_11.45.24_AM.png)

Thanks, Will.
