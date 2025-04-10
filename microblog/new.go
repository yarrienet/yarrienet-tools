package microblog

import (
    "fmt"
    "yarrienet/htmlhelper"
    "golang.org/x/net/html"
    "os"
    "time"
    "strings"
    "slices"
)

// %[1]s post id
// %[2]s iso 8601 datetime (?)
// %[3]s formatted time
// spacing is important and dependant on correct indentation on insertion
const postTemplate = `
        <div class="post" id="%[1]s">
            <div class="date">
                <a href="#%[1]s" class="post-link"><time datetime="%[2]s"><p>%[3]s</p></time></a>
            </div>
            <p></p>
        </div>
`

var formattedMonths = []string{
    "jan", "feb", "march", "april", "may", "june", "july", "aug", "sept", "oct", "nov", "dec",
}

func generatePost(id string, datetime time.Time) string {
    // Post.ID string
    // Post.DatePosted time.Time
    // Nodes []*html.Node
    datetimeStr := datetime.Format(time.RFC3339)
    formattedStr := fmt.Sprintf("%s %d, %d", formattedMonths[datetime.Month()-1], datetime.Day(), datetime.Year())
    return fmt.Sprintf(postTemplate, id, datetimeStr, formattedStr)
}

func InsertNewPost(doc *html.Node, datetime time.Time) error {
    var err error = nil
    htmlhelper.WalkHtmlDoc(doc, func (wn *htmlhelper.NodeWrapper, e htmlhelper.WalkEvent) bool {
        // if err is present from a previous loop frame then exit
        if err != nil {
            return false
        }
        if e != htmlhelper.WalkEnter || wn.ID != "posts" {
            return true
        }
        // have found the posts div
        var postsDiv *html.Node = wn.Node

        postStr := generatePost("exampleid", datetime)

        fragment, err := html.ParseFragment(strings.NewReader(postStr), postsDiv)
        if err != nil {
            return false
        }
        for _, n := range slices.Backward(fragment) {
            postsDiv.InsertBefore(n, wn.Node.FirstChild)
        }

        return false
    })
    if err != nil {
        return err
    }
    return nil
}

// The function inserts an empty post element as the first child in #posts.
// Expects a file descriptor that can read and write to a file. WARNING:
// will truncate all contents of the file with the newly rendered document.
func InsertNewPostFile(f *os.File, datetime time.Time) error {
    doc, err := html.Parse(f)
    if err != nil {
        return err
    }
    err = InsertNewPost(doc, datetime)
    if err != nil {
        return err
    }
    
    f.Truncate(0)
    f.Seek(0, 0)
    err = html.Render(f, doc)
    if err != nil {
        return err
    }
    return nil
}

