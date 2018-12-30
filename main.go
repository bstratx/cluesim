package main

import (
	"project/clue/internal"
	"time"
	"math/rand"
)

func main() {
	rand.Seed(time.Now().UTC().Unix())
	secrets, remainingItems := internal.PickSecretItems()
	players := createPlayers(remainingItems)
	secrets.PrintItems("Secret Items: ", "\n")
	for _, player := range (players) {
		player.Print()
	}

	for _, player := range (players) {
		question, err := player.PoseQuestion()
		if err == nil {
			question.Responses = question.RespondToQuestion(players)
			question.Print(players)
		}
	}

}

func createPlayers(items internal.ItemIdList) ([]internal.Player) {
	return []internal.Player{
		{0, "Player 0", internal.Scarlett, items[0:3]},
		{1, "Player 1", internal.Mustard, items[3:6]},
		{2, "Player 2", internal.White, items[6:9]},
		{3, "Player 3", internal.Green, items[9:12]},
		{4, "Player 4", internal.Peacock, items[12:15]},
		{5, "Player 5", internal.Plum, items[15:18]},
	}
}