package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	// "github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
)

func GetRegion(respw http.ResponseWriter, req *http.Request) {

	var longlat model.LongLat
	json.NewDecoder(req.Body).Decode(&longlat)

	filter := bson.M{
		"border": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": bson.M{ // Menggunakan $geometry, bukan $border
					"type":        "Point",
					"coordinates": []float64{longlat.Longitude, longlat.Latitude},
				},
			},
		},
	}
	
	
	region, err := atdb.GetOneDoc[model.Region](config.MongoconnGeo, "region", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusNotFound, region)
		return
	}
	at.WriteJSON(respw, http.StatusOK, region)
}

// new anton add
func GetRoads(respw http.ResponseWriter, req *http.Request) {


	var longlat model.LongLat
	json.NewDecoder(req.Body).Decode(&longlat)	

	filter := bson.M{
			"geometry": bson.M{
				"$nearSphere": bson.M{
					"$geometry": bson.M{		
						"type":        "Point",
						"coordinates": []float64{longlat.Longitude, longlat.Latitude},
					},
					"$maxDistance": longlat.MaxDistance,
				},
			},
	}

	roads, err := atdb.GetAllDoc[[]model.Roads](config.MongoconnGeo, "roads", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusNotFound, roads)
		return
	}
	at.WriteJSON(respw, http.StatusOK, roads)
}
