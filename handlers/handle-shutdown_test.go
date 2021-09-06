package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestShutdownHandler(t *testing.T) {
	t.Parallel()

	t.Run("Signals the shutdown channel", func(t *testing.T) {
		t.Parallel()

		// Do the setup and request for all the following tests
		req, err := http.NewRequest("POST", "/shutdown", nil)
		if err != nil {
			t.Fatal(err)
		}

		shutdownChannel := make(chan struct{})
		handler := NewShutdownHandler(shutdownChannel)

		rr := httptest.NewRecorder()
		// writes to channel, so we need to be around to receive it
		go handler.ServeHTTP(rr, req)

		timeoutChannel := time.After(500 * time.Millisecond)

		select {
		case <-timeoutChannel:
			t.Error("Got timeout instead of receive on shutdown channel")
			t.Fail()
		case <-shutdownChannel:
			// Pass
		}

	})

}
