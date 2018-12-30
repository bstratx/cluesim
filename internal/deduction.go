package internal

type Deductions struct {
	isEliminated []bool
	playerItemMatrix PlayerItemMatrix
}

type PlayerItemMatrix struct {
	playerIds       []PlayerId
	hasItem         [][]PlayerHasItemType
	hasAtLeastOneOf [][]ItemIdList
}

type PlayerHasItemType int

const (
	Unknown PlayerHasItemType = iota
	Yes
	No
	Maybe
)

type HasSomeItems struct {
	itemIdList ItemIdList
}

func (matrix *PlayerItemMatrix) playerHasItem(playerId PlayerId, itemId ItemId) (bool) {
	return matrix.hasItem[playerId][itemId] == Yes
}

func (matrix *PlayerItemMatrix) RecordQuestionResponses(question Question) {
	for _, response := range question.Responses {
		if response.ItemShown == nil {
			matrix.RecordPlayerDoesntHaveItemList(response.ResponderId, question.ItemIds)
		} else {
			matrix.RecordPlayerHasItem(response.ResponderId, *response.ItemShown)
		}
	}
}

func (matrix *PlayerItemMatrix) RecordPlayerHasItem(playerId PlayerId, itemId ItemId) {
	matrix.hasItem[playerId][itemId] = Yes
	others := matrix.otherPlayers(playerId)
	for _, otherPlayerId := range others {
		matrix.hasItem[otherPlayerId][itemId] = No
	}
}

func (matrix *PlayerItemMatrix) RecordPlayerHasItemList(playerId PlayerId, itemIdList ItemIdList) {
	for _, itemId := range itemIdList {
		matrix.RecordPlayerHasItem(playerId, itemId)
	}
}

func (matrix *PlayerItemMatrix) RecordPlayerDoesntHaveItem(playerId PlayerId, itemId ItemId) {
	matrix.hasItem[playerId][itemId] = No
}

func (matrix *PlayerItemMatrix) RecordPlayerDoesntHaveItemList(playerId PlayerId, itemIdList ItemIdList) {
	for _, itemId := range itemIdList {
		matrix.RecordPlayerDoesntHaveItem(playerId, itemId)
	}
}

func (matrix *PlayerItemMatrix) DeduceNoPlayerHasItem(itemId ItemId) (bool) {
	for _, playerId := range matrix.playerIds {
		if matrix.hasItem[playerId][itemId] != Yes {
			return false
		}
	}
	return true
}

func (matrix *PlayerItemMatrix) Print(playerId PlayerId) {

}

func (matrix *PlayerItemMatrix) otherPlayers(playerId PlayerId) ([]PlayerId) {
	var others []PlayerId
	others = append(others, matrix.playerIds[0:playerId] ...)
	others = append(others, matrix.playerIds[playerId+1:] ...)
	return others
}

func (deductions *Deductions) RecordEliminated(itemId ItemId) {
}
