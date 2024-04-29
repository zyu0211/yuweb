package yu

import "strings"

/*
   前缀树节点
*/
type Node struct {
    path string         //完整路由地址，例如：/p/:lang
    part string         //路由地址的一部分，例如：:lang
    children []*Node    //子节点
    isWild bool         //是否为模糊匹配，若part中含有 : 或 *，则为true
}

/*
    返回第一个匹配成功的节点
*/
func (node *Node) matchChild(part string) *Node {

    for _, child := range node.children {
        if child.part == part || child.isWild {
            return child
        }
    }

    return nil
}

/*
    返回所有匹配成功的节点
*/
func (node *Node) matchChildren(part string) []*Node {

    nodes := make([]*Node ,0)

    for _, child := range node.children {
        if child.part == part || child.isWild {
            nodes = append(nodes, child)
        }
    }

    return nodes
}

/*
    插入节点
*/
func (node *Node) insert(path string, parts []string, height int) {

    if len(parts) == height {
        node.path = path    //只有最后一个节点才会存储完整的路由地址 path
        return
    }

    part := parts[height]
    child := node.matchChild(part)
    if child == nil {
        child = &Node{
            part: part,
            isWild: part[0] == ':' || part[0] == '*',
        }
        node.children = append(node.children, child)
    }

    child.insert(path, parts, height + 1)
}

/*
    查询节点
*/
func (node *Node) search(parts []string, height int) *Node{

    if len(parts) == height || strings.HasPrefix(node.part, "*") {
        if node.path == "" {
            return nil
        }
        return node
    }

    part := parts[height]
    children := node.matchChildren(part)

    for _, child := range children {
        result := child.search(parts, height + 1)
        if result != nil {
            return result
        }
    }

    return nil
}