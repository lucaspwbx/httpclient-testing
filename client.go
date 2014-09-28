package client

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	token      string
	HttpClient *http.Client
}

type Person struct {
	Name string `json:"name"`
}

func (c *Client) do(method, path string, body io.Reader) ([]byte, int, error) {
	req, _ := http.NewRequest(method, path, body)
	resp, _ := c.HttpClient.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	return data, resp.StatusCode, nil
}

func (c *Client) Ping() error {
	_, _, err := c.do("GET", "/teste", nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) PingDois() ([]byte, int, error) {
	body, status, _ := c.do("GET", "/teste", nil)
	return body, status, nil
}

func (c *Client) PingTres() (Person, error) {
	var person Person
	body, _, _ := c.do("GET", "/teste", nil)
	err := json.Unmarshal(body, &person)
	if err != nil {
		return Person{}, err
	}
	return person, nil
}
