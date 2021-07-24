// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pepsi

import (
	"fmt"
	"reflect"
	"strings"
)

// Trie 模块负责实现前缀树 前缀树是动态路由的核心数据结构
// 静态路由表示Path地址与路由函数是一一对应的关系
// 动态路由表示有多个Path地址与路由函数是对应的关系
// URL格式: protocol://hostname[:port]/path/[;parameters][?query]#fragment
// 比如: https://leetcode-cn.com/problems/minimum-window-substring/
// 将每个path部分用 / 分割出来 得到 "problem" "minimum-window-substring"
// 我们就可以利用这些信息进行前缀树的匹配

// Params 定义动态路由参数key-value对
type Params map[string]string

// trieNode 前缀树节点定义
type trieNode struct {
	path      string
	children  []*trieNode
	wildChild bool    // 判断是否通配节点
	handler   Handler // 处理函数接口
	fullPath  string  // 到达该node的路径记录
}

// String 用于打印Trie树节点 方便debug
func (node *trieNode) String() string {
	handlerPtr := reflect.ValueOf(node.handler).Pointer()
	return fmt.Sprintf(
		"trieNode{path=%s, wildChild=%t, fullPath=%s, handler=0x%x}",
		node.path, node.wildChild, node.fullPath, handlerPtr,
	)
}

// addChild 将通配层末尾插入 其余头部插入
// 这样保证*静态*路由优先匹配
func (node *trieNode) addChild(child *trieNode) {
	if len(node.children) == 0 {
		node.children = []*trieNode{child}
		return
	}
	node.children = append(node.children, nil) // 扩充长度
	if child.wildChild {
		node.children[len(node.children)-1] = child
	} else {
		copy(node.children[1:], node.children[:])
		node.children[0] = child
	}
}

// parsePattern 将pattern拆分成规范的path列表
func parsePattern(pattern string) []string {
	paths := make([]string, 0)
	for _, path := range strings.Split(pattern, "/") {
		if path != "" {
			paths = append(paths, path)
			if path[0] == '*' {
				// 只匹配首个*通配符
				break
			}
		}
	}
	return paths
}

// insert 将一条路由pattern插入前缀树中
func (node *trieNode) insert(pattern string, handler Handler) {
	curNode := node
	paths := parsePattern(pattern)
	depth := len(paths)

	var recursion func(level int)
	recursion = func(level int) {
		if level == depth {
			if curNode.fullPath != "" {
				panic(fmt.Sprintf("pattern: %s has the same name with previous router!", pattern))
			}
			curNode.fullPath = pattern
			curNode.handler = handler
			return
		}
		path := paths[level]
		child := curNode.matchFirstchild(path)
		if child == nil {
			// 匹配失败 需要插入新node
			child := &trieNode{
				path:      path,
				wildChild: path[0] == ':' || path[0] == '*',
				handler:   nil,
			}
			curNode.addChild(child)
			curNode = child
		} else {
			curNode = child
		}
		recursion(level + 1)
	}

	recursion(0)
}

// query 从node查询给定路由pattern是否在Trie树中
// 若存在 返回该节点 否则返回nil
func (node *trieNode) query(pattern string) *trieNode {
	curNode := node
	paths := parsePattern(pattern)
	depth := len(paths)

	var recursion func(level int) *trieNode
	recursion = func(level int) *trieNode {
		if level == depth || strings.HasPrefix(curNode.path, "*") {
			if curNode.fullPath == "" {
				return nil
			}
			return curNode
		}
		// curNode可能有多个child与path匹配 那么每条路径都要试一遍
		// 一旦发现成功匹配 立即返回结果 否则返回nil给上层node
		path := paths[level]
		matchedChild := curNode.matchAllChild(path)
		for _, child := range matchedChild {
			curNode = child
			ok := recursion(level + 1)
			if ok != nil {
				return ok
			}
		}
		return nil
	}

	return recursion(0)
}

// matchFirstchild 返回首个匹配path的节点
func (node *trieNode) matchFirstchild(path string) *trieNode {
	for _, child := range node.children {
		if child.path == path {
			return child
		}
	}
	return nil
}

// matchAllChild 返回全部匹配成功的node
func (node *trieNode) matchAllChild(path string) []*trieNode {
	matched := make([]*trieNode, 0)
	for _, child := range node.children {
		if child.path == path || child.wildChild {
			matched = append(matched, child)
		}
	}
	return matched
}
