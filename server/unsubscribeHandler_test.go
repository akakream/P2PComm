package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleUnsubscribe(t *testing.T) {
	s := NewServer("dummy", "libp2p")
	// go s.Client.Start()

	topic := "randomTopic1"
	// Subscribe to a topic
	helperSubscribe(topic, t, s)

	rr := httptest.NewRecorder()

	data, err := json.Marshal(SubRequestBody{
		Topic: topic,
	})
	if err != nil {
		t.Errorf("Could not marhsall subscription data %s", err)
	}
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, "", body)
	if err != nil {
		t.Error(err)
	}

	// Test Unsub
	makeHTTPHandler(s.handleUnsubscribe)(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 but got %d", rr.Result().StatusCode)
	}
	defer rr.Result().Body.Close()

	// Test if Topics include the topics after sub
	jsonBody := UnsubResponseBody{}
	if err := json.NewDecoder(rr.Result().Body).Decode(&jsonBody); err != nil {
		t.Errorf("could not read response body: %s", err)
	}
	if jsonBody.Topic != topic {
		t.Errorf("expected the topic (%s) but got %s topics", topic, jsonBody.Topic)
	}
	if jsonBody.Messages == nil {
		t.Errorf("expected empty array messages but got null")
	}
}

func helperSubscribe(topic string, t *testing.T, s *Server) {
	rr := httptest.NewRecorder()

	data, err := json.Marshal(SubRequestBody{
		Topic: topic,
	})
	if err != nil {
		t.Errorf("Could not marhsall subscription data %s", err)
	}
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, "", body)
	if err != nil {
		t.Error(err)
	}

	// Sub
	makeHTTPHandler(s.handleSubscribe)(rr, req)
}
