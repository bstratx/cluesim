package internal

import (
	"fmt"
	"strings"
	"os"
)

type Deductions struct {
	thisPlayerId     PlayerId
	itemIsEliminated []bool
	playerIsKnown    []bool
	matrix           *playerItemMatrix
}

type playerItemMatrix struct {
	numPlayers      int
	hasItem         [][]PlayerHasItemType
	hasAtLeastOneOf [][]ItemIdList
}

type PlayerHasItemType int

const (
	Unknown PlayerHasItemType = iota // default value
	Yes
	No
	Maybe
)

type HasSomeItems struct {
	itemIdList ItemIdList
}

func NewDeductions(thisPlayerId PlayerId, numPlayers int) (*Deductions) {

	var numItems = len(GetItemsOfType(AllItems))

	matrix := new(playerItemMatrix)
	hasItem := make([][]PlayerHasItemType, numPlayers)
	for playerId := 0; playerId < numPlayers; playerId++ {
		hasItem[playerId] = make([]PlayerHasItemType, numItems)
	}
	matrix.numPlayers = numPlayers
	matrix.hasItem = hasItem
	matrix.hasAtLeastOneOf = nil

	deductions := new(Deductions)
	deductions.matrix = matrix
	deductions.itemIsEliminated = make([]bool, numItems)
	deductions.playerIsKnown = make([]bool, numPlayers)
	deductions.thisPlayerId = thisPlayerId
	return deductions
}

func (deductions *Deductions) RecordPlayerItems(itemList ItemIdList) {
	deductions.recordPlayerHasItemList(deductions.thisPlayerId, itemList)
}

func (deductions *Deductions) RecordQuestionResponses(question *Question) {
	for _, response := range question.Responses {
		if response.ItemShown == nil {
			deductions.recordPlayerDoesntHaveItemList(response.ResponderId, question.ItemIds)
		} else {
			if question.QuestionerId == deductions.thisPlayerId {
				deductions.recordPlayerHasItem(response.ResponderId, *response.ItemShown)
			} else {
				deductions.recordPlayerMaybeHasItemList(response.ResponderId, question.ItemIds)
			}
		}
	}
	deductions.analyze()
}

func (deductions *Deductions) GetEliminatedItems() (ItemIdList) {
	var eliminatedIds ItemIdList
	for itemId, isEliminated := range deductions.itemIsEliminated {
		if isEliminated {
			eliminatedIds = append(eliminatedIds,ItemId(itemId))
		}
	}
	return eliminatedIds
}

func (deductions *Deductions) playerHasItem(playerId PlayerId, itemId ItemId) (bool) {
	return deductions.matrix.hasItem[playerId][itemId] == Yes
}

func (deductions *Deductions) indeterminatePlayerItem(playerId PlayerId, itemId ItemId) (bool) {
	return !(deductions.matrix.hasItem[playerId][itemId] == Yes ||
		deductions.matrix.hasItem[playerId][itemId] == No)
}

func (deductions *Deductions) analyze() {
	secrets, solved := deductions.deduceSecrets()
	if solved {
		fmt.Print("Solved!!! - ")
		for _, secret := range secrets {
			fmt.Print(" ")
			fmt.Print(GetItemName(ItemId(secret)))
		}
		fmt.Println()
		deductions.Print()
		os.Exit(0)
	}
}

func (deductions *Deductions) eliminateItem(itemId ItemId) {
	deductions.itemIsEliminated[itemId] = true
}

func (deductions *Deductions) recordPlayerMaybeHasItem(playerId PlayerId, itemId ItemId) {
	matrix := deductions.matrix
	currentState := matrix.hasItem[playerId][itemId]
	if currentState == Unknown {
		matrix.hasItem[playerId][itemId] = Maybe
	}
}

func (deductions *Deductions) recordPlayerMaybeHasItemList(playerId PlayerId, itemIdList ItemIdList) {
	matrix := deductions.matrix
	count := 0
	var lastItemId ItemId
	for _, itemId := range itemIdList {
		if matrix.hasItem[playerId][itemId] != No {
			deductions.recordPlayerMaybeHasItem(playerId, itemId)
			lastItemId = itemId
			count++
		}
	}
	if count == 1 {
		deductions.recordPlayerHasItem(playerId, lastItemId)
	}
}

func (deductions *Deductions) recordPlayerHasItem(playerId PlayerId, itemId ItemId) {
	matrix := deductions.matrix
	matrix.hasItem[playerId][itemId] = Yes
	deductions.itemIsEliminated[itemId] = true

	others := matrix.otherPlayers(playerId)
	for _, otherPlayerId := range others {
		matrix.hasItem[otherPlayerId][itemId] = No
	}

	if deductions.allPlayerItemsAreKnown(playerId) {
		deductions.markRemainingItemsForPlayer(playerId)
	}
}

func (deductions *Deductions) recordPlayerHasItemList(playerId PlayerId, itemIdList ItemIdList) {
	for _, itemId := range itemIdList {
		deductions.recordPlayerHasItem(playerId, itemId)
	}
}

func (deductions *Deductions) recordPlayerDoesntHaveItem(playerId PlayerId, itemId ItemId) {
	deductions.matrix.hasItem[playerId][itemId] = No
}

func (deductions *Deductions) recordPlayerDoesntHaveItemList(playerId PlayerId, itemIdList ItemIdList) {
	for _, itemId := range itemIdList {
		deductions.recordPlayerDoesntHaveItem(playerId, itemId)
	}
}

func (deductions *Deductions) allPlayerItemsAreKnown(playerId PlayerId) (bool) {

	if deductions.playerIsKnown[playerId] {
		return false
	}

	const maxItemCount = 3 // TODO: generalize
	hasCount := 0
	allItemIds := GetItemsOfType(AllItems)
	for _, itemId := range allItemIds {
		if deductions.playerHasItem(playerId, itemId) {
			hasCount++
		}
	}
	return hasCount == maxItemCount
}

func (deductions *Deductions) markRemainingItemsForPlayer(playerId PlayerId) {
	for _, itemId := range GetItemsOfType(AllItems) {
		if !deductions.playerHasItem(playerId, itemId) {
			deductions.recordPlayerDoesntHaveItem(playerId, itemId)
		}
	}
	deductions.playerIsKnown[playerId] = true
}

func (deductions *Deductions) deduceSecrets() (map[ItemType]ItemId, bool) {
	someSecrets := deductions.deduceFinalItemByType()
	mergedSecrets := deductions.deduceItemsNoPlayerHas()
	for key, value := range someSecrets {
		mergedSecrets[key] = value
	}
	return mergedSecrets, len(mergedSecrets) == len(ItemTypes)
}

func (deductions *Deductions) deduceItemsNoPlayerHas() (map[ItemType]ItemId) {
	secrets := make(map[ItemType]ItemId)
	itemIds := GetItemsOfType(AllItems)
	for _, itemId := range itemIds {
		if deductions.deduceNoPlayerHasItem(itemId) {
			secrets[GetItemType(itemId)] = itemId
		}
	}
	return secrets
}

func (deductions *Deductions) deduceNoPlayerHasItem(itemId ItemId) (bool) {
	matrix := deductions.matrix
	for playerId := 0; playerId < matrix.numPlayers; playerId++ {
		if matrix.hasItem[playerId][itemId] != No {
			return false
		}
	}
	return true
}

func (deductions *Deductions) deduceFinalItemByType() (map[ItemType]ItemId) {
	secrets := make(map[ItemType]ItemId)
	for _, itemType := range ItemTypes {
		itemIds := GetItemsOfType(itemType)
		var uneliminatedCount = 0
		var lastItemId ItemId
		for _, itemId := range itemIds {
			if deductions.itemIsEliminated[itemId] == false {
				lastItemId = itemId
				uneliminatedCount++
			}
		}
		if uneliminatedCount == 1 {
			secrets[GetItemType(lastItemId)] = lastItemId
		}
	}
	return secrets
}


func (matrix *playerItemMatrix) printMatrixHeader(currentPlayerId PlayerId) {
	var builder strings.Builder
	builder.WriteString("\nItem Name           ")
	for playerId := 0; playerId < matrix.numPlayers; playerId++ {
		builder.WriteString(fmt.Sprintf(" P%1d", playerId))
		if int(currentPlayerId) == playerId {
			builder.WriteString("* ")
		} else {
			builder.WriteString("  ")
		}
	}
	fmt.Println(builder.String())
}

func (deductions *Deductions) Print() {
	matrix := deductions.matrix
	matrix.printMatrixHeader(deductions.thisPlayerId)
	var builder strings.Builder
	for _, itemId := range GetItemsOfType(AllItems) {
		if deductions.itemIsEliminated[itemId] {
			const namelen = 20
			eliminated := " (x)"
			builder.WriteString(GetItemName(itemId))
			builder.WriteString(eliminated)
			remainder := namelen - len(GetItemName(itemId)) - len(eliminated)
			for r := 0; r < remainder; r++ {
				builder.WriteByte(' ')
			}
		} else {
			builder.WriteString(fmt.Sprintf("%-20s", GetItemName(itemId)))
		}

		for playerId := 0; playerId < matrix.numPlayers; playerId++ {
			switch matrix.hasItem[playerId][itemId] {
			case Yes:
				builder.WriteString(" 00  ")
			case No:
				builder.WriteString(" XX  ")
			case Maybe:
				builder.WriteString(" ??  ")
			case Unknown:
				builder.WriteString(" --  ")
			}
		}
		fmt.Println(builder.String())
		builder.Reset()
	}
}

func (matrix *playerItemMatrix) otherPlayers(thisPlayerId PlayerId) ([]PlayerId) {
	var others []PlayerId
	for playerId := 0; playerId < matrix.numPlayers; playerId++ {
		if PlayerId(playerId) != thisPlayerId {
			others = append(others, PlayerId(playerId))
		}
	}
	return others
}
