package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const (
	STATUS_CREATED               = 201
	STATUS_DUPLICATE             = 409
	STATUS_OK                    = 200
	STATUS_FOUND                 = 302
	STATUS_BAD_REQUEST           = 400
	STATUS_UNAUTHORIZED          = 401
	STATUS_INTERNAL_SERVER_ERROR = 500
	STATUS_SERVICE_UNREACHABLE   = 503
)

var LookupKeys []string

func SetupProducts() {
	jsonData, err := ioutil.ReadFile("seed.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		panic(err)
	}

	fixtures := data["fixtures"].([]interface{})
	for _, fixture := range fixtures {
		fixtureMap := fixture.(map[string]interface{})
		params := fixtureMap["params"].(map[string]interface{})
		if lookupKey, ok := params["lookup_key"].(string); ok {
			LookupKeys = append(LookupKeys, lookupKey)
		}
	}
}

const (
	META_SUCCESS = 1
	META_FAILED  = 0
)
