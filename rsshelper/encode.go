package rsshelper

import (
    "encoding/xml"
)

func (i *Item) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
    formattedDate := i.PubDate.Format("Mon, 02 Jan 2006 15:04:05 -0700")
    type Alias Item
    aux := &struct{
        PubDate string `xml:"pubDate"`
        *Alias
    }{
        PubDate: formattedDate,
        Alias: (*Alias)(i),
    } 
    return e.EncodeElement(aux, start)
}

