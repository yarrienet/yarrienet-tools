package microblog

import (
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

type RSSMetadata struct {
    Title string
    Author string
    Description string
    BaseUrl string
}

func getNodeClasses(n *html.Node) []string {
    return strings.Fields(htmlhelper.GetNodeAttr(n, "class"))
}

func parseMicroblog(doc *html.Node) []Post {
    var posts []Post

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
                    posts = append(posts, Post{
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

func postToRssItem(post Post, metadata *RSSMetadata) (*rsshelper.Item, error) {
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
    link := fmt.Sprintf("%s#%s", metadata.BaseUrl, post.ID)
    return &rsshelper.Item{
        ID: link,
        Author: metadata.Author,
        Link: link,
        Description: description,
        PubDate: post.DatePosted,
    }, nil
}

func GenRss(doc *html.Node, metadata *RSSMetadata) (string, error) {
    posts := parseMicroblog(doc)
    var items []rsshelper.Item
    for _, post := range posts {
        item, err := postToRssItem(post, metadata)
        if err != nil {
            return "", err
        }
        items = append(items, *item)
    }

    rssData := rsshelper.RSS{
        Version: "2.0",
        Channel: rsshelper.Channel{
            Title: metadata.Title,
            Link: metadata.BaseUrl,
            Description: metadata.Description,
            Items: items,
        },
    }
    data, err := xml.MarshalIndent(rssData, "", "    ")
    if err != nil {
        return "", err
    }
    return string(data), nil
}

func GenRssFromFile(f *os.File, metadata *RSSMetadata) (string, error) {
    doc, err := html.Parse(f)
    if err != nil {
        return "", err
    }
    rss, err := GenRss(doc, metadata)
    if err != nil {
        return "", err
    }
    return rss, nil
}

