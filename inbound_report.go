package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/martini-contrib/binding"
)

// Represents an incoming report from a user
type InboundReport struct {
	URL           string   `json:"url"`
	ParsedURL     *url.URL `json:"-"`
	OriginalText  string   `json:"original"`
	SuggestedText string   `json:"suggest"`
	Note          string   `json:"note"`
}

// Performs validation on the incoming reports
func (r *InboundReport) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	if r.OriginalText == r.SuggestedText {
		errors = append(errors, binding.Error{
			FieldNames:     []string{"suggest"},
			Classification: "SameAsOriginal",
			Message:        "The suggestion cannot be what is already on the page",
		})
	}

	parsedUrl, err := url.Parse(r.URL)
	if err != nil || parsedUrl.Host == "" {
		errors = append(errors, binding.Error{
			FieldNames:     []string{"url"},
			Classification: "HostError",
			Message:        "The URL appears to be in a bad format (1)",
		})
	} else {
		r.ParsedURL = parsedUrl
	}
	r.ParsedURL = parsedUrl

	if r.ParsedURL.Host == "localhost" || strings.HasPrefix(r.ParsedURL.Host, "127.0.") {
		errors = append(errors, binding.Error{
			FieldNames:     []string{"url"},
			Classification: "HostError",
			Message:        "Local hosts are not accepted",
		})
	} else if len(r.URL) < 3 {
		errors = append(errors, binding.Error{
			FieldNames:     []string{"url"},
			Classification: "HostError",
			Message:        "The URL must be a valid hostname and path",
		})
	} else {
		parsedUrl, err := url.Parse(r.URL)
		if err != nil || parsedUrl.Host == "" {
			errors = append(errors, binding.Error{
				FieldNames:     []string{"url"},
				Classification: "HostError",
				Message:        "The URL appears to be in a bad format (2)",
			})
		} else {
			r.ParsedURL = parsedUrl
		}
	}

	return errors
}
