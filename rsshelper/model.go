package rsshelper

import (
    "encoding/xml"
    "time"
)

type RSS struct {
    XMLName xml.Name `xml:"rss"`
    Version string `xml:"version,attr"`
    Channel Channel `xml:"channel"`
    
}
type Channel struct {
    Title string `xml:"title"`
    Link string `xml:"link"`
    Description string `xml:"description"`
    Items []Item `xml:"item"`
}
type Item struct {
    ID string `xml:"guid"`
    Author string `xml:"author"`
    Link string `xml:"link"`
    Description string `xml:"description"`
    PubDate time.Time `xml:"pubDate"`
}

