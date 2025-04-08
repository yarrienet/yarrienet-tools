package main
// package microblog

import (
    "fmt"
    "yarrienet/htmlhelper"
    "golang.org/x/net/html"
    "os"
    "time"
)
const postTemplate = `<div class="post" id="%[1]s">
    <div class="date">
        <a href="#%[1]s" class="post-link"><time datetime="%[2]s"><p>%[3]s</p></time></a>
    </div>
    <p></p>
</div>`

// %[1]s post id
// %[2]s iso 8601 datetime (?)
// %[3]s formatted time
func generatePost(id string, time time.Time) string {
    // Post.ID string
    // Post.DatePosted time.Time
    // Nodes []*html.Node    
    return fmt.Sprintf(postTemplate, id, "iso8601time", "formattedtime")
}

func NewPost(doc *html.Node) {
    htmlhelper.WalkHtmlDoc(doc, func (wn *htmlhelper.NodeWrapper, e htmlhelper.WalkEvent) bool {
        if e != htmlhelper.WalkEnter || wn.ID != "posts" {
            return true
        }
        fmt.Println(generatePost("exampleid", time.Now()))
        return false
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

    NewPost(doc)
}

