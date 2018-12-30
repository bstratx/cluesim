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

func selectItemsOfType(itemType ItemType) ([]ItemId) {
	var itemsOfType []ItemId
	for _, item := range allItems {
		if item.Type == itemType {
			itemsOfType = append(itemsOfType, item.Id)
		}
	}
	return itemsOfType
}

func GetItemsOfType(itemType ItemType) ([]ItemId) {
	switch itemType {
	case Character:
		return characters
	case Weapon:
		return weapons
	case Room:
		return rooms
	default:
		return nil
	}
}

func IsItemOfType(itemId ItemId, itemType ItemType) (bool) {
	item := GetItem(itemId)
	if item == nil {
		return false
	}
	return item.Type == itemType
}

func GetItem(itemId ItemId) (*Item) {
	for _, item := range allItems {
		if item.Id == itemId {
			return &item
		}
	}
	return nil
}

func GetRandomItem(itemType ItemType) (*Item) {
	items := GetItemsOfType(itemType)
	if items == nil {
		return nil
	}
	randItemIndex := items[rand.Intn(len(items))]
	return &allItems[randItemIndex]
}

func GetShuffledItemIdList() ([]ItemId) {
	shuffled := make([]ItemId, len(allItems))
	for i := 0; i < len(allItems); i++ {
		shuffled[i] = allItems[i].Id
	}
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

func PrintItems(itemIds []ItemId) {
	for i, itemId := range itemIds {
		fmt.Printf("%s", GetItem(itemId).Name)
		if i+1 != len(itemIds) {
			fmt.Printf(", ")
		}
	}
}
