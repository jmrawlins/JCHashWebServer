package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHashGetHandler(t *testing.T) {
	t.Parallel()

	t.Run("Called with an existing hash id", func(t *testing.T) {
		t.Parallel()

		// Do the setup and request for all the following tests
		req, err := http.NewRequest("GET", "/42", nil)
		if err != nil {
			t.Fatal(err)
		}

		ds := HashDataStoreMock{}
		ds.GetHashResult = struct {
			H string
			E error
		}{H: "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==", E: nil}

		handler := NewHashGetHandler(&ds)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		t.Run("Returns the expected response", func(t *testing.T) {
			if status := rr.Code; status != int(http.StatusOK) {
				t.Errorf("handler returned wrong status code: got '%v' expected %v",
					status, http.StatusOK)
			}

			expected := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got '%v' expected %v",
					rr.Body.String(), expected)
			}

		})

		t.Run("Calls the datastore with the expected arguments", func(t *testing.T) {
			// calls StoreHash with the expected id and hash
			if ds.GetHash_id != 42 {
				t.Errorf("Expected to call GetHash with %v, got '%v'", 42, ds.GetHash_id)
				t.Fail()
			}
		})
	})

	t.Run("Called with nonexistent hash id", func(t *testing.T) {
		t.Parallel()

		// Do the setup and request for all the following tests
		req, err := http.NewRequest("GET", "/42", nil)
		if err != nil {
			t.Fatal(err)
		}

		ds := HashDataStoreMock{}
		ds.GetHashResult = struct {
			H string
			E error
		}{H: "", E: errors.New("Some Error")}

		handler := NewHashGetHandler(&ds)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		t.Run("Returns the expected response", func(t *testing.T) {
			expected := http.StatusNotFound
			if status := rr.Code; status != expected {
				t.Errorf("handler returned wrong status code: got '%v' expected %v",
					status, expected)
			}
		})

		t.Run("Calls the datastore with the expected arguments", func(t *testing.T) {
			if ds.GetHash_id != 42 {
				t.Fail()
			}
		})
	})

}
