package main

import (
	"fmt"
	"net/http"
	"log"

	"go-badminton-app/config" 
	"go-badminton-app/handlers" 
)

func main() {
	// 1. Inisialisasi Konfigurasi
	config.InitConfig()

	// 2. Routing 
	http.HandleFunc("/api/dates", handlers.GetDatesHandler)
	http.HandleFunc("/api/available-timeslots", handlers.AvailableTimeslotsHandler)
	http.HandleFunc("/api/available-courts", handlers.AvailableCourtsHandler)
	http.HandleFunc("/api/reservations", handlers.CreateReservationHandler)
	http.HandleFunc("/api/payment-confirm", handlers.PaymentConfirmHandler) 

	// 3. Konfigurasi CORS (Next.js)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Menggunakan localhost:3000 karena ini adalah environment Next.js
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") 
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.DefaultServeMux.ServeHTTP(w, r)
	})

	// 4. Jalankan Server
	fmt.Printf("Server Golang berjalan di :%s\n", config.SERVER_PORT)
	if err := http.ListenAndServe(":"+config.SERVER_PORT, handler); err != nil {
		log.Fatal(err)
	}
}