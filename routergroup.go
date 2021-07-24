// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

// routergroup模块负责分组控制
// 分组即对某个路由下的所有路由进行划分
// 分组的好处是实现了功能的划分 如 鉴权功能划分 校验功能划分
// 例如 /admin 下的所有路由访问前均需鉴权
// /api 下的路由需参数校验功能
// 这些功能可以看作一堆插件/中间件 通过Handler实现

type RouterGroup struct {
	prefix      string
	middlewares HandlerChain // 分组下需要使用的中间件
	engine      *Engine      // 提供指向Engine的能力
}

// ForkGroup 创建新分组
// 分组控制必须有嵌套的能力 也即新分组的控制必须继承父分组
// 通过嵌套父Group的prefix实现
func (group *RouterGroup) ForkGroup(prefix string) *RouterGroup {
	childGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: group.engine,
	}
	group.engine.groups = append(group.engine.groups, childGroup)
	return childGroup
}

// AddMiddlewares 为 RouterGroup 插入插件
func (group *RouterGroup) AddMiddlewares(middlewares ...Handler) {
	group.middlewares = append(group.middlewares, middlewares...)
}
