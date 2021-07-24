// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func onlyForV2() Handler {
	return func(c *Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func TestRouterGroup(t *testing.T) {
	engine := CreateEngine()
	engine.GET("/greeting", func(c *Context) {
		c.HTML(http.StatusOK, "<h1>OK</h1>")
	})
	v2 := engine.ForkGroup("/v2")
	v2.AddMiddlewares(onlyForV2())
	v2.GET("/hello/:name", func(c *Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	ts := httptest.NewServer(engine)
	defer ts.Close()

	res1, _ := http.Get(fmt.Sprintf("%s/greeting", ts.URL))
	rsp1, _ := ioutil.ReadAll(res1.Body)
	log.Println(string(rsp1))

	res2, _ := http.Get(fmt.Sprintf("%s/v2/hello/peanutzhen", ts.URL))
	rsp2, _ := ioutil.ReadAll(res2.Body)
	log.Println(string(rsp2))

	assert.Equal(t, "<h1>OK</h1>", string(rsp1))
	assert.Equal(t, http.StatusInternalServerError, res2.StatusCode)
}
