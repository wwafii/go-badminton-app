package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)


var (
	SERVER_PORT         string
	SERVER_URL          string
	PRICE_PER_SLOT      int 
)

func InitConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file. Pastikan file .env ada.")
	}

	SERVER_PORT = os.Getenv("SERVER_PORT")
	SERVER_URL = os.Getenv("SERVER_URL")
	
	priceStr := os.Getenv("PRICE_PER_SLOT")
	if price, err := strconv.Atoi(priceStr); err == nil {
		PRICE_PER_SLOT = price
	} else {
		PRICE_PER_SLOT = 50000 
	}
	
	if SERVER_PORT == "" {
		log.Fatal("SERVER_PORT tidak ditemukan di .env")
	}
}