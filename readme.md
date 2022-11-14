# Antikythera

> A parrallelized Chess engine written in Golang.

This engine is a rebuild of the old engine I built this summer, as I was unhappy with the strength of the previous engine (considering it was built in 2.5 weeks). I've already massively increased the power of the engine through some bugfixing / rewriting.

Here's a quick overview if you want to read the engine code:

The files that begin with "engine\_" denote different engine versions, each with different levels of sophistication. In optimizing the engine, I've started with the simplest possible architecture to gradually improve the system. The parallelized parts are somewhat complex to wrap your head around if you don't understand Goroutines and the way Go works with concurrent operations. A good way to think about is that the engine runs all the different parts at once, then pulls them back in one by one as soon as they finish.

---    

Want to play the engine?

Download this repository.
Run ```brew install go```, if you do not already have Golang installed.
Open your terminal and navigate to the downloaded folder.
Run ```go run .``` or (```go build``` and then ```./main``` for faster speeds).

Thanks, Will.
