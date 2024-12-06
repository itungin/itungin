package model

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Product adalah struct untuk produk
type Product struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Name        string  `bson:"name" json:"name"`
    Price       float64 `bson:"price" json:"price"`
    Category    string  `bson:"category" json:"category"`
    Description string  `bson:"description" json:"description"`
    Stock       int     `bson:"stock" json:"stock"`
    CreatedAt   time.Time   `bson:"createdAt" json:"createdAt"`
    UpdatedAt   time.Time   `bson:"updatedAt" json:"updatedAt"`
}

type Customer struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Phone     string             `bson:"phone" json:"phone"`
	Address   string             `bson:"address" json:"address"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Name     string             `bson:"name" json:"name"`
    Email    string             `bson:"email" json:"email"`
    Password string             `bson:"password" json:"password"`
    Nohp string                 `bson:"no_hp" json:"no_hp"`
}

// ExpenseTransaction adalah struct untuk transaksi pengeluaran
type ExpenseTransaction struct {
    ID            primitive.ObjectID    `bson:"_id,omitempty" json:"id,omitempty"`
    ExpenseName  string    `bson:"expense_name" json:"expense_name"`   // Nama pengeluaran (misalnya: sewa, gaji, dll.)
    Amount        float64   `bson:"amount" json:"amount"`               // Jumlah uang yang dikeluarkan
    Category      string    `bson:"category" json:"category"`           // Kategori pengeluaran (misalnya: operasional, marketing, dll.)
    PaymentMethod string   `bson:"payment_method" json:"payment_method"` // Metode pembayaran (misalnya: transfer bank, tunai)
    ExpenseDate  time.Time `bson:"expense_date" json:"expense_date"`   // Tanggal pengeluaran
    Notes         string    `bson:"notes" json:"notes,omitempty"`       // Catatan tambahan (opsional)
    CreatedAt    time.Time `bson:"created_at" json:"created_at"`       // Waktu transaksi dibuat
    UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`       // Waktu transaksi terakhir diperbarui
}


// SalesTransaction adalah struct untuk transaksi penjualan
type SalesTransaction struct {
    ID            primitive.ObjectID    `bson:"_id,omitempty" json:"id,omitempty"`
    TransactionDate time.Time `bson:"transactionDate" json:"transactionDate"`
    CustomerName  string    `bson:"customer_name" json:"customer_name"`
    Products      []Product `bson:"products" json:"products"`
    TotalAmount   float64   `bson:"total_amount" json:"total_amount"`
    PaymentMethod string    `bson:"payment_method" json:"payment_method"`
    PaymentStatus string    `bson:"payment_status" json:"payment_status"`
}

type LaporanAkuntan struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    StartDate  string             `bson:"startDate" json:"startDate"`   // Menggunakan string
    EndDate    string             `bson:"endDate" json:"endDate"`       // Menggunakan string
    StartDateTime time.Time       `bson:"startDateTime" json:"startDateTime"` // Menyimpan waktu yang sudah diparse
    EndDateTime   time.Time       `bson:"endDateTime" json:"endDateTime"`
	Income    float64   `bson:"income" json:"income"`
	Expenses  float64   `bson:"expenses" json:"expenses"`
	Profit    float64   `bson:"profit" json:"profit"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type Employee struct {
    ID           primitive.ObjectID       `json:"id" bson:"_id"`
    Name    string    `json:"name" bson:"name"`
    Email        string    `json:"email" bson:"email"`
    PhoneNumber  string    `json:"phone_number" bson:"phone_number"`
    Position     string    `json:"position" bson:"position"`
    CreatedAt    time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

// Category adalah struct untuk kategori produk
type Category struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Name        string             `bson:"name" json:"name"`
    Description string             `bson:"description" json:"description"`
    CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
    UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
