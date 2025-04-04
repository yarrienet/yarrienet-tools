package microblog

import (
    "golang.org/x/net/html"
    "time"
)

type Post struct {
    ID string
    DatePosted time.Time
    Nodes []*html.Node 
}

