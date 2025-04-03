package main

import (
    //"yarrienet/microblog"
    "yarrienet/htmlhelper"
    "golang.org/x/net/html"
    "fmt"
    "os"
    "strings"
    "slices"
)

func getNodeClasses(n *html.Node) []string {
    return strings.Fields(htmlhelper.GetNodeAttr(n, "class"))
}

func parseMicroblog(doc *html.Node) /*[]microblog.Post*/ {
    //var postId string

    // 1. walk down to find the posts div
    htmlhelper.WalkHtmlDoc(doc, func (wn *htmlhelper.NodeWrapper, e htmlhelper.WalkEvent) bool {
        if wn.ID == "posts" {
            // 2. loop each post
            n := wn.Node
            for c := n.FirstChild; c != nil; c = c.NextSibling {
                if c.Type != html.ElementNode {
                    continue
                }
                var id string
                var classes []string
                for _, attr := range c.Attr {
                    if attr.Key == "id" {
                        id = attr.Val 
                    } else if attr.Key == "class" {
                        classes = strings.Fields(attr.Val)
                    }
                }
                if c.Data == "div" && slices.Contains(classes, "post") {
                    // 3. root level elements of each post
                    fmt.Printf("got a post: %s\n", id)
                    for c2 := c.FirstChild; c2 != nil; c2 = c2.NextSibling {
                        if c2.Type != html.ElementNode {
                            continue
                        }
                        classes = getNodeClasses(c2)
                        if c2.Data == "div" && slices.Contains(classes, "date") {
                            // 4. walk to find <time datetime>
                            htmlhelper.WalkHtmlDoc(c2, func (wn2 *htmlhelper.NodeWrapper, e2 htmlhelper.WalkEvent) bool {
                                if e2 == htmlhelper.WalkEnter && wn2.Type == "time" {
                                    dateString := htmlhelper.GetNodeAttr(wn2.Node, "datetime")
                                    fmt.Printf("date string: %s\n", dateString)
                                    return false
                                }
                                return true
                            })
                        } else {
                            // 5. else add the node to the post
                            fmt.Printf("have a node: %s\n", c2.Data)
                        }
                    }
                }
            }
            return false
        }
        return true
    })
}

func main() {
    f, err := os.Open("microblog-fix.html")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    doc, err := html.Parse(f)
    if err != nil {
        panic(err)
    }

    parseMicroblog(doc)
}

