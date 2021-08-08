package covidstats

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// HandleFindYearStats is the handler that returns the confirmed cases
// grouped by date for a given year.
func (s *Server) HandleFindYearStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	yr := vars["year"]
	year, _ := strconv.Atoi(yr)
	if year <= 0 {
		year = time.Now().Year()
	}

	cases, findErr := s.casesService.FindByYear(r.Context(), year)
	if findErr != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(cases); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
