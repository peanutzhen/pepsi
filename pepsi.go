// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

// pepsi引擎提供路由功能/上下文封装
// pepsi引擎基于net/http

import (
	"net/http"
)

// Handler 定义路由处理函数接口
type Handler func(*Context)

// Engine Web后台引擎核心
type Engine struct {
	router *router // 路由表
}

// 为Engine的路由表添加路由项目
func (engine *Engine) addRouter(method string, pattern string, handler Handler) {
	engine.router.addRoute(method, pattern, handler)
}

// Get 添加Get方法路由
func (engine *Engine) Get(pattern string, handler Handler) {
	engine.addRouter("GET", pattern, handler)
}

// ServeHTTP 处理Http请求
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := newContext(w, req) // 构造上下文
	engine.router.handle(context) // 分发给router处理
}

// Run 启动pspsi Web服务
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// New 获取Engine实例
func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}
