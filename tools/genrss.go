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
    var postId string
    var nestedInPostDate = false
    htmlhelper.WalkHtmlDoc(doc, func (wn *htmlhelper.NodeWrapper, e htmlhelper.WalkEvent) bool {
        if postId == "" && slices.Contains(wn.Classes, "post") {
            if e == htmlhelper.WalkEnter {
                postId = wn.ID
                fmt.Printf("have a post: %s\n", postId)
            } else {
                postId = ""
                return false
            }
        } else if postId != "" {
            if wn.Type == "div" && slices.Contains(wn.Classes, "date") {
                nestedInPostDate = e == htmlhelper.WalkEnter
            } else if nestedInPostDate {
                if wn.Type == "time" {
                    postDate := htmlhelper.GetNodeAttr(wn.Node, "datetime")
                    fmt.Printf("post date: %s\n", postDate)
                    return false
                }
            } else {
                fmt.Printf("post node: '%s'\n", wn.Type)
                return false 
            }
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

