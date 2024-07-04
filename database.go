package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func initDatabase() error {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		logError("Failed to open database", err)
		return err
	}
	defer db.Close()

	// Create a table
	createTables := `
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY,
			killer_main_hand TEXT,
			killer_off_hand TEXT,
			killer_head TEXT,
			killer_chest TEXT,
			killer_foot TEXT,
			killer_cape TEXT,
			killer_potion TEXT,
			killer_food TEXT,
			killer_mount TEXT,
			killer_bag TEXT,
			killer_average_ip REAL,
			victim_main_hand TEXT,
			victim_off_hand TEXT,
			victim_head TEXT,
			victim_chest TEXT,
			victim_foot TEXT,
			victim_cape TEXT,
			victim_potion TEXT,
			victim_food TEXT,
			victim_mount TEXT,
			victim_bag TEXT,
			victim_average_ip REAL,
			number_of_participants INTEGER,
			timestamp DATETIME
		);
		CREATE TABLE IF NOT EXISTS prices (
			id INTEGER PRIMARY KEY,
			name TEXT,
			tier INTEGER,
			enchantment INTEGER,
			quality INTEGER,
			price REAL,
			timestamp DATETIME
		);
	`
	if _, err := db.Exec(createTables); err != nil {
		logError("", err)
		return err
	}
	return nil
}
