package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ServerError struct {
	Err    string `json:"err"`
	Status int    `json:"status"`
}

func TestHandlePublishWhenNotSubbed(t *testing.T) {
	s := NewServer("dummy", "libp2p")
	// go s.Client.Start()

	topic := "randomTopic1"
	expectedError := "client is not subscribed to the topic"
	expectedStatus := 500

	rr := httptest.NewRecorder()

	// Do sub request
	message := "randomMessage1"
	data, err := json.Marshal(PubRequestBody{
		Topic:   topic,
		Message: message,
	})
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, "", body)
	if err != nil {
		t.Error(err)
	}
	makeHTTPHandler(s.handlePublish)(rr, req)

	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 but got %d", rr.Result().StatusCode)
	}

	// Test if the return values are right
	got := ServerError{}
	if err := json.NewDecoder(rr.Result().Body).Decode(&got); err != nil {
		t.Errorf("could not read response body: %s", err)
	}
	if got.Err != expectedError {
		t.Errorf("expected the error (%s) but got the error (%s)", expectedError, got.Err)
	}
	if got.Status != expectedStatus {
		t.Errorf("expected %d but got %d", expectedStatus, got.Status)
	}
}

func TestHandlePublishWhenSubbed(t *testing.T) {
	s := NewServer("dummy", "libp2p")
	// go s.Client.Start()

	// Subscribe to topic
	topic := "randomTopic1"
	// Subscribe to a topic
	helperSubscribe(topic, t, s)

	rr := httptest.NewRecorder()

	// Do sub request
	message := "randomMessage1"
	data, err := json.Marshal(PubRequestBody{
		Topic:   topic,
		Message: message,
	})
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, "", body)
	if err != nil {
		t.Error(err)
	}
	makeHTTPHandler(s.handlePublish)(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 but got %d", rr.Result().StatusCode)
	}

	// Test if the return values are right
	got := PubResponseBody{}
	if err := json.NewDecoder(rr.Result().Body).Decode(&got); err != nil {
		t.Errorf("could not read response body: %s", err)
	}
	if got.Topic != topic {
		t.Errorf("expected the topic (%s) but got the topic (%s)", topic, got.Topic)
	}
	if got.Messages == nil {
		t.Error("expected messages not to be emtpy")
	}
	if len(got.Messages) != 1 {
		t.Errorf("expected 1 message but got %d messages", len(got.Messages))
	}
	if got.Messages[0] != message {
		t.Errorf("expected the message (%s) but got the message (%s)", message, got.Messages)
	}
}
