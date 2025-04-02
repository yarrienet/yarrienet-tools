package main

import (
    "yarrienet/htmlhelper"
    "yarrienet/rsshelper"
    "fmt"
    "golang.org/x/net/html"
    "os"
    "strings"
    "slices"
    "time"
    "io"
)

func insertDateNodes(doc *html.Node, dates map[string]time.Time) *html.Node {
    // tracks the current id of the post
    var postId string
    // if node is currently nested within <div class="date">
    var parentIsDate = false
    htmlhelper.WalkHtmlDoc(doc, func(wrappedNode *htmlhelper.NodeWrapper, event htmlhelper.WalkEvent) bool {
        if wrappedNode.Type == "div" {
            if slices.Contains(wrappedNode.Classes, "post") {
                if event == htmlhelper.WalkEnter {
                    postId = wrappedNode.ID
                } else {
                    postId = ""
                }
            } else if slices.Contains(wrappedNode.Classes, "date") {
                parentIsDate = event == htmlhelper.WalkEnter
            }
        }

        if postId != "" && parentIsDate && event == htmlhelper.WalkEnter && wrappedNode.Type == "p" {
            // get the date
            postDate, exists := dates[postId]
            if !exists {
                // if doesnt exist in map then skip date insert
                return true
            }

            node := wrappedNode.Node
            nodeParent := node.Parent
            if nodeParent == nil {
                return true
            }
            nodeParent.RemoveChild(node)

            dateNode := htmlhelper.MakeDateNode(postDate)
            dateNode.AppendChild(node)
            
            if node.NextSibling != nil {
                nodeParent.InsertBefore(&dateNode, node.NextSibling)
            } else {
                nodeParent.AppendChild(&dateNode)
            }
        }

        // don't stop walking until end of doc
        return true
    })
    return doc
}

func determineRssDates(data []byte) (map[string]time.Time, error) {
    items, err := rsshelper.Decode(data)
    if err != nil {
        return nil, err
    }

    postDates := make(map[string]time.Time, len(items))
    for _, item := range items {
        postDates[item.ID] = item.PubDate
    }
    return postDates, nil
}

func main() {
    // 1. determine the dates for each id
    rssFile, err := os.Open("rss.xml")
    if err != nil {
        panic(err)
    }
    defer rssFile.Close()

    data, err := io.ReadAll(rssFile)
    if err != nil {
        panic(err)
    }
    dates, err := determineRssDates(data) 
    if err != nil {
        panic(err)
    }

    // 2. update the dates for each post
    f, err := os.Open("microblog.html")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    doc, err := html.Parse(f)
    if err != nil {
        panic(err)
    }

    doc = insertDateNodes(doc, dates)

    var b strings.Builder
    err = html.Render(&b, doc)
    if err != nil {
        panic(err)
    }

    fmt.Println(b.String())
}

