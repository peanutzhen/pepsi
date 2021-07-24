// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

import (
	"log"
	"time"
)

// logger中间件提供日志的功能

func Logger() Handler {
	return func(c *Context) {
		t := time.Now()
		c.NextHandler()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
