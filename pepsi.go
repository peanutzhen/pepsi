// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

import (
	"fmt"
	"net/http"
)

// Handler 定义路由业务处理逻辑
type Handler func(http.ResponseWriter, *http.Request)

// Engine Web后台引擎核心，基于net/http。
type Engine struct {
	router map[string]Handler // 路由表
}

// 为Engine的路由表添加路由项目
func (engine *Engine) addRouter(method string, pattern string, handler Handler) {
	key := method + "_" + pattern
	engine.router[key] = handler
}

// Get 添加Get方法路由
func (engine *Engine) Get(pattern string, handler Handler) {
	engine.addRouter("GET", pattern, handler)
}

// ServeHTTP 处理Http请求
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "_" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND: #{key}\n")
	}
}

// Run 启动pspsi Web服务
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// New 获取Engine实例
func New() *Engine {
	return &Engine{router: make(map[string]Handler)}
}
