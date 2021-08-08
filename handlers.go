package covidstats

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

// HandleFindYearStats is the handler that returns the confirmed cases
// grouped by date for a given year.
func (s *Server) HandleFindYearStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	s.logger.Info("HandleFindYearStats")
	if r.Method == http.MethodOptions {
		return
	}

	vars := mux.Vars(r)
	yr := vars["year"]
	year, _ := strconv.Atoi(yr)

	if year <= 0 {
		year = time.Now().Year()
	}

	cases, findErr := s.casesService.FindByYear(r.Context(), year)
	if findErr != nil {
		s.logger.WithFields(log.Fields{
			"year": year,
			"vars": vars,
		}).
			WithError(findErr).
			Error("FindByYear failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(cases); err != nil {
		s.logger.WithError(err).Error("encoding json response failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
