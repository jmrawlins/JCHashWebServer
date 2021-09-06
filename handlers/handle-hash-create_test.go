package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestHashCreateHandler(t *testing.T) {
	t.Parallel()

	// Set up a multipart form POST
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("password", "angryMonkey")
	err := writer.Close()
	if err != nil {
		t.Fatal(err)
		return
	}
	reqMultipartForm, err := http.NewRequest("POST", "/hash", payload)
	if err != nil {
		t.Fatal(err)
	}
	reqMultipartForm.Header.Set("Content-Type", "multipart/form-data")

	// Set up a form-encoded data POST
	data := url.Values{}
	data.Set("password", "angryMonkey")
	reqFormUrlEncoded, err := http.NewRequest("POST", "/hash", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	reqFormUrlEncoded.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	reqs := []*http.Request{reqMultipartForm, reqFormUrlEncoded}

	setupTest := func(req *http.Request) (*HashDataStoreMock, *HashCreateHandler, *httptest.ResponseRecorder, *http.Request) {
		// Do the setup and request for all the following tests

		wg := &sync.WaitGroup{}
		ds := HashDataStoreMock{}
		ds.GetNextIdResult = struct {
			I uint64
			E error
		}{I: 42, E: nil}

		handler := NewHashCreateHandler(&ds, wg)

		rr := httptest.NewRecorder()

		return &ds, handler, rr, req
	}

	for _, request := range reqs {

		t.Run("Returns the expected response", func(t *testing.T) {
			t.Parallel()
			_, handler, rr, req := setupTest(request)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got '%v' expected '%v'",
					status, http.StatusOK)
			}

			expected := "42"
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got '%v' expected '%v'",
					rr.Body.String(), expected)
			}
		})

		t.Run("Calls the datastore with the expected arguments", func(t *testing.T) {
			t.Parallel()
			ds, handler, rr, req := setupTest(request)
			handler.ServeHTTP(rr, req)

			// TODO We shouldn't sleep inside a test. After refactoring the sleep into a channel
			// with a worker processing the jobs this can be removed. This is also
			// a race condition for our mock's values, but we made our sleep longer than the real code's
			// so it shouldn't bite us.
			time.Sleep(6 * time.Second)

			if ds.StoreHash_id != 42 {
				t.Errorf("Unexpected id argument to StoreHash: got '%v' expected '%v'", ds.StoreHash_id, 42)
				t.Fail()
			}

			expHash := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
			if ds.StoreHash_hash != expHash {
				t.Errorf("Unexpected hash argument to StoreHash: got '%v' expected '%v'", ds.StoreHash_hash, expHash)
				t.Fail()
			}
		})
	}
}
