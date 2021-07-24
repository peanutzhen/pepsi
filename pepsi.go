// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

// pepsi引擎提供路由功能/上下文封装
// pepsi引擎基于net/http

import (
	"log"
	"net/http"
	"strings"
)

type (
	// Handler 定义路由处理函数接口
	Handler func(*Context)
	// HandlerChain 在处理http request时 我们可以通过添加Handler
	// 实现功能上的拓展 即插件模式 这样我们需要Chain保存需调用的Handler
	HandlerChain []Handler
)

// Engine Web后台引擎核心
type Engine struct {
	*RouterGroup         // Engine 组合了分组控制的能力
	router       *router // 路由表
	// 保存指向所有分组控制的指针 用于获取请求path的所有中间件
	groups []*RouterGroup
}

// 为Engine的路由表添加路由项目
// 实际上 这是RouterGroup负责的
func (group *RouterGroup) addRouter(method string, pattern string, handler Handler) {
	pattern = group.prefix + pattern
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// Get 添加Get方法路由
func (group *RouterGroup) Get(pattern string, handler Handler) {
	group.addRouter("GET", pattern, handler)
}

// Post 添加Post方法路由
func (group *RouterGroup) Post(pattern string, handler Handler) {
	group.addRouter("POST", pattern, handler)
}

// ServeHTTP 处理Http请求
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req) // 构造上下文
	for _, group := range engine.groups {
		// 为请求URL添加中间件
		if strings.HasPrefix(c.Path, group.prefix) {
			c.handlers = append(c.handlers, group.middlewares...)
		}
	}
	engine.router.handle(c) // 交给router响应这个上下文
}

// Run 启动pspsi Web服务
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// New 获取Engine实例
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}
