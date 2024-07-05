package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
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

type Event struct {
	EventId              int64
	KillerBuild          Build
	KillerAverageIp      float64
	VictimBuild          Build
	VictimAverageIp      float64
	NumberOfParticipants uint8
	Timestamp            time.Time
}

var itemStringRegex *regexp.Regexp = regexp.MustCompile(`^T(?P<Tier>\d+)_(?P<Name>[A-Z1-3_]+)(?:@(?P<Enchantment>\d+))?$`)

func getKillEventUrls() []string {
	var urls []string
	for offset := 0; offset <= 1000; offset += 50 {
		urls = append(urls, fmt.Sprintf("%s?limit=51&offset=%v", config.KillEventUrl, offset))
	}
	return urls
}

func parseUint8(s string) (uint8, error) {
	// Parse the string to a uint64 first
	val, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0, err
	}
	// Convert the uint64 to a uint8
	return uint8(val), nil
}

func resultToItem(result gjson.Result) (Item, error) {
	var item Item
	var err error

	item_string := result.Get("Type").String()

	if item_string == "" {
		return item, nil
	}

	item.Tier, err = parseUint8(fmt.Sprintf("%c", item_string[1]))
	if err != nil {
		log.Error("Failed to parse item string tier: ", item_string)
		return item, err
	}
	item_string = item_string[3:]

	if item_string[len(item_string)-2] == '@' {
		item.Enchantment, err = parseUint8(fmt.Sprintf("%c", item_string[len(item_string)-1]))
		if err != nil {
			log.Error("Failed to parse item string enchantment: ", item_string)
			return item, err
		}
		item_string = item_string[:len(item_string)-2]
	} else {
		item.Enchantment = 0
	}

	item.Name = item_string

	item.Quality = uint8(result.Get("Quality").Int())

	return item, nil
}

func resultToBuild(result gjson.Result) (Build, error) {
	var build Build
	var err error

	build.MainHand, err = resultToItem(result.Get("MainHand"))
	if err != nil {
		log.Error("Failed to convert main hand to item", result.Get("MainHand"))
		return build, err

	}

	build.OffHand, err = resultToItem(result.Get("OffHand"))
	if err != nil {
		log.Error("Failed to convert offhand to item", result.Get("OffHand"))
		return build, err

	}

	build.Head, err = resultToItem(result.Get("Head"))
	if err != nil {
		log.Error("Failed to convert head to item", result.Get("Head"))
		return build, err

	}

	build.Chest, err = resultToItem(result.Get("Armor"))
	if err != nil {
		log.Error("Failed to convert chest to item", result.Get("Armor"))
		return build, err

	}

	build.Foot, err = resultToItem(result.Get("Shoes"))
	if err != nil {
		log.Error("Failed to convert foot to item", result.Get("Shoes"))
		return build, err

	}

	build.Bag, err = resultToItem(result.Get("Bag"))
	if err != nil {
		log.Error("Failed to convert bag to item", result.Get("Bag"))
		return build, err

	}

	build.Cape, err = resultToItem(result.Get("Cape"))
	if err != nil {
		log.Error("Failed to convert cape to item", result.Get("Cape"))
		return build, err

	}

	build.Potion, err = resultToItem(result.Get("Potion"))
	if err != nil {
		log.Error("Failed to convert potion to item", result.Get("Potion"))
		return build, err

	}

	build.Food, err = resultToItem(result.Get("Food"))
	if err != nil {
		log.Error("Failed to convert food to item", result.Get("Food"))
		return build, err

	}

	return build, nil
}

func resultToEvent(result gjson.Result) (Event, error) {
	var event Event
	var err error
	event.KillerBuild, err = resultToBuild(result.Get("Killer").Get("Equipment"))
	if err != nil {
		log.Error("Failed to convert killer equipment to build")
		return event, err
	}
	event.VictimBuild, err = resultToBuild(result.Get("Victim").Get("Equipment"))
	if err != nil {
		log.Error("Failed to convert victim equipment to build")
		return event, err
	}
	event.EventId = result.Get("EventId").Int()
	event.KillerAverageIp = result.Get("Killer.AverageItemPower").Float()
	event.VictimAverageIp = result.Get("Victim.AverageItemPower").Float()
	event.NumberOfParticipants = uint8(result.Get("numberOfParticipants").Int())
	event.Timestamp = result.Get("TimeStamp").Time()

	return event, nil
}

func getEvents(url string) ([]Event, error) {
	response, err := http.Get(url)
	var events []Event
	if err != nil {
		log.Warn("The HTTP request failed with error ", err)
		return events, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error("Failed to read the response body: ", err)
		return events, err
	}

	// Use Gjson to parse and query the JSON response
	json := string(body)

	if !gjson.Valid(json) {
		log.Error("Invalid json resonse from url: ", url)
		return events, fmt.Errorf("invalid json resonse from url: %s", url)
	}
	// Example: Iterate over all events and print the Killer's Name
	gjson.Parse(json).ForEach(func(_, result gjson.Result) bool {
		var event, err = resultToEvent(result)
		if err != nil {
			log.Error("Failed to convert result to event", result)
		}
		log.Debug("Parsed event: ", event)
		events = append(events, event)
		return true // keep iterating
	})

	return events, nil
}

func eventMonitor() {
	// Make the HTTP GET request
	log.Info("Kill event urls: ", getKillEventUrls())

	for _, url := range getKillEventUrls() {
		getEvents(url)
	}

}
