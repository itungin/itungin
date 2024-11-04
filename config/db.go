package config

import (
	"github.com/gocroot/helper/atdb"
)

// Mendefinisikan MongoString secara langsung
var MongoString string = "mongodb+srv://karamissuu:karamissu1@cluster0.lyovb.mongodb.net/"

// Konfigurasi database dengan nama 'akuntan'
var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "akuntan",
}

// Membuat koneksi ke MongoDB
var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)
