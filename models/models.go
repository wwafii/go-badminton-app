package models


type Court struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Timeslot struct {
	ID        int    `json:"id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type Reservation struct {
	ID               int    `json:"id"`
	CourtID          int    `json:"court_id"`
	TimeslotID       int    `json:"timeslot_id"`
	ReservationDate  string `json:"reservation_date"` // YYYY-MM-DD
	CustomerName     string `json:"customer_name"`
	Status           string `json:"status"` // "Pending", "Confirmed", "Cancelled"
	Amount           int    `json:"amount"` 
}