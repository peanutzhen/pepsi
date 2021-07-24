// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_200OK(t *testing.T) {
	engine := CreateEngine()
	engine.GET("/hello", func(context *Context) {
		context.HTML(http.StatusOK, "<h1>Hello world</h1>")
	})
	ts := httptest.NewServer(engine)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/hello", ts.URL))
	if err != nil {
		fmt.Println(err)
	}
	resp, _ := ioutil.ReadAll(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func Test_404NotFound(t *testing.T) {
	engine := CreateEngine()
	ts := httptest.NewServer(engine)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/not_found", ts.URL))
	if err != nil {
		fmt.Println(err)
	}
	resp, _ := ioutil.ReadAll(res.Body)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Contains(t, string(resp), "404 NOT FOUND")
}
