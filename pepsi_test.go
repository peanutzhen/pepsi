// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

const (
	LOCALHOST = "http://localhost:9999"
	PORT      = ":9999"
)

var engine = New()


func TestEngine_Get(t *testing.T) {
	s := "Hello world!"
	log.Println(engine)
	engine.Get("/hello_world", func(context *Context) {
		context.String(http.StatusOK, s)
	})
	go engine.Run(":9999")

	rsp, err := http.Get("http://localhost:9999/hello_world")
	if err != nil {
		t.Error(err)
		return
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if rsp.StatusCode != 200 {
		t.Error(rsp.StatusCode)
	}
	if string(body) != s {
		t.Error(string(body))
	}

	rsp, err = http.Get("http://localhost:9999/not_found")
	if err != nil {
		t.Error(err)
		return
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != 404 {
		t.Error(rsp.StatusCode)
	}
	body, err = ioutil.ReadAll(rsp.Body)
	if string(body) != "404 NOT FOUND: /not_found\n" {
		t.Error(string(body))
	}
}
