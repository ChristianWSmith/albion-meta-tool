package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func databaseCleanup() {
	for {
		time.Sleep(config.EventCleanupInterval)

		db, err := sql.Open("sqlite3", config.Database)
		if err != nil {
			log.Error("Failed to open database: ", err)
		}

		threshold := time.Now().Add(-config.EventStaleThreshold)

		query := `DELETE FROM events WHERE timestamp < ?`
		result, err := db.Exec(query, threshold)
		if err != nil {
			log.Error("Failed to clean up database: ", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Error("Failed to get rows affected during clean up database: ", err)
		}
		log.Debug("Deleted ", rowsAffected, " old records")

		db.Close()
	}
}

func updatePrices(itemPrices map[Item]float64) error {
	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		log.Error("Failed to open database: ", err)
		return err
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// Rollback the transaction if there's an error
			tx.Rollback()
			log.Println("Transaction rolled back due to error:", err)
		} else {
			// Commit the transaction if successful
			err = tx.Commit()
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	// Prepare the insert or update statement
	stmt, err := tx.Prepare(`INSERT INTO prices (
			name, tier, enchantment, quality, price, timestamp
		) VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(name, tier, enchantment, quality) DO UPDATE SET 
			price = excluded.price,
			timestamp = excluded.timestamp`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Iterate through items and execute the statement
	for item, price := range itemPrices {
		_, err = stmt.Exec(item.Name, item.Tier, item.Enchantment, item.Quality, price, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

func splitArray(arr []Item, maxSize int) [][]Item {
	var result [][]Item

	for i := 0; i < len(arr); i += maxSize {
		end := i + maxSize

		if end > len(arr) {
			end = len(arr)
		}

		result = append(result, arr[i:end])
	}

	return result
}

func queryPrices(items []Item) (map[Item]float64, error) {
	itemPrices := make(map[Item]float64)
	var errs []error

	itemBatches := splitArray(items, 249)

	for _, itemBatch := range itemBatches {
		itemBatchPrices, err := queryPricesBatch(itemBatch)
		if err != nil {
			log.Error("Failed to query price batch: ", itemBatchPrices)
		}
		for batchItem, batchPrice := range itemBatchPrices {
			itemPrices[batchItem] = batchPrice
		}
	}

	if len(errs) != 0 {
		return itemPrices, fmt.Errorf("%v", errs)
	}
	return itemPrices, nil
}

func queryPricesBatch(items []Item) (map[Item]float64, error) {
	itemPrices := make(map[Item]float64)

	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		log.Error("Failed to open database: ", err)
		return itemPrices, err
	}
	defer db.Close()

	// Prepare the query
	var placeholders []string
	var params []interface{}
	for _, item := range items {
		placeholders = append(placeholders, "(?, ?, ?, ?)")
		params = append(params, item.Name, item.Tier, item.Enchantment, item.Quality)
	}
	query := fmt.Sprintf(`SELECT name, tier, enchantment, quality, price, timestamp FROM prices WHERE (name, tier, enchantment, quality) IN (%s)`, strings.Join(placeholders, ","))

	// Execute the query
	rows, err := db.Query(query, params...)
	if err != nil {
		log.Error("Failed to execute query for item prices: ", err)
		return itemPrices, err
	}
	defer rows.Close()

	// Iterate through the result set
	for rows.Next() {
		var item Item
		var price float64
		var timestamp time.Time
		if err := rows.Scan(&item.Name, &item.Tier, &item.Enchantment, &item.Quality, &price, &timestamp); err != nil {
			log.Error("Failed to scan record into item struct: ", err)
			return itemPrices, err
		}
		// Store the price in the map with the item as the key
		if time.Since(timestamp) > config.PriceStaleThreshold {
			itemPrices[item] = 0.0
		} else {
			itemPrices[item] = price
		}
	}
	if err := rows.Err(); err != nil {
		log.Error("Row error encountered while querying item prices: ", err)
		return itemPrices, err
	}

	return itemPrices, nil
}

func insertEvents(events []Event) error {
	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		log.Error("Failed to open database: ", err)
		return err
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Error("Failed to begin transaction: ", err)
		return err
	}

	// Prepare the insert statement within the transaction
	stmt, err := tx.Prepare(`INSERT OR IGNORE INTO events (
		id, 
		killer_main_hand_name, killer_main_hand_tier, killer_main_hand_enchantment, killer_main_hand_quality, 
		killer_off_hand_name, killer_off_hand_tier, killer_off_hand_enchantment, killer_off_hand_quality,
        killer_head_name, killer_head_tier, killer_head_enchantment, killer_head_quality,
        killer_chest_name, killer_chest_tier, killer_chest_enchantment, killer_chest_quality,
        killer_foot_name, killer_foot_tier, killer_foot_enchantment, killer_foot_quality,
        killer_cape_name, killer_cape_tier, killer_cape_enchantment, killer_cape_quality,
        killer_potion_name, killer_potion_tier, killer_potion_enchantment,
        killer_food_name, killer_food_tier, killer_food_enchantment,
		killer_mount_name, killer_mount_tier, killer_mount_enchantment, killer_mount_quality,
		killer_bag_name, killer_bag_tier, killer_bag_enchantment, killer_bag_quality,
		killer_average_ip,
        victim_main_hand_name, victim_main_hand_tier, victim_main_hand_enchantment, victim_main_hand_quality,
        victim_off_hand_name, victim_off_hand_tier, victim_off_hand_enchantment, victim_off_hand_quality,
        victim_head_name, victim_head_tier, victim_head_enchantment, victim_head_quality,
        victim_chest_name, victim_chest_tier, victim_chest_enchantment, victim_chest_quality,
        victim_foot_name, victim_foot_tier, victim_foot_enchantment, victim_foot_quality,
        victim_cape_name, victim_cape_tier, victim_cape_enchantment, victim_cape_quality,
        victim_potion_name, victim_potion_tier, victim_potion_enchantment,
        victim_food_name, victim_food_tier, victim_food_enchantment,
        victim_mount_name, victim_mount_tier, victim_mount_enchantment, victim_mount_quality,
        victim_bag_name, victim_bag_tier, victim_bag_enchantment, victim_bag_quality,
        victim_average_ip,
        number_of_participants,
        timestamp) VALUES 
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Error("Failed to prepare sql statement: ", err)
		return err
	}
	defer stmt.Close()

	// Execute batch insert within the transaction
	for _, event := range events {
		log.Debug("Inserting event: ", event.EventId)
		_, err = stmt.Exec(
			event.EventId,
			event.KillerBuild.MainHand.Name, event.KillerBuild.MainHand.Tier, event.KillerBuild.MainHand.Enchantment, event.KillerBuild.MainHand.Quality,
			event.KillerBuild.OffHand.Name, event.KillerBuild.OffHand.Tier, event.KillerBuild.OffHand.Enchantment, event.KillerBuild.OffHand.Quality,
			event.KillerBuild.Head.Name, event.KillerBuild.Head.Tier, event.KillerBuild.Head.Enchantment, event.KillerBuild.Head.Quality,
			event.KillerBuild.Chest.Name, event.KillerBuild.Chest.Tier, event.KillerBuild.Chest.Enchantment, event.KillerBuild.Chest.Quality,
			event.KillerBuild.Foot.Name, event.KillerBuild.Foot.Tier, event.KillerBuild.Foot.Enchantment, event.KillerBuild.Foot.Quality,
			event.KillerBuild.Cape.Name, event.KillerBuild.Cape.Tier, event.KillerBuild.Cape.Enchantment, event.KillerBuild.Cape.Quality,
			event.KillerBuild.Potion.Name, event.KillerBuild.Potion.Tier, event.KillerBuild.Potion.Enchantment,
			event.KillerBuild.Food.Name, event.KillerBuild.Food.Tier, event.KillerBuild.Food.Enchantment,
			event.KillerBuild.Mount.Name, event.KillerBuild.Mount.Tier, event.KillerBuild.Mount.Enchantment, event.KillerBuild.Mount.Quality,
			event.KillerBuild.Bag.Name, event.KillerBuild.Bag.Tier, event.KillerBuild.Bag.Enchantment, event.KillerBuild.Bag.Quality,
			event.KillerAverageIp,
			event.VictimBuild.MainHand.Name, event.VictimBuild.MainHand.Tier, event.VictimBuild.MainHand.Enchantment, event.VictimBuild.MainHand.Quality,
			event.VictimBuild.OffHand.Name, event.VictimBuild.OffHand.Tier, event.VictimBuild.OffHand.Enchantment, event.VictimBuild.OffHand.Quality,
			event.VictimBuild.Head.Name, event.VictimBuild.Head.Tier, event.VictimBuild.Head.Enchantment, event.VictimBuild.Head.Quality,
			event.VictimBuild.Chest.Name, event.VictimBuild.Chest.Tier, event.VictimBuild.Chest.Enchantment, event.VictimBuild.Chest.Quality,
			event.VictimBuild.Foot.Name, event.VictimBuild.Foot.Tier, event.VictimBuild.Foot.Enchantment, event.VictimBuild.Foot.Quality,
			event.VictimBuild.Cape.Name, event.VictimBuild.Cape.Tier, event.VictimBuild.Cape.Enchantment, event.VictimBuild.Cape.Quality,
			event.VictimBuild.Potion.Name, event.VictimBuild.Potion.Tier, event.VictimBuild.Potion.Enchantment,
			event.VictimBuild.Food.Name, event.VictimBuild.Food.Tier, event.VictimBuild.Food.Enchantment,
			event.VictimBuild.Mount.Name, event.VictimBuild.Mount.Tier, event.VictimBuild.Mount.Enchantment, event.VictimBuild.Mount.Quality,
			event.VictimBuild.Bag.Name, event.VictimBuild.Bag.Tier, event.VictimBuild.Bag.Enchantment, event.VictimBuild.Bag.Quality,
			event.VictimAverageIp,
			event.NumberOfParticipants,
			event.Timestamp,
		)
		if err != nil {
			tx.Rollback() // Rollback the transaction in case of an error
			log.Error("Failed insert for event: ", event, err)
			return err
		}
	}

	// Commit the transaction
	log.Info("Commiting transaction to database")
	err = tx.Commit()
	if err != nil {
		log.Error("Failed to commit: ", err)
		return err
	}

	return nil
}

func initDatabase() error {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		log.Error("Failed to open database: ", err)
		return err
	}
	defer db.Close()

	// Create a table
	createTables := `
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY,
			killer_main_hand_name TEXT, killer_main_hand_tier INTEGER, killer_main_hand_enchantment INTEGER, killer_main_hand_quality INTEGER,
			killer_off_hand_name TEXT, killer_off_hand_tier INTEGER, killer_off_hand_enchantment INTEGER, killer_off_hand_quality INTEGER,
			killer_head_name TEXT, killer_head_tier INTEGER, killer_head_enchantment INTEGER, killer_head_quality INTEGER,
			killer_chest_name TEXT, killer_chest_tier INTEGER, killer_chest_enchantment INTEGER, killer_chest_quality INTEGER,
			killer_foot_name TEXT, killer_foot_tier INTEGER, killer_foot_enchantment INTEGER, killer_foot_quality INTEGER,
			killer_cape_name TEXT, killer_cape_tier INTEGER, killer_cape_enchantment INTEGER, killer_cape_quality INTEGER,
			killer_potion_name TEXT, killer_potion_tier INTEGER, killer_potion_enchantment INTEGER,
			killer_food_name TEXT, killer_food_tier INTEGER, killer_food_enchantment INTEGER,
			killer_mount_name TEXT, killer_mount_tier INTEGER, killer_mount_enchantment INTEGER, killer_mount_quality INTEGER,
			killer_bag_name TEXT, killer_bag_tier INTEGER, killer_bag_enchantment INTEGER, killer_bag_quality INTEGER,
			killer_average_ip REAL,
			victim_main_hand_name TEXT, victim_main_hand_tier INTEGER, victim_main_hand_enchantment INTEGER, victim_main_hand_quality INTEGER,
			victim_off_hand_name TEXT, victim_off_hand_tier INTEGER, victim_off_hand_enchantment INTEGER, victim_off_hand_quality INTEGER,
			victim_head_name TEXT, victim_head_tier INTEGER, victim_head_enchantment INTEGER, victim_head_quality INTEGER,
			victim_chest_name TEXT, victim_chest_tier INTEGER, victim_chest_enchantment INTEGER, victim_chest_quality INTEGER,
			victim_foot_name TEXT, victim_foot_tier INTEGER, victim_foot_enchantment INTEGER, victim_foot_quality INTEGER,
			victim_cape_name TEXT, victim_cape_tier INTEGER, victim_cape_enchantment INTEGER, victim_cape_quality INTEGER,
			victim_potion_name TEXT, victim_potion_tier INTEGER, victim_potion_enchantment INTEGER,
			victim_food_name TEXT, victim_food_tier INTEGER, victim_food_enchantment INTEGER,
			victim_mount_name TEXT, victim_mount_tier INTEGER, victim_mount_enchantment INTEGER, victim_mount_quality INTEGER,
			victim_bag_name TEXT, victim_bag_tier INTEGER, victim_bag_enchantment INTEGER, victim_bag_quality INTEGER,
			victim_average_ip REAL,
			number_of_participants INTEGER,
			timestamp DATETIME
		);
		CREATE TABLE IF NOT EXISTS prices (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			tier INTEGER,
			enchantment INTEGER,
			quality INTEGER,
			price REAL,
			timestamp DATETIME
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_prices_unique ON prices (name, tier, enchantment, quality);
	`
	if _, err := db.Exec(createTables); err != nil {
		log.Error("Failed to create tables: ", err)
		return err
	}
	return nil
}

func queryAllEvents() ([]Event, error) {
	// Connect to the SQLite database
	var err error
	var events []Event
	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		log.Error("Failed to open database: ", err)
		return events, err
	}
	defer db.Close()

	query := `SELECT 
		id,
		killer_main_hand_name, killer_main_hand_tier, killer_main_hand_enchantment, killer_main_hand_quality,
		killer_off_hand_name, killer_off_hand_tier, killer_off_hand_enchantment, killer_off_hand_quality,
		killer_head_name, killer_head_tier, killer_head_enchantment, killer_head_quality,
		killer_chest_name, killer_chest_tier, killer_chest_enchantment, killer_chest_quality,
		killer_foot_name, killer_foot_tier, killer_foot_enchantment, killer_foot_quality,
		killer_cape_name, killer_cape_tier, killer_cape_enchantment, killer_cape_quality,
		killer_potion_name, killer_potion_tier, killer_potion_enchantment,
		killer_food_name, killer_food_tier, killer_food_enchantment,
		killer_mount_name, killer_mount_tier, killer_mount_enchantment, killer_mount_quality,
		killer_bag_name, killer_bag_tier, killer_bag_enchantment, killer_bag_quality,
		killer_average_ip,
		victim_main_hand_name, victim_main_hand_tier, victim_main_hand_enchantment, victim_main_hand_quality,
		victim_off_hand_name, victim_off_hand_tier, victim_off_hand_enchantment, victim_off_hand_quality,
		victim_head_name, victim_head_tier, victim_head_enchantment, victim_head_quality,
		victim_chest_name, victim_chest_tier, victim_chest_enchantment, victim_chest_quality,
		victim_foot_name, victim_foot_tier, victim_foot_enchantment, victim_foot_quality,
		victim_cape_name, victim_cape_tier, victim_cape_enchantment, victim_cape_quality,
		victim_potion_name, victim_potion_tier, victim_potion_enchantment,
		victim_food_name, victim_food_tier, victim_food_enchantment,
		victim_mount_name, victim_mount_tier, victim_mount_enchantment, victim_mount_quality,
		victim_bag_name, victim_bag_tier, victim_bag_enchantment, victim_bag_quality,
		victim_average_ip,
		number_of_participants, timestamp
	FROM events`

	rows, err := db.Query(query)
	if err != nil {
		log.Error("Query for all events failed: ", err)
		return events, err
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.EventId,
			&event.KillerBuild.MainHand.Name, &event.KillerBuild.MainHand.Tier, &event.KillerBuild.MainHand.Enchantment, &event.KillerBuild.MainHand.Quality,
			&event.KillerBuild.OffHand.Name, &event.KillerBuild.OffHand.Tier, &event.KillerBuild.OffHand.Enchantment, &event.KillerBuild.OffHand.Quality,
			&event.KillerBuild.Head.Name, &event.KillerBuild.Head.Tier, &event.KillerBuild.Head.Enchantment, &event.KillerBuild.Head.Quality,
			&event.KillerBuild.Chest.Name, &event.KillerBuild.Chest.Tier, &event.KillerBuild.Chest.Enchantment, &event.KillerBuild.Chest.Quality,
			&event.KillerBuild.Foot.Name, &event.KillerBuild.Foot.Tier, &event.KillerBuild.Foot.Enchantment, &event.KillerBuild.Foot.Quality,
			&event.KillerBuild.Cape.Name, &event.KillerBuild.Cape.Tier, &event.KillerBuild.Cape.Enchantment, &event.KillerBuild.Cape.Quality,
			&event.KillerBuild.Potion.Name, &event.KillerBuild.Potion.Tier, &event.KillerBuild.Potion.Enchantment,
			&event.KillerBuild.Food.Name, &event.KillerBuild.Food.Tier, &event.KillerBuild.Food.Enchantment,
			&event.KillerBuild.Mount.Name, &event.KillerBuild.Mount.Tier, &event.KillerBuild.Mount.Enchantment, &event.KillerBuild.Mount.Quality,
			&event.KillerBuild.Bag.Name, &event.KillerBuild.Bag.Tier, &event.KillerBuild.Bag.Enchantment, &event.KillerBuild.Bag.Quality,
			&event.KillerAverageIp,
			&event.VictimBuild.MainHand.Name, &event.VictimBuild.MainHand.Tier, &event.VictimBuild.MainHand.Enchantment, &event.VictimBuild.MainHand.Quality,
			&event.VictimBuild.OffHand.Name, &event.VictimBuild.OffHand.Tier, &event.VictimBuild.OffHand.Enchantment, &event.VictimBuild.OffHand.Quality,
			&event.VictimBuild.Head.Name, &event.VictimBuild.Head.Tier, &event.VictimBuild.Head.Enchantment, &event.VictimBuild.Head.Quality,
			&event.VictimBuild.Chest.Name, &event.VictimBuild.Chest.Tier, &event.VictimBuild.Chest.Enchantment, &event.VictimBuild.Chest.Quality,
			&event.VictimBuild.Foot.Name, &event.VictimBuild.Foot.Tier, &event.VictimBuild.Foot.Enchantment, &event.VictimBuild.Foot.Quality,
			&event.VictimBuild.Cape.Name, &event.VictimBuild.Cape.Tier, &event.VictimBuild.Cape.Enchantment, &event.VictimBuild.Cape.Quality,
			&event.VictimBuild.Potion.Name, &event.VictimBuild.Potion.Tier, &event.VictimBuild.Potion.Enchantment,
			&event.VictimBuild.Food.Name, &event.VictimBuild.Food.Tier, &event.VictimBuild.Food.Enchantment,
			&event.VictimBuild.Mount.Name, &event.VictimBuild.Mount.Tier, &event.VictimBuild.Mount.Enchantment, &event.VictimBuild.Mount.Quality,
			&event.VictimBuild.Bag.Name, &event.VictimBuild.Bag.Tier, &event.VictimBuild.Bag.Enchantment, &event.VictimBuild.Bag.Quality,
			&event.VictimAverageIp,
			&event.NumberOfParticipants, &event.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func getNumEvents() (int, error) {
	var count int
	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		log.Error("Failed to open database: ", err)
		return count, err
	}
	defer db.Close()

	// Query to count rows in a table
	query := "SELECT COUNT(*) FROM events"

	// Execute the query
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Error("Error while getting number of events: ", err)
		return count, err
	}

	return count, nil
}

func getNumPrices() (int, error) {
	var count int
	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		log.Error("Failed to open database: ", err)
		return count, err
	}
	defer db.Close()

	// Query to count rows in a table
	query := "SELECT COUNT(*) FROM prices"

	// Execute the query
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Error("Error while getting number of prices: ", err)
		return count, err
	}

	return count, nil
}
