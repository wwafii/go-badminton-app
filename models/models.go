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
	ReservationDate  string `json:"reservation_date"` 
	CustomerName     string `json:"customer_name"`
	Status           string `json:"status"` 
	PaymentReference string `json:"payment_reference,omitempty"`
	Amount           int    `json:"amount"`
	MidtransOrderID  string `json:"midtrans_order_id"`
	MidtransToken    string `json:"midtrans_token,omitempty"`
}