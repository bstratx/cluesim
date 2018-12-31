package internal

import "fmt"

type PlayerId int

type Player struct {
	Id          PlayerId
	Name        string
	CharacterId ItemId
	ItemIds     ItemIdList
	Deductions  *Deductions
}

func NewPlayer(playerId PlayerId, name string, character ItemId, items ItemIdList, numPlayers int) (Player) {
	deductions := NewDeductions(playerId, numPlayers)
	player := Player{playerId, name, character, items, deductions}
	player.Deductions.RecordPlayerItems(player.ItemIds)
	return player
}

func (player Player) PoseQuestion() (*Question, error) {
	eliminatedIds := player.Deductions.GetEliminatedItems()
	character := GetRelevantItem(Character, eliminatedIds)
	weapon := GetRelevantItem(Weapon, eliminatedIds)
	room := GetRelevantItem(Room, eliminatedIds)
	return createQuestion(player.Id, ItemIdList{character, weapon, room})
}

func (player Player) RespondToQuestion(question *Question) (*ItemId) {
	var matches ItemIdList
	for _, playerItemId := range player.ItemIds {
		for _, questionItemId := range question.ItemIds {
			if playerItemId == questionItemId {
				matches = append(matches,playerItemId)
			}
		}
	}

	// attempt to reveal least information
	if len(matches) > 0 {
		for _, matchItemId := range matches {
			if GetItemType(matchItemId) == Room {
				return &matchItemId
			}
		}
		return &matches[0]
	}
	return nil
}

func (player Player) Print() {
	header := fmt.Sprintf("%s holds ", player.Name)
	player.ItemIds.PrintItems(header, "\n")
}
