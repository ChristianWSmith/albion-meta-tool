package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

type Event struct {
	EventId                   int64
	KillerMainHandName        string
	KillerMainHandTier        uint8
	KillerMainHandEnchantment uint8
	KillerMainHandQuality     uint8
	KillerOffHandName         string
	KillerOffHandTier         uint8
	KillerOffHandEnchantment  uint8
	KillerOffHandQuality      uint8
	KillerHeadName            string
	KillerHeadTier            uint8
	KillerHeadEnchantment     uint8
	KillerHeadQuality         uint8
	KillerChestName           string
	KillerChestTier           uint8
	KillerChestEnchantment    uint8
	KillerChestQuality        uint8
	KillerFootName            string
	KillerFootTier            uint8
	KillerFootEnchantment     uint8
	KillerFootQuality         uint8
	KillerCapeName            string
	KillerCapeTier            uint8
	KillerCapeEnchantment     uint8
	KillerCapeQuality         uint8
	KillerPotionName          string
	KillerPotionTier          uint8
	KillerPotionEnchantment   uint8
	KillerFoodName            string
	KillerFoodTier            uint8
	KillerFoodEnchantment     uint8
	KillerMountName           string
	KillerMountTier           uint8
	KillerMountEnchantment    uint8
	KillerMountQuality        uint8
	KillerBagName             string
	KillerBagTier             uint8
	KillerBagEnchantment      uint8
	KillerBagQuality          uint8
	KillerAverageIp           float64
	VictimMainHandName        string
	VictimMainHandTier        uint8
	VictimMainHandEnchantment uint8
	VictimMainHandQuality     uint8
	VictimOffHandName         string
	VictimOffHandTier         uint8
	VictimOffHandEnchantment  uint8
	VictimOffHandQuality      uint8
	VictimHeadName            string
	VictimHeadTier            uint8
	VictimHeadEnchantment     uint8
	VictimHeadQuality         uint8
	VictimChestName           string
	VictimChestTier           uint8
	VictimChestEnchantment    uint8
	VictimChestQuality        uint8
	VictimFootName            string
	VictimFootTier            uint8
	VictimFootEnchantment     uint8
	VictimFootQuality         uint8
	VictimCapeName            string
	VictimCapeTier            uint8
	VictimCapeEnchantment     uint8
	VictimCapeQuality         uint8
	VictimPotionName          string
	VictimPotionTier          uint8
	VictimPotionEnchantment   uint8
	VictimFoodName            string
	VictimFoodTier            uint8
	VictimFoodEnchantment     uint8
	VictimMountName           string
	VictimMountTier           uint8
	VictimMountEnchantment    uint8
	VictimMountQuality        uint8
	VictimBagName             string
	VictimBagTier             uint8
	VictimBagEnchantment      uint8
	VictimBagQuality          uint8
	VictimAverageIp           float64
	NumberOfParticipants      uint8
	Timestamp                 time.Time
}

func getKillEventUrls() []string {
	var urls []string
	for offset := 0; offset <= 1000; offset += 50 {
		urls = append(urls, fmt.Sprintf("%s?limit=51&offset=%v", config.KillEventUrl, offset))
	}
	return urls
}

func eventMonitor() {
	// Make the HTTP GET request
	log.Info("Kill event urls: ", getKillEventUrls())

	response, err := http.Get(config.KillEventUrl)
	if err != nil {
		log.Warn("The HTTP request failed with error ", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error("Failed to read the response body: ", err)
		return
	}

	// Use Gjson to parse and query the JSON response
	json := string(body)

	if !gjson.Valid(json) {
		log.Error("Invalid json")
	}
	// Example: Iterate over all events and print the Killer's Name
	gjson.Parse(json).ForEach(func(key, value gjson.Result) bool {
		killerName := value.Get("Killer.Name").String()
		log.Debug("Event: ", key.String(), ", Killer's Name: ", killerName)
		return true // keep iterating
	})

}
