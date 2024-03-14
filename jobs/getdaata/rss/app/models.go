package app

import "encoding/xml"

// input struct(s)
// format = 1
type RssChannelFormat struct {
	XMLName xml.Name `xml:"rss"`
	Rss     Channel  `xml:"channel"`
}
type Channel struct {
	// XMLName xml.Name `xml:"channel"`
	Items []Item `xml:"item"`
}
type Item struct {
	// XMLName     xml.Name     `xml:"item"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Creator     string `xml:"creator"`
	PubDate     string `xml:"pubDate"`
}

// format = 2
type FeedFormat struct {
	Feed xml.Name   `xml:"feed"`
	Rss  []FeedItem `xml:"entry"`
}

type FeedItem struct {
	// XMLName     xml.Name     `xml:"item"`
	Title       string     `xml:"title"`
	Link        string     `xml:"id"`
	Description string     `xml:"summary"`
	Content     string     `xml:"content"`
	Creator     FeedAuthor `xml:"author"`
	PubDate     string     `xml:"published"`
}
type FeedAuthor struct {
	Name []string `xml:"name"`
}

// output struct
type DataFormatOut struct {
	Items []ItemsOut
}
type ItemsOut struct {
	Title       string
	URL         string
	Description string
	Creator     string
	PubDate     string
	ImageURL    string
	ImageCredit string
}
