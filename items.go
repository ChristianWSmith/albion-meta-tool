package main

import "fmt"

type Item struct {
	Name        string
	Tier        uint8
	Enchantment uint8
	Quality     uint8
}

type Build struct {
	MainHand Item
	OffHand  Item
	Head     Item
	Chest    Item
	Foot     Item
	Cape     Item
	Potion   Item
	Food     Item
	Mount    Item
	Bag      Item
}

func typeStringToItem(typeString string, quality uint8) (Item, error) {
	var item Item
	var err error

	if typeString == "" {
		return item, nil
	}

	item.Tier, err = parseUint8(fmt.Sprintf("%c", typeString[1]))
	if err != nil {
		log.Error("Failed to parse item string tier: ", typeString)
		return item, err
	}
	typeString = typeString[3:]

	if typeString[len(typeString)-2] == '@' {
		item.Enchantment, err = parseUint8(fmt.Sprintf("%c", typeString[len(typeString)-1]))
		if err != nil {
			log.Error("Failed to parse item string enchantment: ", typeString)
			return item, err
		}
		typeString = typeString[:len(typeString)-2]
	} else {
		item.Enchantment = 0
	}

	item.Name = typeString

	item.Quality = quality

	return item, nil
}

func itemToTypeString(item Item) string {
	if item.Enchantment == 0 {
		return fmt.Sprintf("T%d_%s", item.Tier, item.Name)
	}
	return fmt.Sprintf("T%d_%s@%d", item.Tier, item.Name, item.Enchantment)
}
