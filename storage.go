package main

import (
	"strings"

	"github.com/nu7hatch/gouuid"
)

// Loads the database into memory and populates the autocomplete structure
func LoadDatabase() {
	db.load()
	for host := range db.Hosts {
		index.Put(host)
	}
}

// Saves the database to a JSON file
func Save() {
	go db.save()
}

// Automatically dumps the database to JSON every so often
func StartAutoSave() {
	go db.autoSave()
}

// Returns up to 10 hostnames starting with the given prefix
func FindHostsByPrefix(prefix string) []string {
	return index.Find(strings.ToLower(prefix), 10)
}

// Returns whether or not the given hostname is found in the data
func HasHost(hostname string) bool {
	_, exists := GetHost(strings.ToLower(hostname))
	return exists
}

// Gets all the data about a hostname from the database
func GetHost(hostname string) (Host, bool) {
	host, exists := db.Hosts[strings.ToLower(hostname)]
	return host, exists
}

// Returns a Version 4 UUID string
func UuidStr() string {
	u, _ := uuid.NewV4()
	return u.String()
}

var (
	// The database holds all the app's data
	db Database

	// This autocomplete index makes it easy to suggest hosts by prefix
	index Set = NewSet()
)
