package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
	"bytes"
	"io"

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

func createMidtransTransaction(res models.Reservation) (string, error) {
	
	payload := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id": res.MidtransOrderID,
			"gross_amount": res.Amount,
		},
		"customer_details": map[string]string{
			"first_name": res.CustomerName,
			"email": "customer@example.com", 
			"phone": "081111111111",
		},
		"callbacks": map[string]string{
			"notification": config.SERVER_URL + "/api/payment-callback", 
		},
	}
	
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("gagal marshall payload: %w", err)
	}

	
	_, err = http.NewRequest(http.MethodPost, config.MIDTRANS_SNAP_URL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("gagal membuat request Midtrans: %w", err)
	}
	
	
	if res.Amount <= 0 {
		return "", fmt.Errorf("jumlah pembayaran tidak valid")
	}
	mockToken := fmt.Sprintf("mock-snap-token-%s", res.MidtransOrderID)

	return mockToken, nil 
}




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
		MidtransOrderID: fmt.Sprintf("DIRO-%d-%d", nextResID, time.Now().Unix()),
	}
	
	
	snapToken, err := createMidtransTransaction(newReservation)
	if err != nil {
		http.Error(w, fmt.Sprintf("Gagal memproses pembayaran: %v", err), http.StatusInternalServerError)
		return
	}
	
	
	newReservation.MidtransToken = snapToken
	Reservations[nextResID] = newReservation
	nextResID++

	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Reservasi berhasil dibuat. Arahkan ke pembayaran.",
		"reservation_id": newReservation.ID,
		"amount": newReservation.Amount,
		"snap_token": snapToken, 
	})
}

func PaymentCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
    
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) 

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	
	var notificationPayload map[string]interface{}
	if err := json.Unmarshal(body, &notificationPayload); err != nil {
		http.Error(w, "Invalid notification payload JSON", http.StatusBadRequest)
		return
	}
	
	orderID, _ := notificationPayload["order_id"].(string)
	transactionStatus, _ := notificationPayload["transaction_status"].(string)

	
	mu.Lock()
	defer mu.Unlock()
	
	var resIDToUpdate int
	resFound := false
	for id, res := range Reservations {
		if res.MidtransOrderID == orderID {
			resIDToUpdate = id
			resFound = true
			break
		}
	}
	
	if !resFound {
		http.Error(w, "Order ID not found in database", http.StatusNotFound)
		return
	}
	
	res := Reservations[resIDToUpdate]
	
	
	if transactionStatus == "capture" || transactionStatus == "settlement" {
		if res.Status != "Confirmed" {
			 res.Status = "Confirmed" 
			 Reservations[resIDToUpdate] = res
			 fmt.Printf("Reservasi %s BERHASIL DIKONFIRMASI (Paid) via Midtrans.\n", orderID)
		}
	} else if transactionStatus == "deny" || transactionStatus == "cancel" || transactionStatus == "expire" {
		res.Status = "Cancelled"
		Reservations[resIDToUpdate] = res
		fmt.Printf("Reservasi %s DIBATALKAN/KADALUARSA.\n", orderID)
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Midtrans notification handled successfully"})
}