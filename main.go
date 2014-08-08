package main

import (
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

func main() {
	LoadDatabase()
	StartAutoSave()

	m := martini.Classic()

	// API endpoints
	m.Get("/app/report/:host", GetReports)
	m.Post("/app/report", binding.Bind(InboundReport{}), FileReport)
	m.Get("/app/suggest/:prefix", SuggestHosts)

	// API endpoints for report actions
	m.Post("/app/update/:reportID/ack", binding.Bind(ActionMeta{}), Acknowledge)
	m.Post("/app/update/:reportID/fixed", binding.Bind(ActionMeta{}), Fixed)
	m.Post("/app/update/:reportID/reject", binding.Bind(ActionMeta{}), Reject)
	m.Post("/app/update/:reportID/delete", binding.Bind(ActionMeta{}), Delete)

	// Serve up the report page
	m.Get("/report/:host", func(req *http.Request, resp http.ResponseWriter) {
		http.ServeFile(resp, req, "public/report.html")
	})

	m.Run()
}
