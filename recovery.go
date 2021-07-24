// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// recovery 模块负责框架在运行时抛出的错误恢复
// 基于runtime的能力获取错误调用堆栈

func Recovery() Handler {
	return func(context *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				context.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		context.NextHandler()
	}
}

func trace(message string) string {
	var pcstack [32]uintptr
	n := runtime.Callers(3, pcstack[:])

	// Using Builder optimize speed.
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcstack[:n] {
		function := runtime.FuncForPC(pc)
		file, line := function.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
