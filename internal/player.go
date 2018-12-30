package internal

import "fmt"

type PlayerId int

type Player struct {
	Id          PlayerId
	Name        string
	CharacterId ItemId
	ItemIds     []ItemId
}

func (player Player) PoseQuestion() (*Question, error) {
	character := GetRandomItem(Character).Id
	weapon := GetRandomItem(Weapon).Id
	room := GetRandomItem(Room).Id
	return createQuestion(player.Id, []ItemId{character, weapon, room})
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
	fmt.Printf("%s holds ", player.Name)
	PrintItems(player.ItemIds)
	fmt.Println()
}
