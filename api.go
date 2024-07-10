package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
)

type ItemStats struct {
	SilverGained float64
	SilverLost   float64
	Kills        int64
	Deaths       int64
	SumAverageIp float64
}

type BuildStats struct {
	SilverGained float64
	SilverLost   float64
	Kills        int64
	Deaths       int64
	SumAverageIp float64
	Prices       []float64
}

func generateItemReport() ([][]string, error) {
	response := [][]string{}

	// get all events
	events, err := queryAllEvents()
	if err != nil {
		log.Error("Failed to query all events: ", err)
		return response, err
	}

	// filter for solo events
	var solo_events []Event
	for _, event := range events {
		if event.NumberOfParticipants == 1 {
			solo_events = append(solo_events, event)
		}
	}
	events = solo_events

	// get all builds
	var builds []Build
	for _, event := range events {
		builds = append(builds, event.KillerBuild, event.VictimBuild)
	}

	// get all items
	var buildFilter = BuildFilter{
		MainHand: true,
		OffHand:  true,
		Head:     true,
		Chest:    true,
		Foot:     true,
		Cape:     true,
		Potion:   false,
		Food:     false,
		Mount:    false,
		Bag:      false,
	}
	items := getItemsFromBuilds(builds, buildFilter)

	// get their prices
	itemPrices, _ := getItemPrices(items)

	// get build prices
	buildPrices := getBuildPrices(builds, itemPrices, buildFilter)

	//
	itemsToStats := make(map[Item]ItemStats)
	for _, event := range events {
		if buildPrices[event.VictimBuild] == 0.0 {
			continue
		}
		if buildFilter.MainHand {
			if event.KillerBuild.MainHand.Name != "" {
				item := event.KillerBuild.MainHand
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				stats.SumAverageIp += event.KillerAverageIp
				itemsToStats[item] = stats
			}
			if event.VictimBuild.MainHand.Name != "" {
				item := event.VictimBuild.MainHand
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				stats.SumAverageIp += event.VictimAverageIp
				itemsToStats[item] = stats
			}
		}
		if buildFilter.OffHand {
			if event.KillerBuild.OffHand.Name != "" {
				item := event.KillerBuild.OffHand
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				stats.SumAverageIp += event.KillerAverageIp
				itemsToStats[item] = stats
			}
			if event.VictimBuild.OffHand.Name != "" {
				item := event.VictimBuild.OffHand
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				stats.SumAverageIp += event.VictimAverageIp
				itemsToStats[item] = stats
			}
		}
		if buildFilter.Head {
			if event.KillerBuild.Head.Name != "" {
				item := event.KillerBuild.Head
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				stats.SumAverageIp += event.KillerAverageIp
				itemsToStats[item] = stats
			}
			if event.VictimBuild.Head.Name != "" {
				item := event.VictimBuild.Head
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				stats.SumAverageIp += event.VictimAverageIp
				itemsToStats[item] = stats
			}
		}
		if buildFilter.Chest {
			if event.KillerBuild.Chest.Name != "" {
				item := event.KillerBuild.Chest
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				stats.SumAverageIp += event.KillerAverageIp
				itemsToStats[item] = stats
			}
			if event.VictimBuild.Chest.Name != "" {
				item := event.VictimBuild.Chest
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				stats.SumAverageIp += event.VictimAverageIp
				itemsToStats[item] = stats
			}
		}
		if buildFilter.Foot {
			if event.KillerBuild.Foot.Name != "" {
				item := event.KillerBuild.Foot
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				stats.SumAverageIp += event.KillerAverageIp
				itemsToStats[item] = stats
			}
			if event.VictimBuild.Foot.Name != "" {
				item := event.VictimBuild.Foot
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				stats.SumAverageIp += event.VictimAverageIp
				itemsToStats[item] = stats
			}
		}
		if buildFilter.Cape {
			if event.KillerBuild.Cape.Name != "" {
				item := event.KillerBuild.Cape
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				stats.SumAverageIp += event.KillerAverageIp
				itemsToStats[item] = stats
			}
			if event.VictimBuild.Cape.Name != "" {
				item := event.VictimBuild.Cape
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				stats.SumAverageIp += event.VictimAverageIp
				itemsToStats[item] = stats
			}
		}
		if buildFilter.Potion {
			if event.KillerBuild.Potion.Name != "" {
				item := event.KillerBuild.Potion
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				stats.SumAverageIp += event.KillerAverageIp
				itemsToStats[item] = stats
			}
			if event.VictimBuild.Potion.Name != "" {
				item := event.VictimBuild.Potion
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				stats.SumAverageIp += event.VictimAverageIp
				itemsToStats[item] = stats
			}
		}
		if buildFilter.Food {
			if event.KillerBuild.Food.Name != "" {
				item := event.KillerBuild.Food
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				stats.SumAverageIp += event.KillerAverageIp
				itemsToStats[item] = stats
			}
			if event.VictimBuild.Food.Name != "" {
				item := event.VictimBuild.Food
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				stats.SumAverageIp += event.VictimAverageIp
				itemsToStats[item] = stats
			}
		}
		if buildFilter.Mount {
			if event.KillerBuild.Mount.Name != "" {
				item := event.KillerBuild.Mount
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				stats.SumAverageIp += event.KillerAverageIp
				itemsToStats[item] = stats
			}
			if event.VictimBuild.Mount.Name != "" {
				item := event.VictimBuild.Mount
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				stats.SumAverageIp += event.VictimAverageIp
				itemsToStats[item] = stats
			}
		}
		if buildFilter.Bag {
			if event.KillerBuild.Bag.Name != "" {
				item := event.KillerBuild.Bag
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				stats.SumAverageIp += event.KillerAverageIp
				itemsToStats[item] = stats
			}
			if event.VictimBuild.Bag.Name != "" {
				item := event.VictimBuild.Bag
				item.Quality = 0
				stats := itemsToStats[item]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				stats.SumAverageIp += event.VictimAverageIp
				itemsToStats[item] = stats
			}

		}

	}

	// get human readable
	var itemsThatHaveStats []Item
	for item := range itemsToStats {
		itemsThatHaveStats = append(itemsThatHaveStats, item)
	}
	humanReadableNamesBatch, err := manyToHumanReadable(itemsThatHaveStats)
	if err != nil {
		log.Error("Failed to fetch some human readable names during report generation: ", err)
	}

	// format to csv
	response = append(response, []string{
		"item_id",
		"tier",
		"enchantment",
		"equivalence",
		"usages",
		"average_build_ip",
		"k/d",
		"silver_ratio",
		"kills",
		"deaths",
		"silver_gained",
		"silver_lost",
	})
	for item, stats := range itemsToStats {
		normalItem := item
		goodItem := item
		outstandingItem := item
		excellentItem := item
		masterpieceItem := item
		normalItem.Quality = 0
		goodItem.Quality = 1
		outstandingItem.Quality = 2
		excellentItem.Quality = 3
		masterpieceItem.Quality = 4
		var qualityPrices []float64
		if itemPrices[normalItem] != 0.0 {
			qualityPrices = append(qualityPrices, itemPrices[normalItem])
		}
		if itemPrices[goodItem] != 0.0 {
			qualityPrices = append(qualityPrices, itemPrices[goodItem])
		}
		if itemPrices[outstandingItem] != 0.0 {
			qualityPrices = append(qualityPrices, itemPrices[outstandingItem])
		}
		if itemPrices[excellentItem] != 0.0 {
			qualityPrices = append(qualityPrices, itemPrices[excellentItem])
		}
		if itemPrices[masterpieceItem] != 0.0 {
			qualityPrices = append(qualityPrices, itemPrices[masterpieceItem])
		}
		medianPrice := calculateMedian(qualityPrices)
		if medianPrice == 0.0 {
			medianPrice = math.Inf(1)
		}
		response = append(response, []string{
			humanReadableNamesBatch[item.Name],
			fmt.Sprintf("%d", item.Tier),
			fmt.Sprintf("%d", item.Enchantment),
			fmt.Sprintf("%d", item.Tier+item.Enchantment),
			fmt.Sprintf("%d", stats.Kills+stats.Deaths),
			fmt.Sprintf("%f", stats.SumAverageIp/float64(stats.Kills+stats.Deaths)),
			fmt.Sprintf("%f", float64(stats.Kills)/math.Max(float64(stats.Deaths), 1.0)),
			fmt.Sprintf("%f", stats.SilverGained*0.7/math.Max(math.Max(stats.SilverLost, 1.0), medianPrice)),
			fmt.Sprintf("%d", stats.Kills),
			fmt.Sprintf("%d", stats.Deaths),
			fmt.Sprintf("%f", stats.SilverGained*0.7),
			fmt.Sprintf("%f", stats.SilverLost),
		})
	}

	return response, nil
}

// Handler function for the endpoint
func itemReportHandler(w http.ResponseWriter, _ *http.Request) {
	// Create a response in CSV format
	response, _ := generateItemReport()

	// Encode and send the CSV response
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=report-%d.csv", time.Now().Unix()))
	writer := csv.NewWriter(w)
	for _, record := range response {
		if err := writer.Write(record); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func generateBuildReport() ([][]string, error) {
	response := [][]string{}

	// get all events
	events, err := queryAllEvents()
	if err != nil {
		log.Error("Failed to query all events: ", err)
		return response, err
	}

	// filter for solo events
	var solo_events []Event
	for _, event := range events {
		if event.NumberOfParticipants == 1 {
			solo_events = append(solo_events, event)
		}
	}
	events = solo_events

	// get all builds
	var builds []Build
	for _, event := range events {
		builds = append(builds, event.KillerBuild, event.VictimBuild)
	}

	// get all items
	var buildFilter = BuildFilter{
		MainHand: true,
		OffHand:  true,
		Head:     true,
		Chest:    true,
		Foot:     true,
		Cape:     true,
		Potion:   false,
		Food:     false,
		Mount:    false,
		Bag:      false,
	}
	items := getItemsFromBuilds(builds, buildFilter)

	// get their prices
	itemPrices, _ := getItemPrices(items)

	// get build prices
	buildPrices := getBuildPrices(builds, itemPrices, buildFilter)

	//
	buildsNamesOnlyToStats := make(map[BuildNamesOnly]BuildStats)
	for _, event := range events {
		if buildPrices[event.VictimBuild] == 0.0 {
			continue
		}
		killerBuildNamesOnly := buildToNamesOnly(event.KillerBuild, buildFilter)
		victimBuildNamesOnly := buildToNamesOnly(event.VictimBuild, buildFilter)

		killerBuildStats := buildsNamesOnlyToStats[killerBuildNamesOnly]
		killerBuildStats.Kills += 1
		killerBuildStats.SilverGained += buildPrices[event.VictimBuild]
		killerBuildStats.SumAverageIp += event.KillerAverageIp
		if buildPrices[event.KillerBuild] != 0.0 {
			killerBuildStats.Prices = append(killerBuildStats.Prices, buildPrices[event.KillerBuild])
		}
		buildsNamesOnlyToStats[killerBuildNamesOnly] = killerBuildStats

		victimBuildStats := buildsNamesOnlyToStats[victimBuildNamesOnly]
		victimBuildStats.Deaths += 1
		victimBuildStats.SilverLost += buildPrices[event.VictimBuild]
		victimBuildStats.SumAverageIp += event.VictimAverageIp
		if buildPrices[event.VictimBuild] != 0.0 {
			victimBuildStats.Prices = append(victimBuildStats.Prices, buildPrices[event.VictimBuild])
		}
		buildsNamesOnlyToStats[victimBuildNamesOnly] = victimBuildStats
	}

	// get human readable
	var batchBuildNamesOnly []BuildNamesOnly
	for buildNameOnly, _ := range buildsNamesOnlyToStats {
		batchBuildNamesOnly = append(batchBuildNamesOnly, buildNameOnly)
	}
	var itemsInBuilds = namesOnlyToItems(batchBuildNamesOnly, buildFilter)
	humanReadableNamesBatch, err := manyToHumanReadable(itemsInBuilds)
	if err != nil {
		log.Error("Failed to fetch some human readable names during report generation: ", err)
	}

	// format to csv
	response = append(response, []string{
		"main_hand",
		"off_hand",
		"head",
		"chest",
		"foot",
		"cape",
		"food",
		"potion",
		"mount",
		"bag",
		"usages",
		"average_ip",
		"k/d",
		"silver_ratio",
		"kills",
		"deaths",
		"silver_gained",
		"silver_lost",
	})
	for buildNamesOnly, stats := range buildsNamesOnlyToStats {
		price := calculateMedian(stats.Prices)
		if price == 0.0 {
			price = math.Inf(1)
		}
		response = append(response, []string{
			humanReadableNamesBatch[buildNamesOnly.MainHand],
			humanReadableNamesBatch[buildNamesOnly.OffHand],
			humanReadableNamesBatch[buildNamesOnly.Head],
			humanReadableNamesBatch[buildNamesOnly.Chest],
			humanReadableNamesBatch[buildNamesOnly.Foot],
			humanReadableNamesBatch[buildNamesOnly.Cape],
			humanReadableNamesBatch[buildNamesOnly.Food],
			humanReadableNamesBatch[buildNamesOnly.Potion],
			humanReadableNamesBatch[buildNamesOnly.Mount],
			humanReadableNamesBatch[buildNamesOnly.Bag],
			fmt.Sprintf("%d", stats.Kills+stats.Deaths),
			fmt.Sprintf("%f", stats.SumAverageIp/float64(stats.Kills+stats.Deaths)),
			fmt.Sprintf("%f", float64(stats.Kills)/math.Max(float64(stats.Deaths), 1.0)),
			fmt.Sprintf("%f", stats.SilverGained/math.Max(math.Max(stats.SilverLost, 1.0), price)),
			fmt.Sprintf("%d", stats.Kills),
			fmt.Sprintf("%d", stats.Deaths),
			fmt.Sprintf("%f", stats.SilverGained),
			fmt.Sprintf("%f", stats.SilverLost),
		})
	}

	return response, nil
}

// Handler function for the endpoint
func buildReportHandler(w http.ResponseWriter, _ *http.Request) {
	// Create a response in CSV format
	response, _ := generateBuildReport()

	// Encode and send the CSV response
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=report-%d.csv", time.Now().Unix()))
	writer := csv.NewWriter(w)
	for _, record := range response {
		if err := writer.Write(record); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type StatsResponse struct {
	NumEvents int           `json:"num_events"`
	NumPrices int           `json:"num_prices"`
	SleepTime time.Duration `json:"sleep_time"`
}

func statsHandler(w http.ResponseWriter, _ *http.Request) {
	// Create sample data

	numEvents, err := getNumEvents()
	if err != nil {
		log.Error("Failed to get number of events during API call")
	}
	numPrices, err := getNumPrices()
	if err != nil {
		log.Error("Failed to get number of prices during API call")
	}

	responseData := StatsResponse{
		NumEvents: numEvents,
		NumPrices: numPrices,
		SleepTime: sleepTime,
	}

	// Marshal responseData into JSON format
	responseJSON, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write JSON response
	w.Write(responseJSON)
}

func startAPI() {
	http.HandleFunc("/itemReport", itemReportHandler)
	http.HandleFunc("/buildReport", buildReportHandler)
	http.HandleFunc("/stats", statsHandler)
	log.Info("Server starting on port ", config.Port, "...")
	log.Error(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
