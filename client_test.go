package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type FakeRoundTripper struct {
	message  string
	status   int
	header   map[string]string
	requests []*http.Request
}

func newTestClient(rt *FakeRoundTripper) *Client {
	client := &Client{
		token:      "foobar",
		HttpClient: &http.Client{Transport: rt},
	}
	return client
}

func (rt *FakeRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	body := strings.NewReader(rt.message)
	rt.requests = append(rt.requests, r)
	res := &http.Response{
		StatusCode: rt.status,
		Body:       ioutil.NopCloser(body),
		Header:     make(http.Header),
	}
	for k, v := range rt.header {
		res.Header.Set(k, v)
	}
	return res, nil
}

func (rt *FakeRoundTripper) Reset() {
	rt.requests = nil
}

func TestPing(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	err := client.Ping()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPingDois(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "lucas", status: http.StatusOK}
	client := newTestClient(fakeRT)
	body, status, _ := client.PingDois()
	if string(body) != "lucas" {
		t.Errorf("Error, got %s", string(body))
	}
	if status != 200 {
		t.Errorf("Error")
	}
}

func TestPingTres(t *testing.T) {
	body := `{
	  "name": "Lucas"
  }`
	var expected Person
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshalling")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.PingTres()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error")
	}
}
