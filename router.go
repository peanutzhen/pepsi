// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

// router模块负责路由功能

import (
	"net/http"
	"strings"
)

// router 路由模块类
type router struct {
	// 处理函数表 我们将所有的handler和trieNode绑定
	// 从而在匹配路由时实现动态路由 注意到每个method维护一颗Trie树
	handlers map[string]*trieNode
}

// router 添加路由
func (r *router) addRoute(method string, pattern string, handler Handler) {
	if _, ok := r.handlers[method]; !ok {
		r.handlers[method] = &trieNode{}
	}
	r.handlers[method].insert(pattern, handler)
}

func (r *router) getRoute(method string, pattern string) (Handler, Params) {
	if _, ok := r.handlers[method]; !ok {
		return nil, nil
	}
	targerNode := r.handlers[method].query(pattern)
	if targerNode != nil {
		// 路由匹配成功 正在构造动态路由参数
		ps := make(Params)
		searchPaths := parsePattern(pattern)
		for i, path := range parsePattern(targerNode.fullPath) {
			if path[0] == ':' {
				ps[path[1:]] = searchPaths[i]
			} else if path[0] == '*' && len(path) > 1 {
				ps[path[1:]] = strings.Join(searchPaths[i:], "/")
				break
			}
		}
		return targerNode.handler, ps
	}
	return nil, nil
}

// handle 负责处理 Context
// 解析 Context 请求path 找到对应的handler 并填充好动态路由的数据
func (r *router) handle(context *Context) {
	handler, params := r.getRoute(context.Method, context.Path)
	if handler != nil {
		context.Params = params // 填充动态路由参数
		// 填充实际业务Handler
		context.handlers = append(context.handlers, handler)
	} else {
		context.handlers = append(context.handlers, func(context *Context) {
			context.String(http.StatusNotFound, "404 NOT FOUND: %s\n", context.Path)
		})
	}
	context.NextHandler()
}

// 构造 router 实例
func newRouter() *router {
	return &router{
		handlers: make(map[string]*trieNode),
	}
}
