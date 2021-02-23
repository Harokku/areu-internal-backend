package main

import (
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
		fmt.Println(string(body)) // => pong
	}

	utils.AssertEqual(t, nil, err, "app.test")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	utils.AssertEqual(t, "pong", string(body), "Response body")
}
