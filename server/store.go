package server

import (
	"github.com/bfv/pascollector/misc"
	"github.com/bfv/pascollector/types"
	bolt "go.etcd.io/bbolt"
)

var storeInitialized = false
var db *bolt.DB

func StoreData(metric types.Metric) {

	if !storeInitialized {
		initStore()
	}

}

func AcknowledgeSent(sent []string) {

}

func initStore() {

	dataDir := misc.GetDatabaseDir()

	dbIn, err := bolt.Open(dataDir, 0666, nil)
	if err == nil {
		db = dbIn
	}
}
