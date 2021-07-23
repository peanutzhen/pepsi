// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

import (
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

// 我们可以使用DeepEqual判断任意两种类型实例是否相等
func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	_, ps := r.getRoute("GET", "/hello/peanutzhen")

	if ps["name"] != "peanutzhen" {
		t.Fatalf("name %s should be equal to 'peanutzhen'", ps["name"])
	}

	_, ps = r.getRoute("GET", "/hello/b")
	if ps["name"] != "b" {
		t.Fatalf("name %s should be equal to 'b'", ps["name"])
	}

	_, ps = r.getRoute("GET", "/hello/b/c")
	if _, ok := ps["name"]; ok {
		t.Fatal("Should not get param 'name'!")
	}

	_, ps = r.getRoute("GET", "/assets/js/i_hate_js.js")
	if ps["filepath"] != "js/i_hate_js.js" {
		t.Fatal("Should get param 'filename' value js/i_hate_js.js!")
	}
}
