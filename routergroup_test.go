// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func Logger() Handler {
	return func(c *Context) {
		t := time.Now()
		c.NextHandler()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func onlyForV2() Handler {
	return func(c *Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func TestRouterGroup(t *testing.T) {
	log.Printf("RouteGroup Engine: %+v\n", engine)
	engine.AddMiddlewares(Logger())
	engine.Get("/greeting", func(c *Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := engine.ForkGroup("/v2")
	v2.AddMiddlewares(onlyForV2())
	{
		log.Printf("v2: %+v\n", v2)
		v2.Get("/hello/:name", func(c *Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	go engine.Run(PORT)

	rsp, _ := http.Get(LOCALHOST + "/greeting")
	defer rsp.Body.Close()
	body, _ := ioutil.ReadAll(rsp.Body)
	fmt.Println(string(body))

	rsp, _ = http.Get(LOCALHOST + "/v2/hello/peanutzhen")
	defer rsp.Body.Close()
	body, _ = ioutil.ReadAll(rsp.Body)
	fmt.Println(string(body))
}
