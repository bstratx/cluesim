package internal

import "fmt"

type PlayerId int

type Player struct {
	Id          PlayerId
	Name        string
	CharacterId ItemId
	ItemIds     ItemIdList
	//deductions  Deductions
}

func (player Player) PoseQuestion() (*Question, error) {
	character := GetRandomItem(Character)
	weapon := GetRandomItem(Weapon)
	room := GetRandomItem(Room)
	return createQuestion(player.Id, ItemIdList{character, weapon, room})
}

func (player Player) RespondToQuestion(question *Question) (*ItemId) {
	for _, playerItemId := range player.ItemIds {
		for _, questionItemId := range question.ItemIds {
			if playerItemId == questionItemId {
				return &playerItemId
			}
		}
	}
	return nil
}

func (player Player) Print() {
	header := fmt.Sprintf("%s holds ", player.Name)
	player.ItemIds.PrintItems(header, "\n")
}
