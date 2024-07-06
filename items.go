package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

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

var humanReadableNames map[string]string

func validHumanReadableItem(name string) bool {
	return strings.HasPrefix(name, "T") &&
		(strings.Contains(name, "_MAIN_") ||
			strings.Contains(name, "_2H_") ||
			strings.Contains(name, "_OFF_") ||
			strings.Contains(name, "_HEAD_") ||
			strings.Contains(name, "_ARMOR_") ||
			strings.Contains(name, "_SHOES_") ||
			strings.Contains(name, "_CAPEITEM_") ||
			strings.Contains(name, "_POTION_") ||
			strings.Contains(name, "_MEAL_") ||
			strings.Contains(name, "_MOUNT_") ||
			strings.Contains(name, "_BAG") ||
			strings.Contains(name, "_BACKPACK_"))
}

func sanitizeItemName(name string) string {
	item, err := typeStringToItem(name, 0)
	if err != nil {
		log.Error("Failed to sanitize item name: ", name)
		return name
	}
	return item.Name
}

func sanitizeHumanReadableItemName(humanReadableName string) string {
	humanReadableName = strings.Replace(humanReadableName, "Beginner's ", "", 1)
	humanReadableName = strings.Replace(humanReadableName, "Novice's ", "", 1)
	humanReadableName = strings.Replace(humanReadableName, "Journeyman's ", "", 1)
	humanReadableName = strings.Replace(humanReadableName, "Adept's ", "", 1)
	humanReadableName = strings.Replace(humanReadableName, "Expert's ", "", 1)
	humanReadableName = strings.Replace(humanReadableName, "Master's ", "", 1)
	humanReadableName = strings.Replace(humanReadableName, "Grandmaster's ", "", 1)
	humanReadableName = strings.Replace(humanReadableName, "Elder's ", "", 1)
	return humanReadableName
}

func updateHumanReadable() {
	if humanReadableNames == nil {
		humanReadableNames = make(map[string]string)
	}

	// Make HTTP GET request
	resp, err := http.Get(config.ItemNamesUrl)
	if err != nil {
		fmt.Printf("Error fetching URL: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	text := string(body)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.SplitN(line, ":", 3)
		if len(parts) != 3 {
			continue
		}
		key := strings.TrimSpace(parts[1])
		value := strings.TrimSpace(parts[2])
		if validHumanReadableItem(key) {
			humanReadableNames[sanitizeItemName(key)] = sanitizeHumanReadableItemName(value)
		}
	}

}

func toHumanReadable(item Item, updateAllowed bool) (string, bool, error) {
	updated := false
	if humanReadableNames == nil {
		updated = true
		updateHumanReadable()
	}
	if humanReadableNames[item.Name] == "" && updateAllowed {
		updated = true
		updateHumanReadable()
	}
	if humanReadableNames[item.Name] == "" {
		log.Error("Failed to get human readable name for item: ", item)
		return item.Name, updated, fmt.Errorf("failed to get human readable name for item: %s", item.Name)
	}
	return humanReadableNames[item.Name], updated, nil
}

func manyToHumanReadable(items []Item) (map[string]string, error) {
	translation := make(map[string]string)
	var errs []error
	updated := false

	for _, item := range items {
		humanReadable, didUpdate, err := toHumanReadable(item, !updated)
		updated = updated || didUpdate
		if err != nil {
			log.Error("Failed to get human readable name for item: ", item)
			errs = append(errs, err)
			translation[item.Name] = item.Name
		} else {
			translation[item.Name] = humanReadable
		}
	}

	if len(errs) != 0 {
		log.Error("Encountered errors while fetching human readable names for items: ", items)
		return translation, fmt.Errorf("%v", errs)
	}
	return translation, nil
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
