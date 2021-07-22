// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestEngine_Get(t *testing.T) {
	engine := New()
	engine.Get("/hello_world", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello world!")
	})
	go engine.Run(":9999")

	rsp, err := http.Get("http://localhost:9999/hello_world")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(*rsp)
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if rsp.StatusCode != 200 {
		t.Error(rsp.StatusCode)
	}
	if string(body) != "Hello world!" {
		t.Error(string(body))
	}

	rsp, err = http.Get("http://localhost:9999/not_found")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(*rsp)
	defer rsp.Body.Close()
	if rsp.StatusCode != 404 {
		t.Error(rsp.StatusCode)
	}
}
