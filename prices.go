package main

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/tidwall/gjson"
)

func getItemPrices(items []Item) (map[Item]float64, error) {
	var discoveredPrices map[Item]float64
	var unknownItems []Item
	var err error
	var itemPrices map[Item]float64

	itemPrices, err = queryPrices(items)
	if err != nil {
		log.Error("Failed to query prices for items: ", items)
		return itemPrices, err
	}

	for _, item := range items {
		if itemPrices[item] == 0.0 {
			unknownItems = append(unknownItems, item)
		}
	}

	if len(unknownItems) > 0 {
		log.Debug("Unknown items: ", unknownItems)
		discoveredPrices, err = callPriceAPI(unknownItems)
		if err != nil {
			log.Error("Failed to call price API for prices for: ", unknownItems)
			return itemPrices, err
		}
		err = updatePrices(discoveredPrices)
		if err != nil {
			log.Error("Failed to update prices for items: ", discoveredPrices)
			return itemPrices, err
		}
		for item, price := range discoveredPrices {
			itemPrices[item] = price
		}
	}

	err = nil
	for _, item := range items {
		if itemPrices[item] == 0.0 {
			log.Error("Failed to query or call API for item price: ", item)
			err = fmt.Errorf("failed to query or call API for at least one item")
		}
	}

	return itemPrices, err
}

func callPriceAPI(items []Item) (map[Item]float64, error) {
	log.Debug("Calling price API for items: ", items)
	prices := make(map[Item]float64)
	qualityToItems := make(map[uint8][]Item)

	for _, item := range items {
		qualityToItems[item.Quality] = append(qualityToItems[item.Quality], item)
	}

	for quality, itemsAtQuality := range qualityToItems {
		qualityPrices, err := callPriceAPIForQuality(itemsAtQuality, quality)
		if err != nil {
			log.Error("Failed to call pricing api for items: ", itemsAtQuality)
		}
		for item, price := range qualityPrices {
			prices[item] = price
		}
	}

	return prices, nil
}

func callPriceAPIForQuality(items []Item, quality uint8) (map[Item]float64, error) {
	prices := make(map[Item]float64)
	urls := getPriceAPIUrls(items, quality)

	for _, url := range urls {
		priceGroups := make(map[string][]float64)
		log.Debug("Calling price url: ", url)
		response, err := http.Get(url)
		if err != nil {
			log.Error("The HTTP request failed with error ", err)
			return prices, fmt.Errorf("the HTTP request failed with error %s", err)
		}
		defer response.Body.Close()

		// Read the response body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Error("Failed to read the response body: ", err)
			return prices, fmt.Errorf("failed to read the response body: %s", err)

		}

		// Use Gjson to parse and query the JSON response
		json := string(body)
		if !gjson.Valid(json) {
			log.Error("Invalid json resonse from url: ", url)
			return prices, fmt.Errorf("invalid json response from url: %s", url)
		}
		// Example: Iterate over all events and print the Killer's Name
		gjson.Parse(json).ForEach(func(_, result gjson.Result) bool {
			price := result.Get("sell_price_min").Float()
			if price == 0.0 {
				return true
			}
			priceGroups[result.Get("item_id").String()] = append(priceGroups[result.Get("item_id").String()], price)
			return true // keep iterating
		})
		for typeString, itemPrices := range priceGroups {
			item, err := typeStringToItem(typeString, quality)
			if err != nil {
				log.Error("Failed to parse type string: ", typeString)
				return prices, err
			}
			prices[item] = calculateMedian(itemPrices)
		}
	}

	return prices, nil
}

func getPriceAPIUrls(items []Item, quality uint8) []string {
	var locations string
	for _, location := range config.PriceLocations {
		locations += location + ","
	}
	locations = locations[:len(locations)-1]
	var urls []string

	i := 0
	var itemList string

	for _, item := range items {
		if i == 50 {
			i = 0
			itemList = itemList[:len(itemList)-1]
			urls = append(urls, config.PriceUrl+"/"+itemList+".json?locations="+locations+"&qualities="+fmt.Sprintf("%d", quality+1))
			itemList = ""
		}
		i += 1
		itemList += itemToTypeString(item) + ","
	}
	if itemList != "" {
		itemList = itemList[:len(itemList)-1]
		urls = append(urls, config.PriceUrl+"/"+itemList+".json?locations="+locations+"&qualities="+fmt.Sprintf("%d", quality+1))
	}

	return urls
}

func calculateMedian(data []float64) float64 {
	// Sort the data
	sort.Float64s(data)

	// Determine the median
	n := len(data)
	if n == 0 {
		return 0.0
	}

	mid := n / 2
	if n%2 == 1 {
		// Odd number of elements
		return data[mid]
	} else {
		// Even number of elements
		return (data[mid-1] + data[mid]) / 2.0
	}
}

func cachePricesFromEvents(events []Event) {
	itemsSet := make(map[Item]bool)
	for _, event := range events {
		itemsSet[event.KillerBuild.MainHand] = event.KillerBuild.MainHand.Name != ""
		itemsSet[event.KillerBuild.OffHand] = event.KillerBuild.OffHand.Name != ""
		itemsSet[event.KillerBuild.Head] = event.KillerBuild.Head.Name != ""
		itemsSet[event.KillerBuild.Chest] = event.KillerBuild.Chest.Name != ""
		itemsSet[event.KillerBuild.Foot] = event.KillerBuild.Foot.Name != ""
		itemsSet[event.KillerBuild.Cape] = event.KillerBuild.Cape.Name != ""
		itemsSet[event.KillerBuild.Potion] = event.KillerBuild.Potion.Name != ""
		itemsSet[event.KillerBuild.Food] = event.KillerBuild.Food.Name != ""
		itemsSet[event.KillerBuild.Mount] = event.KillerBuild.Mount.Name != ""
		itemsSet[event.KillerBuild.Bag] = event.KillerBuild.Bag.Name != ""
		itemsSet[event.VictimBuild.MainHand] = event.VictimBuild.MainHand.Name != ""
		itemsSet[event.VictimBuild.OffHand] = event.VictimBuild.OffHand.Name != ""
		itemsSet[event.VictimBuild.Head] = event.VictimBuild.Head.Name != ""
		itemsSet[event.VictimBuild.Chest] = event.VictimBuild.Chest.Name != ""
		itemsSet[event.VictimBuild.Foot] = event.VictimBuild.Foot.Name != ""
		itemsSet[event.VictimBuild.Cape] = event.VictimBuild.Cape.Name != ""
		itemsSet[event.VictimBuild.Potion] = event.VictimBuild.Potion.Name != ""
		itemsSet[event.VictimBuild.Food] = event.VictimBuild.Food.Name != ""
		itemsSet[event.VictimBuild.Mount] = event.VictimBuild.Mount.Name != ""
		itemsSet[event.VictimBuild.Bag] = event.VictimBuild.Bag.Name != ""
	}

	var items []Item
	for item, included := range itemsSet {
		if included && !strings.HasSuffix(item.Name, "_NONTRADABLE") {
			items = append(items, item)
		}
	}

	_, _ = getItemPrices(items)
}
