package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleListSubscribedTopicsWhenEmpty(t *testing.T) {
	s := NewServer("dummy", "libp2p")
	// go s.Client.Start()

	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Error(err)
	}
	makeHTTPHandler(s.handleListSubscribedTopics)(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 but got %d", rr.Result().StatusCode)
	}
	defer rr.Result().Body.Close()

	// Test if Topics in response is empty
	expectedLen := 0
	jsonBody := ListTopicsRequestBody{}
	if err := json.NewDecoder(rr.Result().Body).Decode(&jsonBody); err != nil {
		t.Errorf("could not read response body: %s", err)
	}
	if len(jsonBody.Topics) != expectedLen {
		t.Errorf("expected empty topics array but got %d topics", len(jsonBody.Topics))
	}
}

func TestHandleListSubscribedTopicsWhenOne(t *testing.T) {
	s := NewServer("dummy", "libp2p")
	// go s.Client.Start()

	topic := "randomTopic1"
	// Subscribe to a topic
	helperSubscribe(topic, t, s)

	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Error(err)
	}
	makeHTTPHandler(s.handleListSubscribedTopics)(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 but got %d", rr.Result().StatusCode)
	}
	defer rr.Result().Body.Close()

	// Test if Topics include the topics after sub
	expected := "randomTopic1"
	jsonBody := ListTopicsRequestBody{}
	if err := json.NewDecoder(rr.Result().Body).Decode(&jsonBody); err != nil {
		t.Errorf("could not read response body: %s", err)
	}
	if jsonBody.Topics[0] != expected {
		t.Errorf("expected %s in the topics but it is not there", jsonBody.Topics[0])
	}
}
