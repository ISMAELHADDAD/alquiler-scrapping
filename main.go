package main

import (
	"log"
	"os"
	"time"

	"alquiler-scrapping/collectors"
	"alquiler-scrapping/database"
	"alquiler-scrapping/telegram"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func checkForUpdates(db *database.Database, tg *telegram.Telegram) {

	// 1. Scrape the listings from the page and save them in the database
	collectors.CollectHabitacliaEntries(db, 700)

	// 2. Iterate through listings that haven't been sent to the Telegram channel yet
	entries, _ := db.ListNotSent()
	var err error
	for _, entry := range entries {
		_, err = tg.SendToChannel(entry)
		if err != nil {
			log.Printf("Failed to send to Telegram channel. Entry url: %s\n", entry.Url)
		}

		if os.Getenv("DEBUG") == "true" {
			log.Printf("Entry sent: %s\n", entry.Url)
		}

		// Mark the listing as sent in the database
		db.MarkAsSent(entry)

		// Wait 3 seconds before sending the next listing
		time.Sleep(3 * time.Second)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	tg, err := telegram.NewTelegramBot()
	if err != nil {
		log.Fatal(err)
	}

	checkForUpdates(db, tg)
}
