package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go-badminton-app/models"
	"go-badminton-app/config"
)


var (
	Courts = []models.Court{
		{ID: 1, Name: "Court Eagle"},
		{ID: 2, Name: "Court Smash"},
	}
	Timeslots = []models.Timeslot{
		{ID: 101, StartTime: "08:00", EndTime: "09:00"},
		{ID: 102, StartTime: "09:00", EndTime: "10:00"},
		{ID: 103, StartTime: "10:00", EndTime: "11:00"},
		{ID: 104, StartTime: "11:00", EndTime: "12:00"},
	}
	Reservations = make(map[int]models.Reservation)
	mu           sync.Mutex 
	nextResID    = 1
)



func isCourtBooked(date string, timeslotID int, courtID int) bool {
	for _, res := range Reservations {
		if res.ReservationDate == date && res.TimeslotID == timeslotID && res.CourtID == courtID && res.Status != "Cancelled" {
			return true
		}
	}
	return false
}

// GetDatesHandler (Logika tetap sama)
func GetDatesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	today := time.Now()
	var dates []string
	
	for i := 0; i < 7; i++ {
		date := today.AddDate(0, 0, i)
		dates = append(dates, date.Format("2006-01-02")) 
	}

	json.NewEncoder(w).Encode(dates)
}

func AvailableTimeslotsHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, "Query parameter 'date' wajib diisi", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	var bookedSlots = make(map[int]bool) 
	
	for _, ts := range Timeslots {
		bookedCourts := 0
		for _, res := range Reservations {
			if res.ReservationDate == dateStr && res.TimeslotID == ts.ID && res.Status != "Cancelled" {
				bookedCourts++
			}
		}
		if bookedCourts >= len(Courts) {
			bookedSlots[ts.ID] = true
		}
	}

	var availableSlots []models.Timeslot
	for _, ts := range Timeslots {
		if !bookedSlots[ts.ID] {
			availableSlots = append(availableSlots, ts)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(availableSlots)
}

// AvailableCourtsHandler (Logika tetap sama)
func AvailableCourtsHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	timeslotIDStr := r.URL.Query().Get("timeslot_id")

	if dateStr == "" || timeslotIDStr == "" {
		http.Error(w, "Parameter 'date' dan 'timeslot_id' wajib diisi", http.StatusBadRequest)
		return
	}
	timeslotID, err := strconv.Atoi(timeslotIDStr)
	if err != nil {
		http.Error(w, "timeslot_id harus berupa angka", http.StatusBadRequest)
		return
	}

	var availableCourts []models.Court

	mu.Lock()
	defer mu.Unlock()

	for _, c := range Courts {
		if !isCourtBooked(dateStr, timeslotID, c.ID) {
			availableCourts = append(availableCourts, c)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(availableCourts)
}

// CreateReservationHandler (Reservasi dibuat dengan status Pending)
func CreateReservationHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CourtID         int    `json:"court_id"`
		TimeslotID      int    `json:"timeslot_id"`
		ReservationDate string `json:"reservation_date"`
		CustomerName    string `json:"customer_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	
	if req.CustomerName == "" || req.CourtID <= 0 || req.TimeslotID <= 0 {
		http.Error(w, "Semua field wajib diisi dan valid.", http.StatusBadRequest)
		return
	}
	if _, err := time.Parse("2006-01-02", req.ReservationDate); err != nil {
		http.Error(w, "Format tanggal reservasi tidak valid (harus YYYY-MM-DD).", http.StatusBadRequest)
		return
	}
	
	
	mu.Lock()
	defer mu.Unlock()
	if isCourtBooked(req.ReservationDate, req.TimeslotID, req.CourtID) {
		http.Error(w, "Lapangan sudah dibooking, silakan pilih yang lain.", http.StatusConflict)
		return
	}

	
	newReservation := models.Reservation{
		ID: nextResID,
		CourtID: req.CourtID,
		TimeslotID: req.TimeslotID,
		ReservationDate: req.ReservationDate,
		CustomerName: req.CustomerName,
		Amount: config.PRICE_PER_SLOT,
		Status: "Pending", 
	}
	
	Reservations[nextResID] = newReservation
	nextResID++

	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Reservasi berhasil dibuat. Menunggu konfirmasi pembayaran.",
		"reservation_id": newReservation.ID,
		"amount": newReservation.Amount,
	})
}


func PaymentConfirmHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ReservationID int `json:"reservation_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	
	mu.Lock()
	defer mu.Unlock()

	if res, ok := Reservations[req.ReservationID]; ok {
		if res.Status == "Pending" {
			res.Status = "Confirmed" // KONFIRMASI BOOKING SLOT
			Reservations[req.ReservationID] = res
			fmt.Printf("Reservasi ID %d BERHASIL DIKONFIRMASI (Simulasi Paid).\n", req.ReservationID)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Pembayaran sukses. Reservasi Confirmed."})
			return
		}
		http.Error(w, "Reservasi sudah dikonfirmasi atau dibatalkan.", http.StatusConflict)
		return
	}
	http.Error(w, "Reservasi tidak ditemukan.", http.StatusNotFound)
}