package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/bfv/pascollector/misc"
	"github.com/bfv/pascollector/types"
	bolt "go.etcd.io/bbolt"
)

// var storeInitialized = false
var db *bolt.DB

func dataStore(ch chan types.Metric) {

	initStore()

	for {
		metric := <-ch

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("DB")).Bucket([]byte("METRICS"))
			err := b.Put([]byte(metric.Id), metricToByteArray(metric))

			if err != nil {
				fmt.Println()
			}

			return nil
		})
	}

}

func metricToByteArray(m interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(m)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func byteArrayToMetric(b []byte) types.Metric {
	m := types.Metric{}
	dec := gob.NewDecoder(bytes.NewReader(b))
	err := dec.Decode(&m)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func CloseDB() {
	fmt.Println("closing database...")
	db.Close()
}

func AcknowledgeSent(sent []string) {

}

func initStore() {

	dbFile := misc.GetDatabaseDir() + "metrics.db"
	dbIn, err := bolt.Open(dbFile, 0666, nil)

	if err == nil {
		db = dbIn
	} else {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {

		root, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			log.Fatalln("error creating DB bucket:", err)
		}

		_, err = root.CreateBucketIfNotExists([]byte("METRICS"))
		if err != nil {
			log.Fatalln("error creating METRICS bucket:", err)
		}

		return nil
	})

}

func getStoredData(chData chan types.Metric, chQuit chan bool) {

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB")).Bucket([]byte("METRICS"))
		b.ForEach(func(k, v []byte) error {
			metric := byteArrayToMetric(v)
			chData <- metric
			return nil
		})
		return nil
	})
	chQuit <- true
}

func deleteMetric(id string) {

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB")).Bucket([]byte("METRICS"))
		b.Delete([]byte(id))
		return nil
	})
}
