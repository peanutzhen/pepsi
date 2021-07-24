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

func TestRecovery(t *testing.T) {
	engine := CreateEngine()
	engine.GET("/recovery", func(context *Context) {
		ar := []int{1}
		context.JSON(http.StatusOK, ar[1]) // 引发数组越界
	})
	ts := httptest.NewServer(engine)
	defer ts.Close()

	res, _ := http.Get(fmt.Sprintf("%s/recovery", ts.URL))
	body, _ := ioutil.ReadAll(res.Body)

	assert.Contains(t, string(body), "message")
}
