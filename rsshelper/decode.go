package rsshelper

import (
    "fmt"
    "encoding/xml"
    "time"
)

func (i *Item) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    type Alias Item
    aux := &struct{
        PubDate string `xml:"pubDate"`
        *Alias
    }{
        Alias: (*Alias)(i),
    }
    err := d.DecodeElement(aux, &start)
    if err != nil {
        return err
    }

    t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 ", aux.PubDate)
    if err != nil {
        return fmt.Errorf("error parsing pubDate: %v", err)
    }

    i.PubDate = t
    return nil
}

func Decode(data []byte) ([]Item, error) {
    var rss RSS
    err := xml.Unmarshal(data, &rss)
    if err != nil {
        return nil, err
    }
    return rss.Channel.Items, nil
}

