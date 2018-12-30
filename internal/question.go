package internal

import (
	"errors"
	"fmt"
)

type Question struct {
	QuestionerId PlayerId
	ItemIds      []ItemId
	Responses    []Response
}

type Response struct {
	ResponderId PlayerId
	ItemShown   *ItemId // nil or an ItemId
}

func createQuestion(playerId PlayerId, itemIds []ItemId) (*Question, error) {
	question := &Question{playerId, itemIds, nil}
	if question.IsValid() {
		return question, nil
	} else {
		return nil, errors.New("invalid question")
	}
}

func (question *Question) IsValid() (bool) {
	return IsItemOfType(question.ItemIds[Character], Character) &&
		IsItemOfType(question.ItemIds[Weapon], Weapon) &&
		IsItemOfType(question.ItemIds[Room], Room)
}

func (question *Question) RespondToQuestion(players []Player) ([]Response) {
	var responses []Response
	var responders []Player
	questionerId := question.QuestionerId
	responders = append(responders, players[questionerId+1:]...)
	responders = append(responders, players[0:questionerId]...)
	for _, player := range responders {
		response := &Response{}
		response.ResponderId = player.Id
		response.ItemShown = player.RespondToQuestion(question)
		responses = append(responses, *response)
		if response.ItemShown != nil {
			return responses
		}
	}
	return responses
}

func (question *Question) Print(players []Player) {
	questioner := players[question.QuestionerId].Name
	character := GetItemName(question.ItemIds[Character])
	weapon := GetItemName(question.ItemIds[Weapon])
	room := GetItemName(question.ItemIds[Room])
	responses := question.Responses
	fmt.Printf("%s asks -- Is it %s with the %s in the %s?\n", questioner, character, weapon, room)
	for _, response := range responses {
		responder := players[response.ResponderId].Name
		var itemShown string
		if response.ItemShown == nil {
			itemShown = "I can't disprove that"
		} else {
			itemShown = fmt.Sprintf("shows %s", GetItemName(*response.ItemShown))
		}
		fmt.Printf("%s responds -- %s\n", responder, itemShown)
	}
}
