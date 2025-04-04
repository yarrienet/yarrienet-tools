package htmlhelper

import (
    "golang.org/x/net/html"
    "strings"
)

type WalkEvent int
const (
    WalkEnter WalkEvent = iota
    WalkExit
)
    
func WalkHtmlDoc(doc *html.Node, cb func(*NodeWrapper, WalkEvent) bool) {
    var f func(n *html.Node)
    f = func(n *html.Node) {
        var wrappedNode *NodeWrapper
        wrappedNode = &NodeWrapper{
            Node: n,
        }
        if n.Type == html.ElementNode {
            wrappedNode.ElementType = n.Data
            for _, attr := range n.Attr {
                if attr.Key == "class" {
                    wrappedNode.Classes = strings.Fields(attr.Val)
                } else if attr.Key == "id" {
                    wrappedNode.ID = attr.Val
                }
            }
        }
        cont := cb(wrappedNode, WalkEnter)
        if (!cont) { return }

        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
        cb(wrappedNode, WalkExit)
    }
    f(doc)
}

