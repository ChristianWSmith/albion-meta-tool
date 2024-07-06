package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"net/http"
)

type ItemStats struct {
	SilverGained float64
	SilverLost   float64
	Kills        int64
	Deaths       int64
}

func generateReport() ([][]string, error) {
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
				stats := itemsToStats[event.KillerBuild.MainHand]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				itemsToStats[event.KillerBuild.MainHand] = stats
			}
			if event.VictimBuild.MainHand.Name != "" {
				stats := itemsToStats[event.VictimBuild.MainHand]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				itemsToStats[event.VictimBuild.MainHand] = stats
			}
		}
		if buildFilter.OffHand {
			if event.KillerBuild.OffHand.Name != "" {
				stats := itemsToStats[event.KillerBuild.OffHand]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				itemsToStats[event.KillerBuild.OffHand] = stats
			}
			if event.VictimBuild.OffHand.Name != "" {
				stats := itemsToStats[event.VictimBuild.OffHand]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				itemsToStats[event.VictimBuild.OffHand] = stats
			}
		}
		if buildFilter.Head {
			if event.KillerBuild.Head.Name != "" {
				stats := itemsToStats[event.KillerBuild.Head]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				itemsToStats[event.KillerBuild.Head] = stats
			}
			if event.VictimBuild.Head.Name != "" {
				stats := itemsToStats[event.VictimBuild.Head]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				itemsToStats[event.VictimBuild.Head] = stats
			}
		}
		if buildFilter.Chest {
			if event.KillerBuild.Chest.Name != "" {
				stats := itemsToStats[event.KillerBuild.Chest]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				itemsToStats[event.KillerBuild.Chest] = stats
			}
			if event.VictimBuild.Chest.Name != "" {
				stats := itemsToStats[event.VictimBuild.Chest]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				itemsToStats[event.VictimBuild.Chest] = stats
			}
		}
		if buildFilter.Foot {
			if event.KillerBuild.Foot.Name != "" {
				stats := itemsToStats[event.KillerBuild.Foot]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				itemsToStats[event.KillerBuild.Foot] = stats
			}
			if event.VictimBuild.Foot.Name != "" {
				stats := itemsToStats[event.VictimBuild.Foot]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				itemsToStats[event.VictimBuild.Foot] = stats
			}
		}
		if buildFilter.Cape {
			if event.KillerBuild.Cape.Name != "" {
				stats := itemsToStats[event.KillerBuild.Cape]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				itemsToStats[event.KillerBuild.Cape] = stats
			}
			if event.VictimBuild.Cape.Name != "" {
				stats := itemsToStats[event.VictimBuild.Cape]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				itemsToStats[event.VictimBuild.Cape] = stats
			}
		}
		if buildFilter.Food {
			if event.KillerBuild.Food.Name != "" {
				stats := itemsToStats[event.KillerBuild.Food]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				itemsToStats[event.KillerBuild.Food] = stats
			}
			if event.VictimBuild.Food.Name != "" {
				stats := itemsToStats[event.VictimBuild.Food]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				itemsToStats[event.VictimBuild.Food] = stats
			}
		}
		if buildFilter.Potion {
			if event.KillerBuild.Potion.Name != "" {
				stats := itemsToStats[event.KillerBuild.Potion]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				itemsToStats[event.KillerBuild.Potion] = stats
			}
			if event.VictimBuild.Potion.Name != "" {
				stats := itemsToStats[event.VictimBuild.Potion]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				itemsToStats[event.VictimBuild.Potion] = stats
			}
		}
		if buildFilter.Mount {
			if event.KillerBuild.Mount.Name != "" {
				stats := itemsToStats[event.KillerBuild.Mount]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				itemsToStats[event.KillerBuild.Mount] = stats
			}
			if event.VictimBuild.Mount.Name != "" {
				stats := itemsToStats[event.VictimBuild.Mount]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				itemsToStats[event.VictimBuild.Mount] = stats
			}
		}
		if buildFilter.Bag {
			if event.KillerBuild.Bag.Name != "" {
				stats := itemsToStats[event.KillerBuild.Bag]
				stats.SilverGained += buildPrices[event.VictimBuild]
				stats.Kills += 1
				itemsToStats[event.KillerBuild.Bag] = stats
			}
			if event.VictimBuild.Bag.Name != "" {
				stats := itemsToStats[event.VictimBuild.Bag]
				stats.SilverLost += buildPrices[event.VictimBuild]
				stats.Deaths += 1
				itemsToStats[event.VictimBuild.Bag] = stats
			}
		}
	}

	// format to csv
	response = append(response, []string{"item_id", "tier", "enchantment", "equivalence", "usages", "silver_ratio", "kills", "deaths", "silver_gained", "silver_lost"})
	for item, stats := range itemsToStats {
		response = append(response, []string{
			item.Name,
			fmt.Sprintf("%d", item.Tier),
			fmt.Sprintf("%d", item.Enchantment),
			fmt.Sprintf("%d", item.Tier+item.Enchantment),
			fmt.Sprintf("%d", stats.Kills+stats.Deaths),
			fmt.Sprintf("%f", stats.SilverGained/math.Max(stats.SilverLost, 1.0)),
			fmt.Sprintf("%d", stats.Kills),
			fmt.Sprintf("%d", stats.Deaths),
			fmt.Sprintf("%f", stats.SilverGained),
			fmt.Sprintf("%f", stats.SilverLost),
		})
	}

	return response, nil
}

// Handler function for the endpoint
func reportHandler(w http.ResponseWriter, r *http.Request) {
	// Create a response in CSV format
	response, _ := generateReport()

	// Encode and send the CSV response
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=report.csv")
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

func startAPI() {
	http.HandleFunc("/report", reportHandler)
	log.Info("Server starting on port ", config.Port, "...")
	log.Error(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
