package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)


var (
	MIDTRANS_SERVER_KEY string
	MIDTRANS_SNAP_URL   string
	SERVER_PORT         string
	SERVER_URL          string
	PRICE_PER_SLOT      int = 50000 
)


func InitConfig() {
	
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file. Pastikan file .env ada di root folder.")
	}


	MIDTRANS_SERVER_KEY = os.Getenv("MIDTRANS_SERVER_KEY")
	MIDTRANS_SNAP_URL = os.Getenv("MIDTRANS_SNAP_URL")
	SERVER_PORT = os.Getenv("SERVER_PORT")
	SERVER_URL = os.Getenv("SERVER_URL")

	if MIDTRANS_SERVER_KEY == "" || SERVER_PORT == "" {
		log.Fatal("SERVER_PORT atau MIDTRANS_SERVER_KEY tidak ditemukan di .env")
	}
}