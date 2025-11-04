# 🏸 Badminton Court Reservation API (Golang Backend)

## 🎯 Project Focus

This repository contains the **Backend API** for a Badminton Court Reservation System, developed for a Full Stack Developer technical test. This service is built with Go (Golang) and adheres to Clean Code principles, focusing on separation of concerns (Handlers, Config, Models).

### Core Requirements Fulfilled

* ✅ **RESTful API:** Provides endpoints for date, timeslot, and court availability checks.
* ✅ **Reservation Logic:** Handles the core business logic, including checking for availability (no double booking).
* ✅ **Bonus Point Achieved (Payment Gateway):** Includes the necessary logic for non-blocking payment processing via Midtrans Webhook Callback.
* ✅ **Clean Code & Robustness:** Logic is separated into modules (`handlers/`, `config/`), and basic input validation is included.

---

## ⚙️ Technology Stack

| Component | Technology | Description |
| :--- | :--- | :--- |
| **Backend (API)** | **Go (Golang)** | Primary language for fast, concurrent, and robust API development. |
| **Configuration** | **`github.com/joho/godotenv`** | Used to load sensitive data and configurations from the `.env` file. |
| **Data Storage** | **In-Memory Map** | Used for simplicity in this test; data is stored in memory and resets upon server restart. |
| **Payment Gateway**| **Midtrans Snap (Sandbox)**| Integration framework for payment processing and asynchronous callback handling. |

---

## 🚀 API Endpoints

The server runs on `http://localhost:8080` (or the port defined in `.env`). All endpoints are prefixed with `/api`.

| Endpoint | Method | Keterangan |
| :--- | :--- | :--- |
| `/api/dates` | `GET` | Mengembalikan daftar tanggal yang valid (7 hari ke depan). |
| `/api/available-timeslots`| `GET` | **Query:** `?date=YYYY-MM-DD`. Mengembalikan slot waktu yang *tersedia*. |
| `/api/available-courts` | `GET` | **Query:** `?date=...&timeslot_id=...`. Mengembalikan Lapangan yang *tersedia*. |
| `/api/reservations` | `POST` | **INPUT:** Customer details, Court ID, Timeslot ID, Date. **OUTPUT:** `snap_token` Midtrans. (Status awal: `Pending`). |
| `/api/payment-callback` | `POST` | **[WEBHOOK HANDLER]** Menerima notifikasi dari Midtrans. Jika sukses (`settlement/capture`), status reservasi diubah menjadi `Confirmed`. |

*(Detail dokumentasi endpoint lengkap dapat dilihat di file: [`docs/endpoints.md`](docs/endpoints.md))*

---

## 🛠️ Getting Started (Local Setup)

### 1. Prerequisites

You must have **Go (1.18+)** installed.

### 2. Setup and Run

Navigate to the root directory of this repository.

**A. Configure Environment Variables**

Create a file named **`.env`** in the root directory:

```env
# .env file

# Server Configuration
SERVER_PORT=8080 
SERVER_URL=http://localhost:8080

# Midtrans Sandbox Credentials (HARUS DIGANTI)
MIDTRANS_SERVER_KEY=SB-Mid-server-YOUR_MIDTRANS_SERVER_KEY 
MIDTRANS_SNAP_URL=[https://app.sandbox.midtrans.com/snap/v1/transactions](https://app.sandbox.midtrans.com/snap/v1/transactions)

# 1. Initialize Go module dependencies
go mod tidy

# 2. Run the backend server
go run main.go