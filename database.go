package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func initDatabase() error {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		log.Error("Failed to open database", err)
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
			name TEXT PRIMARY KEY,
			tier INTEGER,
			enchantment INTEGER,
			quality INTEGER,
			price REAL,
			timestamp DATETIME
		);
	`
	if _, err := db.Exec(createTables); err != nil {
		log.Error("Failed to create tables", err)
		return err
	}
	return nil
}
