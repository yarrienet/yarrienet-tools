package main

import (
    "yarrienet/microblog"
    "yarrienet/htmlhelper"
    "yarrienet/rsshelper"
    "golang.org/x/net/html"
    h "html"
    "fmt"
    "os"
    "strings"
    "slices"
    "time"
    "encoding/xml"
)

var title = "yarrie"
var author = "yarrie"
var description = "yarrie's microblog"
var baseUrl = "http://yarrie.net/microblog"

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

func postToRssItem(post microblog.Post) (*rsshelper.Item, error) {
    var b strings.Builder
    for _, node := range post.Nodes {
        if node.Type != html.ElementNode {
            continue
        }
        err := html.Render(&b, node)
        // render nodes
        if err != nil {
            return nil, err
        }
    }
    rendered := b.String()
    description := h.EscapeString(rendered)

    // assemble item
    link := fmt.Sprintf("%s#%s", baseUrl, post.ID)
    return &rsshelper.Item{
        ID: link,
        Author: author,
        Link: link,
        Description: description,
        PubDate: post.DatePosted,
    }, nil
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
    var items []rsshelper.Item
    for _, post := range posts {
        item, err := postToRssItem(post)
        if err != nil {
            panic(err)
        }
        items = append(items, *item)
    }

    rssData := rsshelper.RSS{
        Version: "2.0",
        Channel: rsshelper.Channel{
            Title: title,
            Link: baseUrl,
            Description: description,
            Items: items,
        },
    }
    data, err := xml.MarshalIndent(rssData, "", "    ")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(data))
}

