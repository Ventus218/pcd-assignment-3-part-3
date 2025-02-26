package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

const MAX = 10000
const N_BOTS = 6

func runBot(id int, guessChannel chan Guess, coordinationChannel chan CoordinationMsg, seeder *rand.Rand) {
	var min = 0
	var max = MAX
	var lastGuess Guess
	var r = rand.New(rand.NewPCG(seeder.Uint64(), seeder.Uint64()))

	var shouldStop = false

	for !shouldStop {
		msg := <-coordinationChannel
		if isReadyForGuessMsg(msg) {
			lastGuess = Guess{BotId: id, Value: r.IntN(max-min) + min}
			guessChannel <- lastGuess
		} else if isSuggestionMsg(msg) {
			if isBiggerSuggestionMsg(msg) {
				min = lastGuess.Value
			} else {
				max = lastGuess.Value
			}
		} else {
			// Must be a winnerMsg
			if msg.WinnerId == id {
				fmt.Println(fmt.Sprint("Bot ", id, ": I WON!!"))
			}
			shouldStop = true
		}
	}
}

func main() {
	// var r = rand.New(rand.NewPCG(123, 456)) // Seeded random
	var r = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano()))) // Unseeded random
	var v = r.IntN(MAX)
	fmt.Println(fmt.Sprint("Value to guess is: ", v))
	fmt.Println(fmt.Sprint(N_BOTS, " bots are playing"))

	var guessChannel = make(chan Guess)

	// Make coordination channels
	var coordinationChannels [N_BOTS]chan CoordinationMsg
	for i := range N_BOTS {
		coordinationChannels[i] = make(chan CoordinationMsg)
	}

	for i, ch := range coordinationChannels {
		go runBot(i, guessChannel, ch, r)
	}

	var winnerId int = -1
	var guesses [N_BOTS]Guess

	for winnerId < 0 {
		// Tell everybody the oracle is ready to receive guesses
		for _, ch := range coordinationChannels {
			ch <- readyForGuessMsg()
		}

		// Receive all guesses
		for range N_BOTS {
			guess := <-guessChannel
			fmt.Println(fmt.Sprint("Bot ", guess.BotId, " guessed ", guess.Value))
			if winnerId < 0 && guess.Value == v {
				winnerId = guess.BotId
			}
			guesses[guess.BotId] = guess
		}
		if winnerId > -1 {
			// If someone won inform everybody about it
			for i := range N_BOTS {
				coordinationChannels[i] <- winnerMsg(winnerId)
			}
		} else {
			// If no one won send suggestions
			for _, guess := range guesses {
				if guess.Value > v {
					coordinationChannels[guess.BotId] <- smallerSuggestionMsg()
				} else {
					coordinationChannels[guess.BotId] <- biggerSuggestionMsg()
				}
			}
		}
	}
	fmt.Println(fmt.Sprint("Winner is bot ", winnerId))

	// Waiting for the winner bot to print that he's won
	time.Sleep(time.Millisecond * 100)
}
