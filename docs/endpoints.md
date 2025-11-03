# 🚀 Backend API Endpoints Documentation

This document provides detailed information for the RESTful API endpoints of the Go (Golang) backend service. All responses are in JSON format.

---

## 1. Availability Endpoints

These endpoints are used for dynamically checking court availability based on date and time.

### 1.1. Get Available Dates

Mengembalikan daftar 7 hari ke depan untuk ditampilkan di kalender/selector tanggal.

* **Endpoint:** `/api/dates`
* **Method:** `GET`
* **Query Parameters:** None
* **Success Response (200 OK):**
    ```json
    [
        "2025-11-03",
        "2025-11-04",
        "2025-11-05",
        "2025-11-06",
        "2025-11-07",
        "2025-11-08",
        "2025-11-09"
    ]
    ```

### 1.2. Get Available Timeslots

Mengembalikan slot waktu yang *masih memiliki setidaknya satu lapangan* yang tersedia pada tanggal tertentu.

* **Endpoint:** `/api/available-timeslots`
* **Method:** `GET`
* **Query Parameters:**
    | Name | Type | Description |
    | :--- | :--- | :--- |
    | `date` | `string` | Tanggal yang dipilih (Format: `YYYY-MM-DD`). **Wajib**. |
* **Success Response (200 OK):**
    ```json
    [
        {
            "id": 101,
            "start_time": "08:00",
            "end_time": "09:00"
        },
        {
            "id": 102,
            "start_time": "09:00",
            "end_time": "10:00"
        }
    ]
    ```

### 1.3. Get Available Courts

Mengembalikan daftar Lapangan yang benar-benar kosong pada slot waktu dan tanggal yang dipilih.

* **Endpoint:** `/api/available-courts`
* **Method:** `GET`
* **Query Parameters:**
    | Name | Type | Description |
    | :--- | :--- | :--- |
    | `date` | `string` | Tanggal yang dipilih (Format: `YYYY-MM-DD`). **Wajib**. |
    | `timeslot_id` | `integer` | ID Slot Waktu yang dipilih (e.g., 101). **Wajib**. |
* **Success Response (200 OK):**
    ```json
    [
        {
            "id": 1,
            "name": "Court Eagle"
        },
        {
            "id": 2,
            "name": "Court Smash"
        }
    ]
    ```

---

## 2. Transaction Endpoints

Ini adalah *endpoint* untuk membuat reservasi dan menangani pembayaran.

### 2.1. Create Reservation

Membuat reservasi baru dengan status `Pending` dan memicu pembuatan transaksi di Midtrans.

* **Endpoint:** `/api/reservations`
* **Method:** `POST`
* **Headers:** `Content-Type: application/json`
* **Request Body:**
    | Name | Type | Description |
    | :--- | :--- | :--- |
    | `court_id` | `integer` | ID Lapangan yang dipilih. |
    | `timeslot_id` | `integer` | ID Slot Waktu yang dipilih. |
    | `reservation_date`| `string` | Tanggal reservasi (`YYYY-MM-DD`). |
    | `customer_name` | `string` | Nama pelanggan. |
* **Success Response (200 OK):**
    ```json
    {
        "message": "Reservasi berhasil dibuat. Arahkan ke pembayaran.",
        "reservation_id": 3,
        "amount": 50000,
        "snap_token": "mock-snap-token-DIRO-3-1698822000" 
        // Frontend harus menggunakan token ini untuk memunculkan pop-up Midtrans Snap.
    }
    ```
* **Error Responses:**
    * `400 Bad Request`: Validasi input gagal (mis. Nama kosong, format tanggal salah).
    * `409 Conflict`: Lapangan sudah dibooking.
    * `500 Internal Server Error`: Kegagalan saat memanggil API Midtrans.

### 2.2. Payment Callback (Webhook Handler)

**Catatan PENTING:** *Endpoint* ini hanya dipanggil oleh *server* Midtrans, BUKAN oleh *frontend*!

* **Endpoint:** `/api/payment-callback`
* **Method:** `POST`
* **Headers:** Dipanggil oleh Midtrans.
* **Fungsi:** Menerima notifikasi status pembayaran. Jika status `settlement` atau `capture`, *backend* mengubah status reservasi menjadi **`Confirmed`** (mengamankan *slot* untuk pelanggan).

---