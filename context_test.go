// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestContext_Data(t *testing.T) {
	data := []byte("I am peanutzhen!.")
	log.Println(engine)
	engine.Get("/data", func(context *Context) {
		context.Data(http.StatusOK, data)
	})
	go engine.Run(PORT)

	rsp, err := http.Get(LOCALHOST + "/data")
	defer rsp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if string(body) != string(data) {
		t.Error(string(body))
	}
}

func TestContext_HTML(t *testing.T) {
	html := "<h1>TestContext_HTML</h1>"
	log.Println(engine)
	engine.Get("/html", func(context *Context) {
		context.HTML(http.StatusOK, html)
	})
	go engine.Run(PORT)

	rsp, err := http.Get(LOCALHOST + "/html")
	defer rsp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if string(body) != html {
		t.Errorf("Actual: %s\nExpected: %s\n", string(body), html)
	}
}

func TestContext_JSON(t *testing.T) {
	obj := map[string]interface{}{
		"username": "peanutzhen",
		"age":      20,
	}
	log.Println(engine)
	engine.Get("/json", func(context *Context) {
		context.JSON(http.StatusOK, obj)
	})
	go engine.Run(PORT)

	rsp, err := http.Get(LOCALHOST + "/json")
	defer rsp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	body, err := ioutil.ReadAll(rsp.Body)
	jsonObj, err := json.Marshal(obj)
	if string(body) != string(jsonObj)+"\n" {
		t.Errorf("Actual: %s\nExpected: %s\n", string(body), string(jsonObj))
	}
}
