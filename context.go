// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

// context模块封装了一次Http请求与响应，提供了强大的接口
// 方便用户快速构造返回体

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Context 定义一次http交互的信息
type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	// 定义request所含信息
	Path   string // 路由
	Method string // 请求方法(GET/POST/PUT)
	Params Params // 动态路由参数

	handlers HandlerChain
	i        int // 表示执行到Handlers[i]

	// 定义response所含信息
	StatusCode int // 响应http状态码
}

// NextHandler 告诉 Context 进入下一个处理逻辑
func (context *Context) NextHandler() {
	context.i++
	for context.i < len(context.handlers) {
		context.handlers[context.i](context)
		context.i++
	}
}

// PostForm 返回 Context 中 post 数据对应 key 的 value。
// 若 value 不存在，则返回 ""
func (context *Context) PostForm(key string) string {
	return context.Req.FormValue(key)
}

// Query 返回 Context 中 get 数据对应 key 的 value。
// 若 value 不存在，则返回 ""
func (context *Context) Query(key string) string {
	return context.Req.URL.Query().Get(key)
}

// SetStatus 设置http返回码
func (context *Context) SetStatus(code int) {
	context.StatusCode = code
	context.Writer.WriteHeader(code)
}

// SetHeader 设置http响应的header
func (context *Context) SetHeader(key string, value string) {
	context.Writer.Header().Set(key, value)
}

// Data 返回字节流数据
func (context *Context) Data(code int, data []byte) {
	context.SetStatus(code)
	context.Writer.Write(data)
}

// String 快速构造返回类型为 text/plain 的 Body
func (context *Context) String(code int, format string, values ...interface{}) {
	context.SetHeader("Content-Type", "text/plain")
	context.SetStatus(code)
	context.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 快速构造返回类型为 application/json 的 Body
func (context *Context) JSON(code int, obj interface{}) {
	context.SetHeader("Content-Type", "application/json")
	context.SetStatus(code)
	encoder := json.NewEncoder(context.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(context.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// HTML 通过HTML内容快速构造返回Body
func (context *Context) HTML(code int, html string) {
	context.SetHeader("Content-Type", "text/html")
	context.SetStatus(code)
	context.Writer.Write([]byte(html))
}

func (context *Context) Fail(code int, err string) {
	context.i = len(context.handlers)
	context.JSON(code, map[string]interface{}{"message": err})
}

// 构造 Context 实例
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		i: -1,
	}
}
