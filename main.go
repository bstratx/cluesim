package main

import (
	"github.com/bstratx/cluesim/internal"
	"time"
	"math/rand"
	"fmt"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	hiddenSecrets, remainingItems := internal.PickSecretItems()
	players := createPlayers(remainingItems)
	hiddenSecrets.PrintItems("Secret Items: ", "\n")
	for _, player := range players {
		player.Print()
	}

	solved := false
	var secrets map[internal.ItemType]internal.ItemId
	for rounds := 1; !solved && rounds < 10; rounds++ {
		fmt.Println()
		fmt.Print("Round ")
		fmt.Println(rounds)
		for _, player := range players {

			secrets, solved = player.Deductions.HasSolution()
			if solved {
				player.Deductions.PrintSolution(secrets)
				break
			}

			question, err := player.PoseQuestion()
			if err == nil {
				question.Responses = question.RespondToQuestion(players)
				question.Print(players)
				for _, player := range players {
					player.Deductions.RecordQuestionResponses(question)
				}
			}

		}

	}
	fmt.Println()
	fmt.Println("------- SOLUTION -----------")
	for _, player := range players {
		_, solved := player.Deductions.HasSolution()
		if solved {
			player.Deductions.Print()
		}
	}
	fmt.Println()
	fmt.Println("\n------- NO SOLUTION -----------")
	for _, player := range players {
		_, solved := player.Deductions.HasSolution()
		if !solved {
			player.Deductions.Print()
		}
	}
}

func createPlayers(items internal.ItemIdList) ([]internal.Player) {
	const numPlayers = 6
	players := []internal.Player{
		internal.NewPlayer(0, "Player 0", internal.Scarlett, items[0:3], numPlayers),
		internal.NewPlayer(1, "Player 1", internal.Mustard, items[3:6], numPlayers),
		internal.NewPlayer(2, "Player 2", internal.White, items[6:9], numPlayers),
		internal.NewPlayer(3, "Player 3", internal.Green, items[9:12], numPlayers),
		internal.NewPlayer(4, "Player 4", internal.Peacock, items[12:15], numPlayers),
		internal.NewPlayer(5, "Player 5", internal.Plum, items[15:18], numPlayers),
	}
	return players
}
