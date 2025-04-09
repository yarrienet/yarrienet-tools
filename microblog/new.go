package main
// package microblog

import (
    "fmt"
    "yarrienet/htmlhelper"
    "golang.org/x/net/html"
    "os"
    "time"
    "strings"
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

func NewPost(doc *html.Node) error {
    var err error = nil
    htmlhelper.WalkHtmlDoc(doc, func (wn *htmlhelper.NodeWrapper, e htmlhelper.WalkEvent) bool {
        // if err is present from a previous loop frame then exit
        if err != nil {
            return false
        }
        if e != htmlhelper.WalkEnter || wn.ID != "posts" {
            return true
        }
        postStr := generatePost("exampleid", time.Now())

        fragment, err := html.ParseFragment(strings.NewReader(postStr), wn.Node)
        if err != nil {
            return false
        }
        for _, n := range fragment {
            // TODO should insert at the first TextElement?
            wn.Node.InsertBefore(n, wn.Node.FirstChild)
        }

        return false
    })
    if err != nil {
        return err
    }
    return nil
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

    var b strings.Builder
    err = html.Render(&b, doc)
    if err != nil {
        panic(err)
    }
    fmt.Println(b.String())
}

