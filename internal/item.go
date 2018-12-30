package internal

import (
	"math/rand"
	"fmt"
)

type ItemType int

const (
	Character ItemType = iota
	Weapon
	Room
)

type ItemId int

const (
	Scarlett     ItemId = iota
	Mustard
	White
	Peacock
	Green
	Plum
	Candlestick
	Knife
	LeadPipe
	Revolver
	Rope
	Wrench
	Ballroom
	BilliardRoom
	Conservatory
	DiningRoom
	Hall
	Kitchen
	Library
	Lounge
	Study
)

var ItemTypes = []ItemType{
	Character,
	Weapon,
	Room,
}

type Item struct {
	Id   ItemId
	Type ItemType
	Name string
}

type ItemIdList []ItemId

var allItems = []Item{
	{Scarlett, Character, "Miss Scarlett"},
	{Mustard, Character, "Col. Mustard"},
	{White, Character, "Mrs. White"},
	{Peacock, Character, "Mr. Green"},
	{Green, Character, "Mrs. Peacock"},
	{Plum, Character, "Prof. Plum"},
	{Candlestick, Weapon, "Candlestick"},
	{Knife, Weapon, "Knife"},
	{LeadPipe, Weapon, "Lead Pipe"},
	{Revolver, Weapon, "Revolver"},
	{Rope, Weapon, "Rope"},
	{Wrench, Weapon, "Wrench"},
	{Ballroom, Room, "Ballroom"},
	{BilliardRoom, Room, "Billiard Room"},
	{Conservatory, Room, "Conservatory"},
	{DiningRoom, Room, "Dining Room"},
	{Hall, Room, "Hall"},
	{Kitchen, Room, "Kitchen"},
	{Library, Room, "Library"},
	{Lounge, Room, "Lounge"},
	{Study, Room, "Study"},
}

var characters = selectItemsOfType(Character)
var weapons = selectItemsOfType(Weapon)
var rooms = selectItemsOfType(Room)

func GetItemsOfType(itemType ItemType) (ItemIdList) {
	var empty ItemIdList
	switch itemType {
	case Character:
		return characters
	case Weapon:
		return weapons
	case Room:
		return rooms
	default:
		return empty
	}
}

func IsItemOfType(itemId ItemId, itemType ItemType) (bool) {
	item := getItem(itemId)
	if item == nil {
		return false
	}
	return item.Type == itemType
}

func GetItemName(itemId ItemId) (string) {
	return getItem(itemId).Name
}

func PickSecretItems() (ItemIdList, ItemIdList) {
	var secrets ItemIdList
	remaining := getShuffledItemIdList()
	secrets = append(secrets, GetRandomItem(Character))
	secrets = append(secrets, GetRandomItem(Weapon))
	secrets = append(secrets, GetRandomItem(Room))
	for _, secret := range secrets {
		remaining = removeItem(remaining, secret)
	}
	return secrets, remaining
}

func GetRandomItem(itemType ItemType) (ItemId) {
	items := GetItemsOfType(itemType)
	randItemIndex := items[rand.Intn(len(items))]
	return allItems[randItemIndex].Id
}

func (itemIds ItemIdList) PrintItems(header string, trailer string) {
	fmt.Print(header)
	for i, itemId := range itemIds {
		fmt.Print(getItem(itemId).Name)
		if i+1 != len(itemIds) {
			fmt.Print(", ")
		}
	}
	fmt.Print(trailer)
}

func getItem(itemId ItemId) (*Item) {
	for _, item := range allItems {
		if item.Id == itemId {
			return &item
		}
	}
	return nil
}

func getShuffledItemIdList() (ItemIdList) {
	shuffled := make(ItemIdList, len(allItems))
	for i := 0; i < len(allItems); i++ {
		shuffled[i] = allItems[i].Id
	}
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

func removeItem(itemList ItemIdList, removeItem ItemId) (ItemIdList) {
	var newList ItemIdList
	for i, item := range itemList {
		if item == removeItem {
			newList = append(newList, itemList[0:i] ...)
			newList = append(newList, itemList[i+1:] ...)
			return newList
		}
	}
	return itemList
}

func selectItemsOfType(itemType ItemType) (ItemIdList) {
	var itemsOfType []ItemId
	for _, item := range allItems {
		if item.Type == itemType {
			itemsOfType = append(itemsOfType, item.Id)
		}
	}
	return itemsOfType
}
