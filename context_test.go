// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	engine := CreateEngine()
	engine.GET("/data", func(context *Context) {
		context.Data(http.StatusOK, []byte("Hello data"))
	})
	engine.GET("/html", func(context *Context) {
		context.HTML(http.StatusOK, "<h1>Hello world</h1>")
	})
	engine.GET("/json", func(context *Context) {
		context.JSON(http.StatusOK, map[string]string{"hello": "world"})
	})
	engine.GET("/fail", func(context *Context) {
		context.Fail(http.StatusInternalServerError, "Internal Serber Error")
	})
	ts := httptest.NewServer(engine)
	defer ts.Close()

	res1, _ := http.Get(fmt.Sprintf("%s/data", ts.URL))
	res2, _ := http.Get(fmt.Sprintf("%s/html", ts.URL))
	res3, _ := http.Get(fmt.Sprintf("%s/json", ts.URL))
	res4, _ := http.Get(fmt.Sprintf("%s/fail", ts.URL))

	//resp1, _ := ioutil.ReadAll(res1.Body)
	//resp2, _ := ioutil.ReadAll(res2.Body)
	//resp3, _ := ioutil.ReadAll(res3.Body)
	//resp4, _ := ioutil.ReadAll(res4.Body)

	assert.Equal(t, http.StatusOK, res1.StatusCode)
	assert.Equal(t, http.StatusOK, res2.StatusCode)
	assert.Equal(t, http.StatusOK, res3.StatusCode)
	assert.Equal(t, http.StatusInternalServerError, res4.StatusCode)

	//assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}
