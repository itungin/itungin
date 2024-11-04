package config

import (
	"context"
	"log"
	"time"

	"github.com/gocroot/helper/atdb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mendefinisikan MongoString secara langsung
var MongoString string = "mongodb+srv://karamissuu:karamissu1@cluster0.lyovb.mongodb.net/"

// Konfigurasi database dengan nama 'akuntan'
var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "akuntan",
}

var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)


// Membuat variabel untuk client MongoDB dan collections
var Client *mongo.Client
var UserCollection *mongo.Collection
var ProductCollection *mongo.Collection
var SalesTransactionCollection *mongo.Collection
var ExpenseTransactionCollection *mongo.Collection
var CustomerCollection *mongo.Collection
var ReportCollection *mongo.Collection

// Fungsi untuk menginisialisasi koneksi ke MongoDB
func InitMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(MongoString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Gagal terhubung ke MongoDB: %v", err)
	}

	// Memastikan koneksi berhasil
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Gagal ping ke MongoDB: %v", err)
	}

	log.Println("MongoDB connected")

	// Menetapkan variabel global client dan collections
	Client = client
	UserCollection = client.Database("akuntan").Collection("user")
	ProductCollection = client.Database("akuntan").Collection("produk")
	SalesTransactionCollection = client.Database("akuntan").Collection("transaksi_penjualan")
	ExpenseTransactionCollection = client.Database("akuntan").Collection("transaksi_pengeluaran")
	CustomerCollection = client.Database("akuntan").Collection("pelanggan")
	ReportCollection = client.Database("akuntan").Collection("laporan")
}


