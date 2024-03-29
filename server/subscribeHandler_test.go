package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleSubscribe(t *testing.T) {
	s := NewServer("dummy", "dummy", false)
	// go s.Client.Start()
	rr := httptest.NewRecorder()

	// Do sub request
	topic := "randomTopic1"
	data, err := json.Marshal(SubRequestBody{
		Topic: topic,
	})
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, "", body)
	if err != nil {
		t.Error(err)
	}
	makeHTTPHandler(s.handleSubscribe)(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 but got %d", rr.Result().StatusCode)
	}

	// Test if Topics include the topics after sub
	expected := fmt.Sprintf("\"Subscribed to the topic %s\"", topic)
	b, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Error(err)
	}
	if strings.TrimRight(string(b), "\n") != expected {
		t.Errorf("expected %s but got %s", expected, string(b))
	}
}

func TestHandleSubscribeWithPrivateKey(t *testing.T) {
	s := NewServer("dummy", "./data", false)
	// go s.Client.Start()
	rr := httptest.NewRecorder()

	// Do sub request
	topic := "randomTopic1"
	data, err := json.Marshal(SubRequestBody{
		Topic: topic,
	})
	if err != nil {
		t.Error(err)
	}
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, "", body)
	if err != nil {
		t.Error(err)
	}
	makeHTTPHandler(s.handleSubscribe)(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 but got %d", rr.Result().StatusCode)
	}

	// Test if Topics include the topics after sub
	expected := fmt.Sprintf("\"Subscribed to the topic %s\"", topic)
	b, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Error(err)
	}
	if strings.TrimRight(string(b), "\n") != expected {
		t.Errorf("expected %s but got %s", expected, string(b))
	}
}
