package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type (
	Database struct {
		Hosts  map[string]Host // keyed by hostname
		Counts struct {
			Total        uint
			New          uint
			Acknowledged uint
			Fixed        uint
			Rejected     uint
			Deleted      uint
		}
	}
	Host struct {
		ID           string
		Hostname     string
		Subscribers  []Subscriber
		URIs         map[string]Reports // keyed by full URI (path+query string)
		Created      time.Time
		Reported     uint
		Acknowledged uint
		Fixed        uint
		Rejected     uint
		Deleted      uint
	}
	Reports map[string]Report
	Report  struct {
		ID          string
		Original    string
		Status      uint
		Suggestions []Suggestion
		Comments    []Comment
	}
	Suggestion struct {
		Suggestion string
		Note       string
		Submitter  string
		Date       time.Time
		Comments   []Comment
	}
	Subscriber struct {
		ID               string
		Email            string
		SubscriptionDate time.Time
	}
	Comment struct {
		ID      string
		Author  string
		Body    string
		Created time.Time
		Edited  time.Time
	}
)

// Auto-saves the database to disk every interval
func (d *Database) autoSave() {
	for {
		time.Sleep(saveInterval)
		Save()
	}
}

// Loads the database from disk, or prepares a new one to be used
// if it doesn't already exist on disk
func (d *Database) load() {
	if !fileExists(dumpPath) {
		d.Hosts = make(map[string]Host)
		return
	}
	contents, err := ioutil.ReadFile(dumpPath)
	if err != nil {
		fmt.Println("**DB ERROR*** Could not load data")
		panic(err)
	}
	err = json.Unmarshal(contents, d)
	if err != nil {
		fmt.Println("**DB ERROR*** Could not unmarshal data")
		panic(err)
	}
}

// Saves the database to disk in JSON format
func (d *Database) save() {
	saveMutex.Lock()

	if !fileExists(dumpPath) {
		if _, err := os.Create(dumpPath); err != nil {
			fmt.Println("***ERROR*** Could not create dump file")
			panic(err)
		}
	}
	jsonBytes, err := json.Marshal(d)
	if err != nil {
		fmt.Println("**ERROR*** Could not marshal data")
		panic(err)
	}
	err = ioutil.WriteFile(dumpPath, jsonBytes, 0644)
	if err != nil {
		fmt.Println("**ERROR*** Could not write data")
		panic(err)
	}

	saveMutex.Unlock()
}

// Gets the host entry from the DB. It creates it if it doesn't already exist.
func (d *Database) getOrCreateHost(hostname string) Host {
	if _, exists := d.Hosts[hostname]; !exists {
		d.Hosts[hostname] = Host{
			ID:       UuidStr(),
			Hostname: hostname,
			Created:  time.Now().UTC(),
			URIs:     make(map[string]Reports),
		}
	}
	return d.Hosts[hostname]
}

// Returns whether a file exists or not
func fileExists(file string) bool {
	if _, err := os.Stat(dumpPath); os.IsNotExist(err) {
		return false
	}
	return true
}

// Sets a status flag on a report
func SetFlag(host, uri, id string, flag uint) bool {
	if hostentry, exists := db.Hosts[host]; exists {
		if reports, exists := hostentry.URIs[uri]; exists {
			for i, report := range reports {
				if report.ID == id {
					report.Status |= flag
					db.Counts.New--
					switch flag {
					case RPT_ACK:
						db.Counts.Acknowledged++
						hostentry.Acknowledged++
					case RPT_FIX:
						db.Counts.Fixed++
						hostentry.Fixed++
					case RPT_REJ:
						db.Counts.Rejected++
						hostentry.Rejected++
					case RPT_DEL:
						db.Counts.Deleted++
						hostentry.Deleted++
					}
					reports[i] = report
					hostentry.URIs[uri] = reports
					db.Hosts[host] = hostentry
					return true
				}
			}
		}
	}
	return false
}

// Clears a status flag on a report
func ClearFlag(host, uri, id string, flag uint) bool {
	if hostentry, exists := db.Hosts[host]; exists {
		if reports, exists := hostentry.URIs[uri]; exists {
			for i, report := range reports {
				if report.ID == id {
					report.Status ^= flag
					if report.Status == 0 {
						db.Counts.New++
					}
					switch flag {
					case RPT_ACK:
						db.Counts.Acknowledged--
						hostentry.Acknowledged--
					case RPT_FIX:
						db.Counts.Fixed--
						hostentry.Fixed--
					case RPT_REJ:
						db.Counts.Rejected--
						hostentry.Rejected--
					case RPT_DEL:
						db.Counts.Deleted--
						hostentry.Deleted--
					}
					reports[i] = report
					hostentry.URIs[uri] = reports
					db.Hosts[host] = hostentry
					return true
				}
			}
		}
	}
	return false
}

// Only one save operation at a time
var saveMutex = &sync.Mutex{}

// Database file
const dumpPath = "dump.json"

// Time between database dumps (other than filing new reports)
const saveInterval = 5 * time.Minute

const (
	RPT_NEW uint = 1<<iota - 1 // New, unacknowledged report
	RPT_ACK                    // Acknowledged
	RPT_FIX                    // Fixed
	RPT_REJ                    // Rejected (not used on the client right now, as it seems redundant to deleted)
	RPT_DEL                    // Deleted
)
