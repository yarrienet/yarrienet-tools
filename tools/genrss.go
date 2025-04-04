package main

import (
    "yarrienet/microblog"
    "yarrienet/htmlhelper"
    "golang.org/x/net/html"
    "fmt"
    "os"
    "strings"
    "slices"
    "time"
)

func getNodeClasses(n *html.Node) []string {
    return strings.Fields(htmlhelper.GetNodeAttr(n, "class"))
}

func parseMicroblog(doc *html.Node) []microblog.Post {
    var posts []microblog.Post

    var postId string
    var postDate time.Time
    var postNodes []*html.Node

    var nestedInPostDate = false

    htmlhelper.WalkHtmlDoc(doc, func (wn *htmlhelper.NodeWrapper, e htmlhelper.WalkEvent) bool {
        if slices.Contains(wn.Classes, "post") {
            if e == htmlhelper.WalkEnter {
                postId = wn.ID
            } else {
                if !postDate.IsZero() {
                    posts = append(posts, microblog.Post{
                        ID: postId,
                        DatePosted: postDate,
                        Nodes: postNodes,
                    })
                } else {
                    fmt.Printf("warning, post %s is missing date, skipping\n", postId)
                }
                postId = ""
                postDate = time.Time{}
                postNodes = nil
                return false
            }
        } else if postId != "" {
            if wn.ElementType == "div" && slices.Contains(wn.Classes, "date") {
                nestedInPostDate = e == htmlhelper.WalkEnter
            } else if e == htmlhelper.WalkEnter {
                if nestedInPostDate {
                    if wn.ElementType == "time" {
                        postDateString := htmlhelper.GetNodeAttr(wn.Node, "datetime")
                        postDate, _ = time.Parse(time.RFC3339, postDateString)
                        return false
                    }
                } else {
                    postNodes = append(postNodes, wn.Node)
                    return false 
                }
            }
        }
        return true
    })
    return posts
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

    posts := parseMicroblog(doc)
    for _, post := range posts {
        fmt.Printf("%s - %s\n", post.ID, post.DatePosted.Format(time.RFC3339))
        for _, node := range post.Nodes {
            fmt.Printf("  %s\n", node.Data)
        } 
    }
}

