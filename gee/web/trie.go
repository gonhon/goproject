package web

import "strings"

type node struct {
	//待匹配的路由 /a/:b
	pattern string
	//路由中的一部分 :b
	part string
	//子节点 [c,d]
	children []*node
	//是否模糊匹配  : 或 * 时为true
	isWild bool
}

// 匹配第一个节点
func (n *node) matchChildren(part string) *node {
	for _, children := range n.children {
		if children.part == part || children.isWild {
			return children
		}
	}
	return nil
}

// 匹配所有的
func (n *node) matchChildrens(part string) []*node {
	// childrens := make([]*node, len(n.children))//error
	childrens := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			childrens = append(childrens, child)
		}
	}
	return childrens
}

//将url按照parts插入到tire中
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		//最后插入pattern
		n.pattern = pattern
		return
	}
	part := parts[height]
	//找到第一个匹配的结点
	child := n.matchChildren(part)
	//没有匹配的结点就作为当前元素的子节点
	if child == nil {
		// :|* isWild为true
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	//递归插入
	child.insert(pattern, parts, height+1)
}

//根据parts从根逐个匹配
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	//根据part找到所有匹配的节点
	children := n.matchChildrens(part)

	for _, child := range children {
		//逐个递归查询
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
