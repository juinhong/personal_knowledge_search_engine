package models

type Note struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	Type    string   `json:"type"`
}

type SearchQuery struct {
	Query MultiMatchContainer `json:"query"`
}

type MultiMatchContainer struct {
	MultiMatch MultiMatchQuery `json:"multi_match"`
}

type MultiMatchQuery struct {
	Query  string   `json:"query"`
	Fields []string `json:"fields"`
}
