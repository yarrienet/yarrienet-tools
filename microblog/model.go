package microblog

import (
    "time"
)

type Post struct {
    ID string
    DatePosted time.Time
    Content string
}

