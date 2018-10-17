package backend

import (
	"encoding/json"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	http.HandleFunc("/", handler)
}

// SpannerAPIListRequest
type SpannerAPIListRequest struct {
	ProjectID string `json:"projectId"`
	Instance  string `json:"instance"`
	Database  string `json:"database"`
	SQL       string `json:"sql"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var form SpannerAPIListRequest
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		log.Errorf(ctx, "%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	store, err := NewSpannerStore(ctx, form.ProjectID, form.Instance, form.Database)
	if err != nil {
		log.Errorf(ctx, "%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rows, err := store.ExecuteQuery(ctx, form.SQL)
	if err != nil {
		log.Errorf(ctx, "%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Errorf(ctx, "%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
