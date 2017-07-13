package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

const (
	persistentFile        = "data.gob" // File to persist data to disk
	persistentSeconds     = 10         // Period in seconds to persist data to disk
	requestTimeoutSeconds = 20         // Seconds since last request per transaction to timeout message
)

type Fragment struct {
	Offset  int    `json:"offset"`
	TransID int    `json:"trans_id"`
	Payload string `json:"payload"`
	Size    int    `json:"size"`
}

// Decodes payload from Base64
func (f *Fragment) decodePayload() error {
	payload, err := base64.StdEncoding.DecodeString(f.Payload)
	if err != nil {
		return err
	}
	f.Payload = string(payload[:])
	return nil
}

func init() {
	d = NewDatabase()
	dbLock = &sync.Mutex{}

	// Load data if available
	file, err := os.Open(persistentFile)
	if err != nil {
		log.Println("[INFO] No existing data to load")
		d = NewDatabase()
		return
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)

	dbLock.Lock()
	defer dbLock.Unlock()
	var db Database

	decoder.Decode(&db)

	if len(db.Data) == 0 {
		db.Data = map[int]map[int]string{}
	}

	log.Printf("[INFO] Successfully loaded %v transactions from disk\n", len(db.Data))

	// Set loaded db as global db instance
	d = &db
}

// Globals
var d *Database
var dbLock *sync.Mutex

func main() {

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		var f Fragment
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			log.Println("Unabled to decode body")
		}
		f.decodePayload()

		d.Add(&f)
	})

	// Persist to disk every `persistentSeconds` seconds
	ticker := time.NewTicker(persistentSeconds * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				d.Persist()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// Start the server
	log.Fatal(http.ListenAndServe(":1234", nil))

}

type Database struct {
	Data map[int]map[int]string
}

// Creates a new Database struct
func NewDatabase() *Database {
	return &Database{Data: map[int]map[int]string{}}
}

// Adds a message fragment to the database
func (d *Database) Add(f *Fragment) {

	timers := map[int]*time.Timer{}

	// Acquire lock to write to the database
	dbLock.Lock()
	defer dbLock.Unlock()

	// Check if new transaction
	if _, ok := d.Data[f.TransID]; !ok {
		// Create new data
		d.Data[f.TransID] = make(map[int]string)
		// Create timer
		timers[f.TransID] = time.AfterFunc(requestTimeoutSeconds*time.Second, func() {
			var buffer bytes.Buffer
			message := d.CheckForCompletedMessage(f.TransID)

			if len(message) == 0 {
				log.Printf("Successfully received full message [transaction ID: %v]!\n", f.TransID)
				
				return
			}

			for i := 0; i < len(message); i += 2 {
				buffer.WriteString(fmt.Sprintf("%v-%v, ", message[i], message[i+1]))
			}

			log.Println("[ERROR] Never recieved the full message. Hole(s) at:")
			log.Println(buffer.String())
		})
	} else {
		// Otherwise we've seen this transaction before, reset the timer
		timer := timers[f.TransID]
		if timer != nil {
			timers[f.TransID].Reset(requestTimeoutSeconds * time.Second)
		}
	}

	// Add payload to database
	for i := 1; i <= f.Size; i++ {
		d.Data[f.TransID][f.Offset+i] = f.Payload
	}
}

// Checks if a message has been fully received
func (d *Database) CheckForCompletedMessage(transactionID int) []int {
	// Get location of all bytes
	var keys []int
	dbLock.Lock()
	defer dbLock.Unlock()
	for k := range d.Data[transactionID] {
		keys = append(keys, k)
	}

	// Sort byte locations
	sort.Ints(keys)

	// Completed messages should be fully contiguous
	holes := []int{}
	startIndex := keys[0]
	for idx, item := range keys {
		if idx+startIndex != item {
			holes = append(holes, keys[idx-1])
			holes = append(holes, item)
			startIndex += item - (idx + startIndex)
		}
	}
	return holes
}

// Persist the database to disk
func (d *Database) Persist() {
	log.Println("[INFO] Persisting to disk...")
	outFile, err := os.Create(persistentFile)

	if err != nil {
		log.Println("No such file")
	} else {
		defer outFile.Close()
		encoder := gob.NewEncoder(outFile)
		dbLock.Lock()
		defer dbLock.Unlock()
		if err := encoder.Encode(d); err != nil {
			log.Println("Unable to encode")
		}
		log.Printf("[INFO] Successfully persisted %v transactions to disk\n", len(d.Data))
	}
}
