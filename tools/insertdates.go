package main

import (
    "yarrienet/htmlhelper"
    "yarrienet/rsshelper"
    "golang.org/x/net/html"
    "fmt"
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
    var nestedInDateDiv = false
    // if a time element has already been encountered before the p tag,
    // don't add another
    var existingTimeElement = false
    htmlhelper.WalkHtmlDoc(doc, func(wrappedNode *htmlhelper.NodeWrapper, event htmlhelper.WalkEvent) bool {
        if wrappedNode.ElementType == "div" {
            if slices.Contains(wrappedNode.Classes, "post") {
                if event == htmlhelper.WalkEnter {
                    // entering a new post
                    postId = wrappedNode.ID
                } else {
                    // exiting a new post, reset state
                    postId = ""
                }
            } else if slices.Contains(wrappedNode.Classes, "date") {
                nestedInDateDiv = event == htmlhelper.WalkEnter
                if event == htmlhelper.WalkExit {
                    // exiting a time element, reset state
                    existingTimeElement = false
                    return true
                }
            }
        }

        if postId != "" && nestedInDateDiv && event == htmlhelper.WalkEnter {
            if wrappedNode.ElementType == "time" {
                existingTimeElement = true
            } else if wrappedNode.ElementType == "p" && !existingTimeElement {
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
                    nodeParent.InsertBefore(dateNode, node.NextSibling)
                } else {
                    nodeParent.AppendChild(dateNode)
                }
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

const usageString = "insertdates <microblog file> <rss file> [output]"
func printUsage() {
    fmt.Println("USAGE")
    fmt.Println("  " + usageString)
    fmt.Println("")
    fmt.Println("OUTPUT")
    fmt.Println("  Wrap each yarrie.net microblog date <p> in <time datetime>")
    fmt.Println("  sourced from a linked RSS feed.")
    fmt.Println("")
    fmt.Println("  If no output file is provided then stdout is used.")
}

func printMissingArg(arg string) {
    fmt.Fprintln(os.Stderr, fmt.Sprintf("missing argument: %s", arg))
    fmt.Fprintln(os.Stderr, fmt.Sprintf("usage: %s", usageString))
}

func main() {
    if len(os.Args) < 2 {
        printMissingArg("microblog file")
        os.Exit(1)
    }
    if len(os.Args) < 3 {
        printMissingArg("rss file")
        os.Exit(1)
    }
    microblogFilePath := os.Args[1]
    rssFilePath := os.Args[2]

    // 1. determine the dates for each id
    rssFile, err := os.Open(rssFilePath)
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
    f, err := os.Open(microblogFilePath)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    doc, err := html.Parse(f)
    if err != nil {
        panic(err)
    }

    doc = insertDateNodes(doc, dates)

    if len(os.Args) > 3 {
        // output into a file
        outputFile, err := os.OpenFile(os.Args[3], os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
        if err != nil {
            panic(err)
        }
        defer outputFile.Close()

        err = html.Render(outputFile, doc)
        if err != nil {
            panic(err)
        }
    } else {
        // output into stdout
        var b strings.Builder
        err = html.Render(&b, doc)
        if err != nil {
            panic(err)
        }

        fmt.Println(b.String())
    }
}

