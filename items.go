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

type BuildFilter struct {
	MainHand bool
	OffHand  bool
	Head     bool
	Chest    bool
	Foot     bool
	Cape     bool
	Potion   bool
	Food     bool
	Mount    bool
	Bag      bool
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

func getItemsFromBuilds(builds []Build, filter BuildFilter) []Item {
	var items []Item
	for _, build := range builds {
		if filter.MainHand && build.MainHand.Name != "" {
			items = append(items, build.MainHand)
		}
		if filter.OffHand && build.OffHand.Name != "" {
			items = append(items, build.OffHand)
		}
		if filter.Head && build.Head.Name != "" {
			items = append(items, build.Head)
		}
		if filter.Chest && build.Chest.Name != "" {
			items = append(items, build.Chest)
		}
		if filter.Foot && build.Foot.Name != "" {
			items = append(items, build.Foot)
		}
		if filter.Cape && build.Cape.Name != "" {
			items = append(items, build.Cape)
		}
		if filter.Potion && build.Potion.Name != "" {
			items = append(items, build.Potion)
		}
		if filter.Food && build.Food.Name != "" {
			items = append(items, build.Food)
		}
		if filter.Mount && build.Mount.Name != "" {
			items = append(items, build.Mount)
		}
		if filter.Bag && build.Bag.Name != "" {
			items = append(items, build.Bag)
		}
	}
	return items
}
