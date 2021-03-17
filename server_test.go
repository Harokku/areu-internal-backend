package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/utils"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	app := fiberApp()
	var body []byte

	req := httptest.NewRequest("GET", "/ping", nil)
	resp, err := app.Test(req)

	if resp.StatusCode == 200 {
		body, _ = ioutil.ReadAll(resp.Body)
	}

	utils.AssertEqual(t, nil, err, "app.test")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	utils.AssertEqual(t, "pong", string(body), "Response body")
}

func TestLogin(t *testing.T) {
	app := fiberApp()
	var body []byte
	postBody := map[string]string{
		"identity": "user",
		"password": "pass",
	}
	marshPostBody, _ := json.Marshal(postBody)
	fmt.Printf("marshalled body: %s\n", marshPostBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(marshPostBody))
	req.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode == 200 {
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Printf("Body %s", body)
	}
}
