package main

import (
    "yarrienet/htmlhelper"
    "fmt"
    "golang.org/x/net/html"
    "os"
    "strings"
    "slices"
    "time"
)

func insertDateNodes(doc *html.Node) *html.Node {
    var parentIsDate = false
    htmlhelper.WalkHtmlDoc(doc, func(wrappedNode *htmlhelper.NodeWrapper, event htmlhelper.WalkEvent) bool {
        if wrappedNode.Type == "div" && slices.Contains(wrappedNode.Classes, "date") {
            parentIsDate = event == htmlhelper.WalkEnter
        }
        if parentIsDate && event == htmlhelper.WalkEnter {
            if wrappedNode.Type == "p" {
                node := wrappedNode.Node
                nodeParent := node.Parent
                if nodeParent != nil {
                    nodeParent.RemoveChild(node)

                    dateNode := htmlhelper.MakeDateNode(time.Now())
                    dateNode.AppendChild(node)
                    
                    if node.NextSibling != nil {
                        nodeParent.InsertBefore(&dateNode, node.NextSibling)
                    } else {
                        nodeParent.AppendChild(&dateNode)
                    }
                }
            }
        }

        // don't stop walking until end of doc
        return true
    })
    return doc
}
func main() {
    f, err := os.Open("microblog.html")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    doc, err := html.Parse(f)
    if err != nil {
        panic(err)
    }

    doc = insertDateNodes(doc)

    var b strings.Builder
    err = html.Render(&b, doc)
    if err != nil {
        panic(err)
    }

    fmt.Println(b.String())
}

