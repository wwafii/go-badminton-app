package main

import (
	"fmt"
	"net/http"
	"log"

	"go-badminton-app/config" 
	"go-badminton-app/handlers" 
)

func main() {
	
	config.InitConfig()


	http.HandleFunc("/api/dates", handlers.GetDatesHandler)
	http.HandleFunc("/api/available-timeslots", handlers.AvailableTimeslotsHandler)
	http.HandleFunc("/api/available-courts", handlers.AvailableCourtsHandler)
	http.HandleFunc("/api/reservations", handlers.CreateReservationHandler)
	http.HandleFunc("/api/payment-callback", handlers.PaymentCallbackHandler)

	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") 
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.DefaultServeMux.ServeHTTP(w, r)
	})

	
	fmt.Printf("Server Golang berjalan di :%s\n", config.SERVER_PORT)
	if err := http.ListenAndServe(":"+config.SERVER_PORT, handler); err != nil {
		log.Fatal(err)
	}
}