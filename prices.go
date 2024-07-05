package main

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

	return itemPrices, nil
}

func callPriceAPI(items []Item) (map[Item]float64, error) {
	log.Debug("Calling price API for items: ", items)
	prices := make(map[Item]float64)

	for _, item := range items {
		prices[item] = 1.0
	}

	return prices, nil
}
