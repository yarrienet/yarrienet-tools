package rsshelper

import (
    "testing"
)

var exampleFeed = `
<rss version="2.0">
    <channel>
        <title>yarrie</title>
        <link>http://yarrie.net/microblog</link>
        <description>yarrie&#39;s microblog</description>
        <item>
            <pubDate>Mon, 14 Apr 2025 12:26:44 +0100</pubDate>
            <guid>http://yarrie.net/microblog#exampleid</guid>
            <author>yarrie</author>
            <link>http://yarrie.net/microblog#exampleid</link>
            <description>&amp;lt;p&amp;gt;&amp;lt;/p&amp;gt;</description>
        </item>
    </channel>
</rss>`

func TestDecodeRss(t *testing.T) {
    // testing valid rss with date

    // Decode(data []byte) ([]Item, error)
}

