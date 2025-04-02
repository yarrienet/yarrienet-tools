package rsshelper

import (
    "time"
)

type RSS struct {
    Channel Channel `xml:"channel"`
}
type Channel struct {
    Items []Item `xml:"item"`
}
type Item struct {
    ID string `xml:"guid"`
    Title string `xml:"title"`
    Link string `xml:"link"`
    Description string `xml:"description"`
    PubDate time.Time `xml:"pubDate"`
}

