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


var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)



// Mendefinisikan MongoString secara langsung
var MongoStringGeo string = "mongodb+srv://Cito:w.cito.a@cluster0.svl9a.mongodb.net/"

// Konfigurasi database dengan nama 'akuntan'
var mongoinfoGeo = atdb.DBInfo{
	DBString: MongoStringGeo,
	DBName:   "Geo",
}

var MongoconnGeo, ErrorMongoconnGeo = atdb.MongoConnect(mongoinfoGeo)



