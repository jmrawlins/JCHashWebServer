package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jmrawlins/JCHashWebServer/datastore"
)

func TestStatsHandler(t *testing.T) {
	t.Parallel()

	t.Run("Called without ?all", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest("GET", "/stats", nil)
		if err != nil {
			t.Fatal(err)
		}

		ds := datastore.StatsDataStoreMock{}
		ds.GetUriStatsResults.S = datastore.RequestStats{URI: "/hash", Total: 42, Average: 10.5}
		handler := NewStatsHandler(&ds)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		// Check that it called the datastore with the right request
		if ds.GetUriStats_uri != "/hash" {
			t.Errorf("handler called GetUriStats with the wrong uri: got '%v', expected '%v'",
				ds.GetUriStats_uri, "/hash")
			t.Fail()
		}

		// Check that the response code is correct
		if rr.Code != http.StatusOK {
			t.Errorf("handler returned unexpected response code: got '%v', expected '%v'",
				rr.Code, http.StatusOK)
		}

		// Check that the response body is correct
		expected := `{"request":"/hash","total":42,"average":10.5}`
		if strings.Trim(rr.Body.String(), "\n") != expected {
			t.Errorf("handler returned unexpected body: got '%v' expected '%v'",
				rr.Body.String(), expected)
		}
	})

	t.Run("Called with ?all", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest("GET", "/stats?all", nil)
		if err != nil {
			t.Fatal(err)
		}

		ds := datastore.StatsDataStoreMock{}
		ds.GetStatsResult.S = ``
		handler := NewStatsHandler(&ds)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		// Check that it called the datastore with the right request
		if ds.GetUriStats_uri != "/hash" {
			t.Errorf("handler called GetUriStats with the wrong uri: got '%v', expected '%v'",
				ds.GetUriStats_uri, "/hash")
			t.Fail()
		}

		// Check that the response code is correct
		if rr.Code != http.StatusOK {
			t.Errorf("handler returned unexpected response code: got '%v', expected '%v'",
				rr.Code, http.StatusOK)
		}

		// Check that the response body is correct
		expected := `{"request":"/hash","total":42,"average":10.5}`
		if strings.Trim(rr.Body.String(), "\n") != expected {
			t.Errorf("handler returned unexpected body: got '%v' expected '%v'",
				rr.Body.String(), expected)
		}
	})

}
