package main

func getItemPrices(items []Item) (map[Item]float64, error) {
	var discovered_prices map[Item]float64
	var unknown_items []Item
	var price float64
	var err error

	prices := make(map[Item]float64)

	for _, item := range items {
		price, err = queryPrice(item) // TODO: batch this
		if err != nil {
			log.Error("Failed to query for price for item: ", item)
			return prices, err
		}
		if price == 0.0 {
			unknown_items = append(unknown_items, item)
		} else {
			prices[item] = price
		}
	}

	discovered_prices, err = callPriceAPI(unknown_items)
	if err != nil {
		log.Error("Failed to get prices for: ", unknown_items)
		return prices, err
	}
	for item, price := range discovered_prices {
		err = updatePrice(item, price) // TODO: batch this
		if err != nil {
			log.Error("Failed to update prices for item: ", item)
			return prices, err
		}
		prices[item] = price
	}

	return prices, nil
}

func callPriceAPI(items []Item) (map[Item]float64, error) {
	prices := make(map[Item]float64)

	for _, item := range items {
		prices[item] = 1.0
	}

	return prices, nil
}
