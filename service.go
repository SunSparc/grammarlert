package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-martini/martini"
)

// Autocomplete endpoint; suggests hosts based on prefix
func SuggestHosts(params martini.Params, resp http.ResponseWriter) {
	bytes, err := json.Marshal(FindHostsByPrefix(params["prefix"]))
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	} else {
		resp.Header().Set("Content-Type", "application/json")
		resp.Write(bytes)
	}
}

// Save a user's incoming report in the database as a full report
func FileReport(report InboundReport, req *http.Request, resp http.ResponseWriter) {
	lchost := strings.ToLower(report.ParsedURL.Host)
	uri := report.ParsedURL.RequestURI()

	host := db.getOrCreateHost(lchost)

	if _, exists := host.URIs[uri]; !exists {
		host.URIs[uri] = make(map[string]Report)
	}
	page := host.URIs[uri]

	if _, exists := page[report.OriginalText]; !exists {
		page[report.OriginalText] = Report{
			ID:       UuidStr(),
			Original: report.OriginalText,
			Status:   RPT_NEW,
		}
	}
	rpt := page[report.OriginalText]

	rpt.Suggestions = append(rpt.Suggestions, Suggestion{
		Suggestion: report.SuggestedText,
		Note:       report.Note,
		Submitter:  req.RemoteAddr,
		Date:       time.Now().UTC(),
	})

	host.Reported++

	page[report.OriginalText] = rpt
	host.URIs[uri] = page
	db.Hosts[lchost] = host

	index.Put(lchost)
	db.Counts.Total++
	db.Counts.New++
	Save()

	resp.WriteHeader(http.StatusOK)
}

// Get all reports for a certain host
func GetReports(params martini.Params, resp http.ResponseWriter) {
	host, exists := GetHost(strings.ToLower(params["host"]))
	if !exists {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(host.URIs)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	} else {
		resp.Header().Set("Content-Type", "application/json")
		resp.Write(bytes)
	}
}

func Acknowledge(action ActionMeta, params martini.Params, resp http.ResponseWriter) {
	if SetFlag(action.Host, action.Page, params["reportID"], RPT_ACK) {
		resp.WriteHeader(http.StatusOK)
	} else {
		resp.WriteHeader(http.StatusNotFound)
	}
}

func Fixed(action ActionMeta, params martini.Params, resp http.ResponseWriter) {
	if SetFlag(action.Host, action.Page, params["reportID"], RPT_FIX) {
		resp.WriteHeader(http.StatusOK)
	} else {
		resp.WriteHeader(http.StatusNotFound)
	}
}

func Reject(action ActionMeta, params martini.Params, resp http.ResponseWriter) {
	if SetFlag(action.Host, action.Page, params["reportID"], RPT_REJ) {
		resp.WriteHeader(http.StatusOK)
	} else {
		resp.WriteHeader(http.StatusNotFound)
	}
}

func Delete(action ActionMeta, params martini.Params, resp http.ResponseWriter) {
	if SetFlag(action.Host, action.Page, params["reportID"], RPT_DEL) {
		resp.WriteHeader(http.StatusOK)
	} else {
		resp.WriteHeader(http.StatusNotFound)
	}
}

// Used to assist with finding a particular report in the database
type ActionMeta struct {
	Host string `form:"host"`
	Page string `form:"uri"`
}
