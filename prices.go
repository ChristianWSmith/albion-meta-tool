package main

func getBuildPrice(build Build) (float64, error) {
	var err error
	var main_hand_price float64
	var off_hand_price float64
	var head_price float64
	var chest_price float64
	var foot_price float64
	var cape_price float64
	var potion_price float64
	var food_price float64
	var bag_price float64
	var mount_price float64

	main_hand_price, err = getItemPrice(build.MainHand)
	if err != nil {
		log.Error("Failed to get build price: ", build)
		return 0.0, err
	}

	off_hand_price, err = getItemPrice(build.OffHand)
	if err != nil {
		log.Error("Failed to get build price: ", build)
		return 0.0, err
	}

	head_price, err = getItemPrice(build.Head)
	if err != nil {
		log.Error("Failed to get build price: ", build)
		return 0.0, err
	}

	chest_price, err = getItemPrice(build.Chest)
	if err != nil {
		log.Error("Failed to get build price: ", build)
		return 0.0, err
	}

	foot_price, err = getItemPrice(build.Foot)
	if err != nil {
		log.Error("Failed to get build price: ", build)
		return 0.0, err
	}

	cape_price, err = getItemPrice(build.Cape)
	if err != nil {
		log.Error("Failed to get build price: ", build)
		return 0.0, err
	}

	potion_price, err = getItemPrice(build.Potion)
	if err != nil {
		log.Error("Failed to get build price: ", build)
		return 0.0, err
	}

	food_price, err = getItemPrice(build.Food)
	if err != nil {
		log.Error("Failed to get build price: ", build)
		return 0.0, err
	}

	bag_price, err = getItemPrice(build.Bag)
	if err != nil {
		log.Error("Failed to get build price: ", build)
		return 0.0, err
	}

	mount_price, err = getItemPrice(build.Mount)
	if err != nil {
		log.Error("Failed to get build price: ", build)
		return 0.0, err
	}

	return main_hand_price + off_hand_price + head_price + chest_price + foot_price + cape_price + potion_price + food_price + bag_price + mount_price, nil
}

func getItemPrice(item Item) (float64, error) {
	var price float64
	var err error

	price, err = queryPrice(item)
	if err != nil {
		log.Error("Failed to query for price for item: ", item)
		return price, err
	}
	if price == 0.0 {
		price, err = callPriceAPI(item)
		if err != nil {
			log.Error("Failed to call price API: ", err)
			return price, err
		}
		err = updatePrice(item, price)
		if err != nil {
			log.Error("Failed to update price in Database: ", err)
			return price, err
		}
	}

	return price, nil
}

func callPriceAPI(item Item) (float64, error) {
	// TODO: actually call price API
	return 1.0, nil
}
