package helper

import (
	"encoding/json"
	"festajuninaweb/areas/commonstruct"
	"fmt"
	"io/ioutil"

	"github.com/go-redis/redis"
)

var redisclient *redis.Client
var SYSID string
var databaseEV commonstruct.DatabaseX
var envirvar commonstruct.RestEnvVariables

// Readfileintostruct is
func Readfileintostruct() commonstruct.RestEnvVariables {
	dat, err := ioutil.ReadFile("feijoadajusu.ini")
	check(err)
	fmt.Print(string(dat))

	var restenv commonstruct.RestEnvVariables

	json.Unmarshal(dat, &restenv)

	return restenv
}

// GetSYSID is just returning the System ID directly from file
// It is happening to enable multiple usage of Redis Keys ("SYSID" + "APIURL" for instance)
func GetSYSID() string {

	if SYSID == "" {

		dat, err := ioutil.ReadFile("feijoadajusu.ini")
		check(err)
		fmt.Print(string(dat))

		var restenv commonstruct.RestEnvVariables

		json.Unmarshal(dat, &restenv)

		SYSID = restenv.SYSID

		return restenv.SYSID
	}

	return SYSID

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// GetRedisPointer returns
func GetRedisPointer(bucket int) *redis.Client {

	bucket = 0

	if redisclient == nil {
		redisclient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",     // no password set
			DB:       bucket, // use default DB
		})
	}

	return redisclient
}

// Getvaluefromcache returns the value of a key from cache
func Getvaluefromcache(key string) string {

	// bucket is ZERO for now
	// I am allowing it to be setup now
	rp := GetRedisPointer(0)

	sysid := GetSYSID()

	valuetoreturn, _ := rp.Get(sysid + key).Result()

	return valuetoreturn
}

// GetDBParmFromCache returns the value of a key from cache
func GetDBParmFromCache(collection string) *commonstruct.DatabaseX {

	database := new(commonstruct.DatabaseX)

	database.Collection = Getvaluefromcache(collection)
	database.Database = Getvaluefromcache("API.MongoDB.Database")
	database.Location = Getvaluefromcache("API.MongoDB.Location")

	return database
}

// GetDBParmAllFromCache returns all data from cache
func GetDBParmAllFromCache() commonstruct.RestEnvVariables {

	if envirvar.APIAPIServerIPAddress == "" {
		envirvar.APIAPIServerIPAddress = Getvaluefromcache("Web.APIServer.IPAddress")
		envirvar.APIAPIServerPort = Getvaluefromcache("Web.APIServer.Port")
		envirvar.WEBServerPort = Getvaluefromcache("WEBServerPort")
		envirvar.WEBDebug = Getvaluefromcache("Web.Debug")
		envirvar.RecordCurrencyTick = Getvaluefromcache("RecordCurrencyTick")
		envirvar.RunningFromServer = Getvaluefromcache("RunningFromServer")
		envirvar.AppBelnorthEnabled = Getvaluefromcache("AppBelnorthEnabled")
		envirvar.AppBitcoinEnabled = Getvaluefromcache("AppBitcoinEnabled")
		envirvar.AppFestaJuninaEnabled = Getvaluefromcache("AppFestaJuninaEnabled")
		envirvar.MSAPImainPort = Getvaluefromcache("MSAPImainPort")
		envirvar.MSAPIdishesPort = Getvaluefromcache("MSAPIdishesPort")
		envirvar.MSAPIordersPort = Getvaluefromcache("MSAPIordersPort")
		envirvar.SecurityMicroservice = Getvaluefromcache("SecurityMicroservice")
		envirvar.SecurityMicroserviceURL = Getvaluefromcache("SecurityMicroserviceURL")
		envirvar.APIMongoDBLocation = Getvaluefromcache("APIMongoDBLocation")
		envirvar.APIMongoDBDatabase = Getvaluefromcache("APIMongoDBDatabase")
	}

	return envirvar

}
