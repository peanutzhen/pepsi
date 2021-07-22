// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

// router模块负责路由功能

import (
	"log"
	"net/http"
)

// router 路由模块类
type router struct {
	handlers map[string]Handler // 处理函数表
}

// router 添加路由
func (r *router) addRoute(method string, pattern string, handler Handler) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "_" + pattern
	r.handlers[key] = handler
}

// router 负责处理 Context
func (r *router) handle(context *Context) {
	key := context.Method + "_" + context.Path
	if handler, isExists := r.handlers[key]; isExists {
		handler(context)
	} else {
		context.String(http.StatusNotFound, "404 NOT FOUND: %s\n", context.Path)
	}
}

// 构造 router 实例
func newRouter() *router {
	return &router{
		handlers: make(map[string]Handler),
	}
}
