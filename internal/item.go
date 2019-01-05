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
	AllItems
)

var ItemTypes = []ItemType{
	Character,
	Weapon,
	Room,
}

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

type Item struct {
	Id   ItemId
	Type ItemType
	Name string
}

type ItemIdList []ItemId

var itemDetails = []Item{
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
var allItems = selectItemsOfType(AllItems)

func GetItemsOfType(itemType ItemType) (ItemIdList) {
	var empty ItemIdList
	switch itemType {
	case Character:
		return characters
	case Weapon:
		return weapons
	case Room:
		return rooms
	case AllItems:
		return allItems
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

func GetItemType(itemId ItemId) (ItemType) {
	return getItem(itemId).Type
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
	return itemDetails[randItemIndex].Id
}

func GetRelevantItem(itemType ItemType, eliminated ItemIdList) (ItemId) {
	itemsOfType := GetItemsOfType(itemType)
	for _, itemId := range eliminated {
		itemsOfType = removeItem(itemsOfType,itemId)
	}
	randItemIndex := itemsOfType[rand.Intn(len(itemsOfType))]
	return itemDetails[randItemIndex].Id
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
	for _, item := range itemDetails {
		if item.Id == itemId {
			return &item
		}
	}
	return nil
}

func getShuffledItemIdList() (ItemIdList) {
	shuffled := make(ItemIdList, len(itemDetails))
	for i := 0; i < len(itemDetails); i++ {
		shuffled[i] = itemDetails[i].Id
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
	for _, item := range itemDetails {
		if itemType == item.Type || itemType == AllItems {
			itemsOfType = append(itemsOfType, item.Id)
		}
	}
	return itemsOfType
}
