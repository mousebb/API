package products

// Content ...
type Content struct {
	Text        string      `bson:"text" json:"text" xml:"text"`
	ContentType ContentType `json:"contentType" xml:"contentType"`
	Sort        int         `json:"sort" xml:"sort"`
}

// ContentType ...
type ContentType struct {
	Id         int    `json:"id" xml:"id"`
	Type       string `bson:"type" json:"type" xml:"type"`
	AllowsHTML bool   `bson:"allows_html" json:"allows_html" xml:"allows_html"`
}
