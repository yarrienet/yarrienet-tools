package htmlhelper

import (
    "golang.org/x/net/html"
    "time"
)

func MakeDateNode(datetime time.Time) html.Node {
    formattedTime := datetime.Format(time.RFC3339)
    return html.Node{
        Type: html.ElementNode,
        Data: "time",
        Attr: []html.Attribute{
            { Key: "datetime", Val: formattedTime },
        },
    }
}

// NodeWrapper is a helper struct that wraps an html.Node and provides
// convenience fields such as the node type and associated classes.
type NodeWrapper struct {
    Node *html.Node
    Type string
    Classes []string    
}

type WalkEvent int
const (
    WalkEnter WalkEvent = iota
    WalkExit
)

