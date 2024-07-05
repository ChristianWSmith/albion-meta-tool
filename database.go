package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func updatePrice(item Item, price float64) error {
	var err error

	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		log.Error("Failed to open database: ", err)
		return err
	}
	defer db.Close()

	// Prepare the insert or update statement
	query := `INSERT INTO prices (
			name, tier, enchantment, quality, price, timestamp
		) VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(name, tier, enchantment, quality) DO UPDATE SET 
			price = excluded.price,
			timestamp = excluded.timestamp`

	// Execute the insert or update statement
	_, err = db.Exec(query, item.Name, item.Tier, item.Enchantment, item.Quality, price, time.Now())
	if err != nil {
		log.Error("Failed to insert price record for item: ", item)
		return err
	}

	return nil
}

func queryPrice(item Item) (float64, error) {
	var price float64
	var timestamp time.Time

	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		log.Error("Failed to open database: ", err)
		return price, err
	}
	defer db.Close()

	// Prepare the query
	query := `SELECT price, timestamp FROM prices WHERE name = ? AND tier = ? AND enchantment = ? AND quality = ?`

	// Execute the query
	err = db.QueryRow(query, item.Name, item.Tier, item.Enchantment, item.Quality).Scan(&price, &timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info("No record found for item: ", item)
		} else {
			log.Error("Query failed for item: ", item)
			return price, err
		}
	} else {
		fmt.Printf("The price is: %.2f\n", price)
	}
	if time.Since(timestamp) > config.PriceStaleThreshold {
		log.Info("Price is stale for item: ", item)
		return 0.0, nil
	}
	return price, nil
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
		killer_main_hand_name, 
		killer_main_hand_tier, 
		killer_main_hand_enchantment, 
		killer_main_hand_quality, 
		killer_off_hand_name, 
		killer_off_hand_tier, 
		killer_off_hand_enchantment, 
		killer_off_hand_quality,
        killer_head_name,
        killer_head_tier,
		killer_head_enchantment,
        killer_head_quality,
        killer_chest_name,
        killer_chest_tier,
		killer_chest_enchantment,
        killer_chest_quality,
        killer_foot_name,
        killer_foot_tier,
		killer_foot_enchantment,
        killer_foot_quality,
        killer_cape_name,
        killer_cape_tier,
		killer_cape_enchantment,
        killer_cape_quality,
        killer_potion_name,
        killer_potion_tier,
		killer_potion_enchantment,
        killer_food_name,
        killer_food_tier,
        killer_food_enchantment,
		killer_mount_name,
        killer_mount_tier,
        killer_mount_enchantment,
        killer_mount_quality,
		killer_bag_name,
        killer_bag_tier,
        killer_bag_enchantment,
        killer_bag_quality,
		killer_average_ip,
        victim_main_hand_name,
        victim_main_hand_tier,
		victim_main_hand_enchantment,
        victim_main_hand_quality,
        victim_off_hand_name,
		victim_off_hand_tier,
        victim_off_hand_enchantment,
        victim_off_hand_quality,
        victim_head_name,
		victim_head_tier,
        victim_head_enchantment,
        victim_head_quality,
        victim_chest_name,
		victim_chest_tier,
        victim_chest_enchantment,
        victim_chest_quality,
        victim_foot_name,
		victim_foot_tier,
        victim_foot_enchantment,
        victim_foot_quality,
        victim_cape_name,
		victim_cape_tier,
        victim_cape_enchantment,
        victim_cape_quality,
        victim_potion_name,
		victim_potion_tier,
        victim_potion_enchantment,
        victim_food_name,
        victim_food_tier,
		victim_food_enchantment,
        victim_mount_name,
        victim_mount_tier,
        victim_mount_enchantment,
        victim_mount_quality,
        victim_bag_name,
        victim_bag_tier,
        victim_bag_enchantment,
		victim_bag_quality,
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
			event.KillerBuild.MainHand.Name,
			event.KillerBuild.MainHand.Tier,
			event.KillerBuild.MainHand.Enchantment,
			event.KillerBuild.MainHand.Quality,
			event.KillerBuild.OffHand.Name,
			event.KillerBuild.OffHand.Tier,
			event.KillerBuild.OffHand.Enchantment,
			event.KillerBuild.OffHand.Quality,
			event.KillerBuild.Head.Name,
			event.KillerBuild.Head.Tier,
			event.KillerBuild.Head.Enchantment,
			event.KillerBuild.Head.Quality,
			event.KillerBuild.Chest.Name,
			event.KillerBuild.Chest.Tier,
			event.KillerBuild.Chest.Enchantment,
			event.KillerBuild.Chest.Quality,
			event.KillerBuild.Foot.Name,
			event.KillerBuild.Foot.Tier,
			event.KillerBuild.Foot.Enchantment,
			event.KillerBuild.Foot.Quality,
			event.KillerBuild.Cape.Name,
			event.KillerBuild.Cape.Tier,
			event.KillerBuild.Cape.Enchantment,
			event.KillerBuild.Cape.Quality,
			event.KillerBuild.Potion.Name,
			event.KillerBuild.Potion.Tier,
			event.KillerBuild.Potion.Enchantment,
			event.KillerBuild.Food.Name,
			event.KillerBuild.Food.Tier,
			event.KillerBuild.Food.Enchantment,
			event.KillerBuild.Mount.Name,
			event.KillerBuild.Mount.Tier,
			event.KillerBuild.Mount.Enchantment,
			event.KillerBuild.Mount.Quality,
			event.KillerBuild.Bag.Name,
			event.KillerBuild.Bag.Tier,
			event.KillerBuild.Bag.Enchantment,
			event.KillerBuild.Bag.Quality,
			event.KillerAverageIp,
			event.VictimBuild.MainHand.Name,
			event.VictimBuild.MainHand.Tier,
			event.VictimBuild.MainHand.Enchantment,
			event.VictimBuild.MainHand.Quality,
			event.VictimBuild.OffHand.Name,
			event.VictimBuild.OffHand.Tier,
			event.VictimBuild.OffHand.Enchantment,
			event.VictimBuild.OffHand.Quality,
			event.VictimBuild.Head.Name,
			event.VictimBuild.Head.Tier,
			event.VictimBuild.Head.Enchantment,
			event.VictimBuild.Head.Quality,
			event.VictimBuild.Chest.Name,
			event.VictimBuild.Chest.Tier,
			event.VictimBuild.Chest.Enchantment,
			event.VictimBuild.Chest.Quality,
			event.VictimBuild.Foot.Name,
			event.VictimBuild.Foot.Tier,
			event.VictimBuild.Foot.Enchantment,
			event.VictimBuild.Foot.Quality,
			event.VictimBuild.Cape.Name,
			event.VictimBuild.Cape.Tier,
			event.VictimBuild.Cape.Enchantment,
			event.VictimBuild.Cape.Quality,
			event.VictimBuild.Potion.Name,
			event.VictimBuild.Potion.Tier,
			event.VictimBuild.Potion.Enchantment,
			event.VictimBuild.Food.Name,
			event.VictimBuild.Food.Tier,
			event.VictimBuild.Food.Enchantment,
			event.VictimBuild.Mount.Name,
			event.VictimBuild.Mount.Tier,
			event.VictimBuild.Mount.Enchantment,
			event.VictimBuild.Mount.Quality,
			event.VictimBuild.Bag.Name,
			event.VictimBuild.Bag.Tier,
			event.VictimBuild.Bag.Enchantment,
			event.VictimBuild.Bag.Quality,
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
			killer_main_hand_name TEXT,
			killer_main_hand_tier INTEGER,
			killer_main_hand_enchantment INTEGER,
			killer_main_hand_quality INTEGER,
			killer_off_hand_name TEXT,
			killer_off_hand_tier INTEGER,
			killer_off_hand_enchantment INTEGER,
			killer_off_hand_quality INTEGER,
			killer_head_name TEXT,
			killer_head_tier INTEGER,
			killer_head_enchantment INTEGER,
			killer_head_quality INTEGER,
			killer_chest_name TEXT,
			killer_chest_tier INTEGER,
			killer_chest_enchantment INTEGER,
			killer_chest_quality INTEGER,
			killer_foot_name TEXT,
			killer_foot_tier INTEGER,
			killer_foot_enchantment INTEGER,
			killer_foot_quality INTEGER,
			killer_cape_name TEXT,
			killer_cape_tier INTEGER,
			killer_cape_enchantment INTEGER,
			killer_cape_quality INTEGER,
			killer_potion_name TEXT,
			killer_potion_tier INTEGER,
			killer_potion_enchantment INTEGER,
			killer_food_name TEXT,
			killer_food_tier INTEGER,
			killer_food_enchantment INTEGER,
			killer_mount_name TEXT,
			killer_mount_tier INTEGER,
			killer_mount_enchantment INTEGER,
			killer_mount_quality INTEGER,
			killer_bag_name TEXT,
			killer_bag_tier INTEGER,
			killer_bag_enchantment INTEGER,
			killer_bag_quality INTEGER,
			killer_average_ip REAL,
			victim_main_hand_name TEXT,
			victim_main_hand_tier INTEGER,
			victim_main_hand_enchantment INTEGER,
			victim_main_hand_quality INTEGER,
			victim_off_hand_name TEXT,
			victim_off_hand_tier INTEGER,
			victim_off_hand_enchantment INTEGER,
			victim_off_hand_quality INTEGER,
			victim_head_name TEXT,
			victim_head_tier INTEGER,
			victim_head_enchantment INTEGER,
			victim_head_quality INTEGER,
			victim_chest_name TEXT,
			victim_chest_tier INTEGER,
			victim_chest_enchantment INTEGER,
			victim_chest_quality INTEGER,
			victim_foot_name TEXT,
			victim_foot_tier INTEGER,
			victim_foot_enchantment INTEGER,
			victim_foot_quality INTEGER,
			victim_cape_name TEXT,
			victim_cape_tier INTEGER,
			victim_cape_enchantment INTEGER,
			victim_cape_quality INTEGER,
			victim_potion_name TEXT,
			victim_potion_tier INTEGER,
			victim_potion_enchantment INTEGER,
			victim_food_name TEXT,
			victim_food_tier INTEGER,
			victim_food_enchantment INTEGER,
			victim_mount_name TEXT,
			victim_mount_tier INTEGER,
			victim_mount_enchantment INTEGER,
			victim_mount_quality INTEGER,
			victim_bag_name TEXT,
			victim_bag_tier INTEGER,
			victim_bag_enchantment INTEGER,
			victim_bag_quality INTEGER,
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
	`
	if _, err := db.Exec(createTables); err != nil {
		log.Error("Failed to create tables: ", err)
		return err
	}
	priceIndex := "CREATE UNIQUE INDEX IF NOT EXISTS idx_prices_unique ON prices (name, tier, enchantment, quality);"
	if _, err := db.Exec(priceIndex); err != nil {
		log.Error("Failed to create price index: ", err)
		return err
	}
	return nil
}
