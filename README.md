# Pepsi Web Framework
Pepsi是跟随[7days-golang](https://github.com/geektutu/7days-golang)开发的简单Web框架，并在此基础上增添了自己的想法，也算是学习了Gin框架的原理。Pepsi取名灵感源自于女朋友网名叫pepsi。

## Feature

- 前缀树`Trie`作为路由功能
- `handlerChain`实现与业务逻辑无关的功能拓展
- 提供路由分组中间件控制，实现基于基底路径的逻辑划分。如`/api`下提供A中间件，而`user`提供B中间件功能。

## Further

- 使用自己实现的`net/http`库作为基础网络API

## Installation

环境需安装Go 1.16以上，然后执行：

```sh
$ go get -u github.com/peanutzhen/pepsi
```

## Quick start

```sh
$ cat example.go
```

`example.go`输出如下：

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/peanutzhen/pepsi"
)

func main() {
	r := pepsi.CreateEngine()
	r.GET("/hello", func(c *pepsi.Context) {
		c.HTML(http.StatusOK, "<h1>Hello world!</h1>")
	})
	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
```

运行HTTP服务进程：

```sh
$ go run example.go
```
